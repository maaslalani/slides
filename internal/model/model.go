package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/maaslalani/slides/styles"
)

type Model struct {
	Slides   []string
	Page     int
	Author   string
	Date     string
	Theme    glamour.TermRendererOption
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
		case " ", "down", "k", "right", "l", "enter", "n":
			if m.Page < len(m.Slides)-1 {
				m.Page++
			}
		case "up", "j", "left", "h", "p":
			if m.Page > 0 {
				m.Page--
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	r, _ := glamour.NewTermRenderer(m.Theme, glamour.WithWordWrap(0))
	slide, err := r.Render(m.Slides[m.Page])
	if err != nil {
		slide = fmt.Sprintf("Error: Could not render markdown! (%v)", err)
	}
	slide = styles.Slide.Render(slide)

	left := styles.Author.Render(m.Author) + styles.Date.Render(m.Date)
	right := styles.Page.Render(fmt.Sprintf("Slide %d / %d", m.Page+1, len(m.Slides)))
	status := styles.Status.Render(styles.JoinHorizontal(left, right, m.viewport.Width))
	return styles.JoinVertical(slide, status, m.viewport.Height)
}
