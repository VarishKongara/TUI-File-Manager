package cmdline

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ID int
	Height int

	// Keybinds
	KeyMap KeyMap
}

func New(id int, keyMap KeyMap) Model {
	return Model{
		ID:  id,
        Height: 3,
        KeyMap: keyMap,
	}
}

type KeyMap struct {
	Escape       key.Binding
}



func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "\n\n\n"
}
