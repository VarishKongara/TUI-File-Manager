package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/VarishKongara/TUI-File-Manager/filemanager"
)

type model struct {
	filemanager filemanager.Model
}

func (m model) Init() tea.Cmd {
	return m.filemanager.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.filemanager, cmd = m.filemanager.Update(msg)
			return m, cmd
		}
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
	return str.String()
}

func main() {
	filemanager := filemanager.New(1, ".")
	app := tea.NewProgram(model{filemanager: filemanager})
	if _, err := app.Run(); err != nil {
		fmt.Print("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Finished")
}
