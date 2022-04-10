package server

import (
	"context"
	"fmt"

	"github.com/charmbracelet/wish"
	"github.com/gliderlabs/ssh"
	"github.com/maaslalani/slides/internal/model"
)

type Server struct {
	host  string
	port  int
	srv   *ssh.Server
	presentation model.Model
}

// NewServer creates a new server.
func NewServer(keyPath, host string, port int, presentation model.Model) (*Server, error) {
	s := &Server{
		host:  host,
		port:  port,
		presentation: presentation,
	}
	srv, err := wish.NewServer(
		wish.WithHostKeyPath(keyPath),
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithMiddleware(
			slidesMiddleware(s),
		),
	)
	if err != nil {
		return nil, err
	}
	s.srv = srv
	return s, nil
}

// Start starts the ssh server.
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// Shutdown shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
