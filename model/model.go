package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
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
		m.viewport.Height = msg.Height
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
	out, _ := glamour.Render(m.Slides[m.Page], "dark")

	slideStyle := lipgloss.NewStyle().
		Width(m.viewport.Width).
		Height(m.viewport.Height - 1).
		Align(lipgloss.Left)

	authorStyle := lipgloss.NewStyle().
		Height(1).
		Foreground(lipgloss.Color("#E8B4BC")).
		Align(lipgloss.Left)

	dateStyle := lipgloss.NewStyle().
		Height(1).
		Faint(true).
		Align(lipgloss.Left)

	pageStyle := lipgloss.NewStyle().
		Height(1).
		Foreground(lipgloss.Color("#E8B4BC")).
		Align(lipgloss.Right)

	return lipgloss.JoinVertical(lipgloss.Left,
		slideStyle.Render(out),
		lipgloss.JoinHorizontal(lipgloss.Center,
			authorStyle.Render(m.Author),
			dateStyle.Render(m.Date),
			pageStyle.Render(fmt.Sprintf("Slide %d / %d", m.Page, len(m.Slides)-1)),
		),
	)
}
