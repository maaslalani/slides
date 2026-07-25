package server

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	bm "charm.land/wish/v2/bubbletea"
	"github.com/charmbracelet/ssh"
)

func slidesMiddleware(srv *Server) wish.Middleware {
	newProg := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		p := tea.NewProgram(m, opts...)
		return p
	}
	teaHandler := func(s ssh.Session) *tea.Program {
		_, _, active := s.Pty()
		if !active {
			fmt.Println("no active terminal, skipping")
			err := s.Exit(1)
			if err != nil {
				fmt.Println("Error exiting session")
			}
			return nil
		}
		return newProg(srv.presentation, bm.MakeOptions(s)...)
	}
	return bm.MiddlewareWithProgramHandler(teaHandler)
}
