package bubbletea

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/wish"
	"github.com/gliderlabs/ssh"
	"github.com/muesli/termenv"
)

// BubbleTeaHandler is the function Bubble Tea apps implement to hook into the
// SSH Middleware. This will create a new tea.Program for every connection and
// start it with the tea.ProgramOptions returned.
//
// Deprecated: use Handler instead.
type BubbleTeaHandler = Handler // nolint: revive

// Handler is the function Bubble Tea apps implement to hook into the
// SSH Middleware. This will create a new tea.Program for every connection and
// start it with the tea.ProgramOptions returned.
type Handler func(ssh.Session) (tea.Model, []tea.ProgramOption)

// ProgramHandler is the function Bubble Tea apps implement to hook into the SSH
// Middleware. This should return a new tea.Program. This handler is different
// from the default handler in that it returns a tea.Program instead of
// (tea.Model, tea.ProgramOptions).
//
// Make sure to set the tea.WithInput and tea.WithOutput to the ssh.Session
// otherwise the program will not function properly.
type ProgramHandler func(ssh.Session) *tea.Program

// Middleware takes a Handler and hooks the input and output for the
// ssh.Session into the tea.Program. It also captures window resize events and
// sends them to the tea.Program as tea.WindowSizeMsgs. By default a 256 color
// profile will be used when rendering with Lip Gloss.
func Middleware(bth Handler) wish.Middleware {
	return MiddlewareWithColorProfile(bth, termenv.ANSI256)
}

// MiddlewareWithColorProfile allows you to specify the number of colors
// returned by the server when using Lip Gloss. The number of colors supported
// by an SSH client's terminal cannot be detected by the server but this will
// allow for manually setting the color profile on all SSH connections.
func MiddlewareWithColorProfile(bth Handler, cp termenv.Profile) wish.Middleware {
	h := func(s ssh.Session) *tea.Program {
		m, opts := bth(s)
		if m == nil {
			return nil
		}
		opts = append(opts, tea.WithInput(s), tea.WithOutput(s))
		return tea.NewProgram(m, opts...)
	}
	return MiddlewareWithProgramHandler(h, cp)
}

// MiddlewareWithProgramHandler allows you to specify the ProgramHandler to be
// able to access the underlying tea.Program. This is useful for creating custom
// middlewars that need access to tea.Program for instance to use p.Send() to
// send messages to tea.Program.
//
// Make sure to set the tea.WithInput and tea.WithOutput to the ssh.Session
// otherwise the program will not function properly.
func MiddlewareWithProgramHandler(bth ProgramHandler, cp termenv.Profile) wish.Middleware {
	return func(sh ssh.Handler) ssh.Handler {
		lipgloss.SetColorProfile(cp)
		return func(s ssh.Session) {
			errc := make(chan error, 1)
			p := bth(s)
			if p != nil {
				_, windowChanges, _ := s.Pty()
				go func() {
					for {
						select {
						case <-s.Context().Done():
							if p != nil {
								p.Quit()
							}
						case w := <-windowChanges:
							if p != nil {
								p.Send(tea.WindowSizeMsg{Width: w.Width, Height: w.Height})
							}
						case err := <-errc:
							if err != nil {
								log.Print(err)
							}
						}
					}
				}()
				errc <- p.Start()
				// p.Kill() will force kill the program if it's still running,
				// and restore the terminal to its original state in case of a
				// tui crash
				p.Kill()
			}
			sh(s)
		}
	}
}
