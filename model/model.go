package model

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type Model struct {
	Slides   []string
	Page     int
	Author   string
	Date     string
	viewport viewport.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		return m, nil

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
	g, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(m.viewport.Width),
		glamour.WithAutoStyle(),
	)
	out, err := g.Render(m.Slides[m.Page])

	if err != nil {
		return `Error: Invalid Markdown!`
	}

	return out
}
