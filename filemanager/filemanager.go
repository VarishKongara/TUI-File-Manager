package filemanager

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ID    int
	CWD   string //current working directory
	Files []os.DirEntry
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
	case direntMsg:
		if msg.id != m.ID {
			break
		}
		m.Files = msg.dirents
	}

	return m, nil
}

func (m Model) View() string {
	var str strings.Builder
	for _, file := range m.Files {
		name := file.Name()
		str.WriteString(name)
		str.WriteRune('\n')
	}

	return str.String()
}
