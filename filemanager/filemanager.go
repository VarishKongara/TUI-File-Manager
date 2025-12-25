package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
	}
}

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	OpenFile key.Binding
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
	OpenFile: key.NewBinding(key.WithKeys("l", "right"),
		key.WithHelp("l", "open file"),
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
		case key.Matches(msg, DefaultKeyMap.OpenFile):
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
		}
	}

	return m, nil
}

func (m Model) View() string {
	var str strings.Builder
	if len(m.Files) > 0 {
		str.WriteString(m.Files[m.Selected].Name())
		str.WriteRune('\n')
		str.WriteRune('\n')
	}
	for i, file := range m.Files {
		if i < m.Top || i > m.Bottom {
			continue
		}

		name := file.Name()
		str.WriteString(name)
		str.WriteRune('\n')
	}

	return str.String()
}
