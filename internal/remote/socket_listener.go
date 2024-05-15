package remote

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// Default UNIX Socket file to be used as Socket Remote Listener
const SocketRemoteListenerDefaultPath string = "/tmp/slides.sock"

// Sane maximum socket buffer lengths as per current usage
const socketMaxReadBufLen int = 16
const socketMaxWriteBufLen int = 64

var socketCommandSlideFirst = []byte("s:first")
var socketCommandSlideNext = []byte("s:next")
var socketCommandSlidePrev = []byte("s:prev")
var socketCommandSlideLast = []byte("s:last")
var socketCommandCodeExec = []byte("c:exec")
var socketCommandCodeCopy = []byte("c:copy")
var socketCommandQuit = []byte("quit")

type SocketRemoteListener struct {
	net.Listener
	relay *CommandRelay
}

// Start listening on socket
func (s *SocketRemoteListener) Start() {
	go func() {
		for {
			var conn net.Conn
			var err error

			conn, err = s.Accept()
			if err != nil {
				// Nowhere to log or report error that
				// may happen.
				// Neither it makes sense to impact
				// the slides session for issue in
				// remote listening.
				continue
			}

			// handle accepted connection
			go s.handleConnection(conn)
		}
	}()
}

func (s *SocketRemoteListener) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, socketMaxReadBufLen)
	n, err := conn.Read(buf)
	if err != nil {
		writeSocketError(conn, err)
		return
	}

	if n == 0 {
		writeSocketError(conn, errors.New("invalid 0 length command"))
		return
	}

	args := strings.Split(string(buf[:n]), " ")
	command := args[0]

	switch command {
	case string(socketCommandSlideNext):
		s.relay.SlideNext()
	case string(socketCommandSlidePrev):
		s.relay.SlidePrev()
	case string(socketCommandSlideFirst):
		s.relay.SlideFirst()
	case string(socketCommandSlideLast):
		s.relay.SlideLast()
	case string(socketCommandCodeExec):
		s.relay.CodeExecute()
	case string(socketCommandCodeCopy):
		s.relay.CodeCopy()
	case string(socketCommandQuit):
		s.relay.Quit()
	default:
		writeSocketError(conn, errors.New("invalid command"))
		return
	}

	conn.Write([]byte("OK"))
}

// write error string on the connection
// this is meant to be a feedback to the client
func writeSocketError(conn net.Conn, err error) {
	conn.Write([]byte(fmt.Sprintf("ERR:%s", err)))
}

func NewSocketRemoteListener(socketFile string, relay *CommandRelay) (sock *SocketRemoteListener, err error) {
	socket, err := net.Listen("unix", socketFile)
	if err != nil {
		return nil, err
	}

	return &SocketRemoteListener{
		Listener: socket,
		relay:    relay,
	}, nil
}
