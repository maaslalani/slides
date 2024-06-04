package remote

import tea "github.com/charmbracelet/bubbletea"

// CommandRelay is meant to expose slide interaction to external
// processes that can work as a remote for the slides.
type CommandRelay struct {
	*tea.Program
}

func NewCommandRelay(p *tea.Program) *CommandRelay {
	return &CommandRelay{
		Program: p,
	}
}

func (r *CommandRelay) SlideNext() {
	r.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'n'},
	})
}

func (r *CommandRelay) SlidePrev() {
	r.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'p'},
	})
}

func (r *CommandRelay) SlideFirst() {
	// Requires 2 keystrokes to actually
	// move to first slide
	r.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'g'},
	})
	r.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'g'},
	})
}

func (r *CommandRelay) SlideLast() {
	r.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'G'},
	})
}

func (r *CommandRelay) CodeExecute() {
	r.Send(tea.KeyMsg{
		Type: tea.KeyCtrlE,
	})
}

func (r *CommandRelay) CodeCopy() {
	r.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'y'},
	})
}

func (r *CommandRelay) Quit() {
	r.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	})
}
