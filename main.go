package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	str string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			m.str += msg.String()
		}
	}

	return m, nil
}

func (m model) View() string {
	return m.str
}

func main() {
	app := tea.NewProgram(model{str: "Hello, World!"})
	if _, err := app.Run(); err != nil {
		fmt.Print("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Finished")
}
