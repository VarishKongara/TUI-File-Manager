package cmdline

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ID int
	Height int
}

func New(id int) Model {
	return Model{
		ID:  id,
        Height: 3,
	}
}

type KeyMap struct {
	Escape       key.Binding
}

var DefaultKeyMap = KeyMap{
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "escape"),
	),
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	switch {
	// 	case key.Matches(msg, DefaultKeyMap.Escape):
	// 	return m,
	// 	}
	// }
	return m, nil
}

func (m Model) View() string {
	return "\n\n\n"
}
