package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ID    int
	CWD   string //current working directory
	Files []os.DirEntry

	Size     int
	Top      int // top item displayed
	Bottom   int // bottom item displayed
	Selected int // selected item

	Margin int // number of lines above and below the file display

	// Styles
	PermStyles PermStyles
}

// Styles for permissions string output
type PermStyles struct {
	Dir     lipgloss.Style
	File    lipgloss.Style
	Symlink lipgloss.Style

	Read  lipgloss.Style
	Write lipgloss.Style
	Exec  lipgloss.Style
	None  lipgloss.Style

	Special lipgloss.Style
}

func DefaultPermStyles() PermStyles {
	return PermStyles{
		Dir: lipgloss.NewStyle().Foreground(lipgloss.Color("69")). // blue
										Bold(true),
		File: lipgloss.NewStyle().Foreground(lipgloss.Color("252")). // light gray
										Bold(true),
		Symlink: lipgloss.NewStyle().Foreground(lipgloss.Color("81")). // cyan
										Bold(true),

		Read:  lipgloss.NewStyle().Foreground(lipgloss.Color("214")), // yellow
		Write: lipgloss.NewStyle().Foreground(lipgloss.Color("196")), // red
		Exec:  lipgloss.NewStyle().Foreground(lipgloss.Color("76")),  // green
		None:  lipgloss.NewStyle().Foreground(lipgloss.Color("240")), // dim gray

		Special: lipgloss.NewStyle().
			Foreground(lipgloss.Color("197")). // red/magenta
			Bold(true),
	}
}

func New(id int, cwd string) Model {
	return Model{
		ID:  id,
		CWD: cwd,

		Size:     0,
		Top:      0,
		Bottom:   0,
		Selected: 0,

		Margin: 4,

		PermStyles: DefaultPermStyles(),
	}
}

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Open   key.Binding
	Parent key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Open: key.NewBinding(key.WithKeys("l", "right"),
		key.WithHelp("l", "open file"),
	),
	Parent: key.NewBinding(key.WithKeys("h", "left"),
		key.WithHelp("l", "parent directory"),
	),
}

type direntMsg struct {
	id      int
	dirents []os.DirEntry
}

func (m Model) readDir(dir string) tea.Cmd {
	return func() tea.Msg {
		dirents, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
		}

		return direntMsg{id: m.ID, dirents: dirents}
	}
}

func (m Model) Init() tea.Cmd {
	return m.readDir(m.CWD)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg.Height - m.Margin
		m.Bottom = m.Size - 1
	case direntMsg:
		if msg.id != m.ID {
			break
		}
		m.Files = msg.dirents
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Up):
			if m.Selected > 0 {
				m.Selected--
			}
			if m.Selected < m.Top {
				m.Top--
				m.Bottom--
			}
		case key.Matches(msg, DefaultKeyMap.Down):
			m.Selected++
			if m.Selected >= len(m.Files) {
				m.Selected = len(m.Files) - 1
			}
			if m.Selected > m.Bottom {
				m.Top++
				m.Bottom++
			}
		case key.Matches(msg, DefaultKeyMap.Open):
			if len(m.Files) == 0 {
				break
			}
			if m.Files[m.Selected].IsDir() {

				m.CWD = filepath.Join(m.CWD, m.Files[m.Selected].Name())
				m.Selected = 0
				m.Top = 0
				m.Bottom = m.Size - 1
				return m, m.readDir(m.CWD)
			}
		case key.Matches(msg, DefaultKeyMap.Parent):
			m.CWD = filepath.Dir(m.CWD)
			m.Selected = 0
			m.Top = 0
			m.Bottom = m.Size - 1
			return m, m.readDir(m.CWD)
		}
	}

	return m, nil
}

func (m Model) View() string {
	var out strings.Builder
	if len(m.Files) > 0 && m.Selected < len(m.Files) {
		out.WriteString(m.Files[m.Selected].Name())
		out.WriteRune('\n')
		out.WriteRune('\n')
	}

	for i, file := range m.Files {
		if i < m.Top || i > m.Bottom {
			continue
		}

		info, err := os.Stat(filepath.Join(m.CWD, file.Name()))
		if err != nil {
			out.WriteString(strings.Repeat("?", 10))
			continue
		} else {
			writePerms(&out, info.Mode().String(), m.PermStyles)
		}
		out.WriteRune(' ')

		name := file.Name()
		if file.IsDir() {
			name = m.PermStyles.Dir.Render(name)
		}
		out.WriteString(name)
		out.WriteRune('\n')
	}

	return out.String()
}

func writePerms(out *strings.Builder, perms string, style PermStyles) {
	if out == nil {
		panic("String Builder is nil. In filemanager writePerms")
	}

	switch perms[0] {
	case 'd':
		out.WriteString(style.Dir.Render("d"))
	case 'l':
		out.WriteString(style.Symlink.Render("l"))
	default:
		out.WriteString(style.File.Render(string(perms[0])))
	}

	for _, char := range perms[1:] {
		switch char {
		case 'r':
			out.WriteString(style.Read.Render("r"))
		case 'w':
			out.WriteString(style.Write.Render("w"))
		case 'x':
			out.WriteString(style.Exec.Render("x"))
		case 's', 'S', 't', 'T':
			out.WriteString(style.Special.Render(string(char)))
		case '-':
			out.WriteString(style.None.Render("-"))
		default:
			out.WriteRune(char)
		}
	}
}
