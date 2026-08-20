package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/VarishKongara/TUI-File-Manager/cmdline"
	"github.com/VarishKongara/TUI-File-Manager/filemanager"
)

type model struct {
	filemanager filemanager.Model
	cmdline cmdline.Model
	focus int
}

type KeyMap struct {
	SwitchFocus key.Binding
	Quit key.Binding
}

var DefaultKeybinds = KeyMap {
	SwitchFocus: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "next component"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "exit program"),
	),

}

var DefaultFileManagerKeyMap = filemanager.KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	EnterDir: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l", "EnterDir file"),
	),
	Parent: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("l", "parent directory"),
	),
	Open: key.NewBinding(key.WithKeys("enter", "ctrl+o"),
		key.WithHelp("enter/ctrl+o", "open selection"),
	),
}

var DefaultCmdLineKeyMap = cmdline.KeyMap{
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "escape"),
	),
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.filemanager.Init(),
		m.cmdline.Init(),
	)
}

// IDs for each component
const (
	FileManager = iota
	CmdLine
	numComponents
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeybinds.Quit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeybinds.SwitchFocus):
			m.focus = (m.focus+1) % numComponents
			fallthrough
		default:
			var cmd tea.Cmd
			switch m.focus{
			case FileManager:
				m.filemanager, cmd = m.filemanager.Update(msg)

			case CmdLine:
				m.cmdline, cmd = m.cmdline.Update(msg)
			}
			return m, cmd
		}
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
        msg.Height = msg.Height - m.cmdline.Height;
		m.filemanager, cmd = m.filemanager.Update(msg)
		return m, cmd
	default:
		var cmd tea.Cmd
		m.filemanager, cmd = m.filemanager.Update(msg)
		return m, cmd
	}

	// return m, nil
}

func (m model) View() string {
	var str strings.Builder
	str.WriteString("\n" + m.filemanager.View() + "\n")
	str.WriteString(m.cmdline.View())
	return str.String()
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filemanager := filemanager.New(FileManager, cwd, DefaultFileManagerKeyMap)
	cmdline := cmdline.New(CmdLine, DefaultCmdLineKeyMap)
	app := tea.NewProgram(model{filemanager: filemanager, cmdline: cmdline, focus: FileManager}, tea.WithAltScreen())
	if _, err := app.Run(); err != nil {
		fmt.Print("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Finished")
}
