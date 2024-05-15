package remote

import (
	"errors"
	"net"
)

// SocketRemote is the client (remote controller) that
// communicates with the SocketRemoteListener UNIX socket
type SocketRemote struct {
	net.Conn
}

func NewSocketRemote(socketFile string) (remote *SocketRemote, err error) {
	conn, err := net.Dial("unix", socketFile)
	if err != nil {
		return nil, err
	}

	return &SocketRemote{
		Conn: conn,
	}, nil
}

func (r *SocketRemote) writeCommand(command []byte) error {
	n, err := r.Write(command)
	if err != nil {
		return err
	}
	if n != len(command) {
		return errors.New("could not send complete data")
	}
	return nil
}

func (r *SocketRemote) SlideNext() error {
	return r.writeCommand(socketCommandSlideNext)
}

func (r *SocketRemote) SlidePrevious() error {
	return r.writeCommand(socketCommandSlidePrev)
}

func (r *SocketRemote) SlideFirst() error {
	return r.writeCommand(socketCommandSlideFirst)
}

func (r *SocketRemote) SlideLast() error {
	return r.writeCommand(socketCommandSlideLast)
}

func (r *SocketRemote) CodeExec() error {
	return r.writeCommand(socketCommandCodeExec)
}

func (r *SocketRemote) CodeCopy() error {
	return r.writeCommand(socketCommandCodeCopy)
}

func (r *SocketRemote) Quit() error {
	return r.writeCommand(socketCommandQuit)
}
