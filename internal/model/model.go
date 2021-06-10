package model

import (
	"fmt"
	"os/exec"
	"regexp"

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

	Code struct {
		Language string
		Block    string
		Result   string
		Display  bool
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

var LanguageCommands = map[string]struct {
	Command string
	Args    []string
}{
	"ruby": {
		Command: "ruby",
		Args:    []string{"-e"},
	},
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		m.Code.Display = false
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
		case "e":
			// Executing arbitrary code blocks and displaying
			// the output of the program
			m.Code.Language = "ruby"
			re := regexp.MustCompile("```.*\n(.*)\n```")
			slide := m.Slides[m.Page]
			match := re.FindStringSubmatch(slide)

			// There is no code block on the screen
			// skipping...
			if len(match) < 2 {
				return m, nil
			}

			m.Code.Block = match[1]
			lang := LanguageCommands[m.Code.Language]
			args := append(lang.Args, m.Code.Block)
			cmd := exec.Command(lang.Command, args...)
			out, err := cmd.Output()

			if err != nil {
				m.Code.Result = "Error: failed to execute code block"
			} else {
				m.Code.Result = string(out)
			}

			m.Code.Display = true
		}
	}
	return m, nil
}

func (m Model) View() string {
	r, _ := glamour.NewTermRenderer(m.Theme, glamour.WithWordWrap(0))
	slide := m.Slides[m.Page]
	if m.Code.Display {
		slide += "\n```\nResult: " + m.Code.Result + "\n```"
	}
	slide, err := r.Render(slide)
	if err != nil {
		slide = fmt.Sprintf("Error: Could not render markdown! (%v)", err)
	}
	slide = styles.Slide.Render(slide)

	left := styles.Author.Render(m.Author) + styles.Date.Render(m.Date)
	right := styles.Page.Render(fmt.Sprintf("Slide %d / %d", m.Page+1, len(m.Slides)))
	status := styles.Status.Render(styles.JoinHorizontal(left, right, m.viewport.Width))
	return styles.JoinVertical(slide, status, m.viewport.Height)
}
