package remote

import (
	"errors"
	"strings"
)

// parse socketPath to type and addr
// eg: tcp:localhost:8091 => "tcp", "localhost:8091"
// eg: unix:/tmp/slides.sock => "unix", "/tmp/slides.sock"
func parseSocketPath(path string) (socketType string, socketAddr string, err error) {
	parts := strings.Split(path, ":")
	if len(parts) < 2 {
		err = errors.New("invalid socket path")
		return socketType, socketAddr, err
	}

	socketType = parts[0]
	socketAddr = strings.Join(parts[1:], ":")
	switch socketType {
	case "unix", "tcp", "tcp4", "tcp6":
		break
	default:
		err = errors.New("unsupported socket type")
	}
	return socketType, socketAddr, err
}
