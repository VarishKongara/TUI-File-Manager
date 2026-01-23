package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	ID    int
	CWD   string //current working directory
	Files []os.DirEntry

	Selected int // selected item
	Margin   int // number of lines above and below the file display

	viewport    viewport.Model
	vport_ready bool

	// Styles
	PermStyles    PermStyles
	SelectedStyle SelectedStyle
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

// Inherit ligloss styles for PermStyles
func (ps PermStyles) Inherit(base lipgloss.Style) PermStyles {
	return PermStyles{
		Dir:     ps.Dir.Inherit(base),
		File:    ps.File.Inherit(base),
		Symlink: ps.Symlink.Inherit(base),
		Read:    ps.Read.Inherit(base),
		Write:   ps.Write.Inherit(base),
		Exec:    ps.Exec.Inherit(base),
		None:    ps.None.Inherit(base),
		Special: ps.Special.Inherit(base),
	}
}

func DefaultPermStyles() PermStyles {
	return PermStyles{
		Dir: lipgloss.NewStyle().Foreground(lipgloss.Color("69")). // blue
										Bold(true),
		File: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
			Dark:  "252", // light gray
			Light: "240", // dim gray
		},
		),
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

type SelectedStyle struct {
	Background lipgloss.Style
}

func DefaultSelectedStyle() SelectedStyle {
	return SelectedStyle{
		Background: lipgloss.NewStyle().Background(lipgloss.CompleteAdaptiveColor{
			Dark:  lipgloss.CompleteColor{ANSI256: "238", TrueColor: "#2e2e3e"}, // darker + blue tint
			Light: lipgloss.CompleteColor{ANSI256: "254", TrueColor: "#dcdcdc"}, // grey
		}),
	}
}

func New(id int, cwd string) Model {
	return Model{
		ID:  id,
		CWD: cwd,

		Selected: 0,
		Margin:   4,

		PermStyles:    DefaultPermStyles(),
		SelectedStyle: DefaultSelectedStyle(),
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
		height := msg.Height - m.Margin

		if !m.vport_ready {
			m.viewport = viewport.New(msg.Width, height)
			m.viewport.YPosition = m.Margin / 2
			m.vport_ready = true

			m.viewport.SetContent(m.renderFiles())
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = height
		}
	case direntMsg:
		if msg.id != m.ID {
			break
		}
		m.Files = msg.dirents
		m.viewport.SetContent(m.renderFiles())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Up):
			if m.Selected > 0 {
				m.Selected--
			}

			m.viewport.SetYOffset(max(0, m.Selected-m.viewport.Height/2))

			m.viewport.SetContent(m.renderFiles())
		case key.Matches(msg, DefaultKeyMap.Down):
			m.Selected++
			if m.Selected >= len(m.Files) {
				m.Selected = len(m.Files) - 1
			}

			m.viewport.SetYOffset(max(0, m.Selected-m.viewport.Height/2))

			m.viewport.SetContent(m.renderFiles())
		case key.Matches(msg, DefaultKeyMap.Open):
			if len(m.Files) == 0 {
				break
			}
			if m.Files[m.Selected].IsDir() {
				m.CWD = filepath.Join(m.CWD, m.Files[m.Selected].Name())
				m.Selected = 0
				m.viewport.GotoTop()
				return m, m.readDir(m.CWD)
			}
		case key.Matches(msg, DefaultKeyMap.Parent):
			m.CWD = filepath.Dir(m.CWD)
			m.viewport.GotoTop()
			return m, m.readDir(m.CWD)
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) renderFiles() string {
	var out strings.Builder
	for i, file := range m.Files {
		name := file.Name()

		// Get and style perms
		var permStr string
		info, err := os.Stat(filepath.Join(m.CWD, name))
		if err != nil {
			permStr = strings.Repeat("?", 10)
			continue
		} else {
			permStr = info.Mode().String()
		}

		// Style the output
		if i == m.Selected {
			// Highlight selected
			style := m.PermStyles.Inherit(m.SelectedStyle.Background)
			permStr = writePerms(permStr, style)
			if file.IsDir() {
				name = style.Dir.Render(name)
			} else {
				name = style.File.Render(name)
			}

			out.WriteString(permStr + m.SelectedStyle.Background.Render(" ") + name)
		} else {
			permStr = writePerms(permStr, m.PermStyles)
			if file.IsDir() {
				name = m.PermStyles.Dir.Render(name)
			} else {
				name = m.PermStyles.File.Render(name)
			}
			out.WriteString(permStr + " " + name)
		}

		out.WriteRune('\n')
	}

	return out.String()
}

func writePerms(perms string, style PermStyles) string {
	var out strings.Builder

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

	return out.String()
}

func (m Model) View() string {
	if !m.vport_ready {
		return "Loading..."
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.CWD,
		m.viewport.View(),
	)
}
