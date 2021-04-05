package model

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	Slides []string
	Page   int
	Author string
	Date   string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k", "right", "l", "enter", "n":
			if m.Page < len(m.Slides)-1 {
				m.Page++
			}
		case "down", "j", "left", "h", "p":
			if m.Page > 0 {
				m.Page--
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	return m.Slides[m.Page]
}
