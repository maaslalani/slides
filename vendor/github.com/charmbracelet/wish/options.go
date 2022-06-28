package wish

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/keygen"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// WithAddress returns an ssh.Option that sets the address to listen on.
func WithAddress(addr string) ssh.Option {
	return func(s *ssh.Server) error {
		s.Addr = addr
		return nil
	}
}

// WithVersion returns an ssh.Option that sets the server version.
func WithVersion(version string) ssh.Option {
	return func(s *ssh.Server) error {
		s.Version = version
		return nil
	}
}

// WithMiddleware composes the provided Middleware and return a ssh.Option.
// This useful if you manually create an ssh.Server and want to set the
// Server.Handler.
//
// Notice that middlewares are composed from first to last, which means the last one is executed first.
func WithMiddleware(mw ...Middleware) ssh.Option {
	return func(s *ssh.Server) error {
		h := func(s ssh.Session) {}
		for _, m := range mw {
			h = m(h)
		}
		s.Handler = h
		return nil
	}
}

// WithHostKeyFile returns an ssh.Option that sets the path to the private.
func WithHostKeyPath(path string) ssh.Option {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir, f := filepath.Split(path)
		n := strings.TrimSuffix(f, "_ed25519")
		_, err := keygen.NewWithWrite(filepath.Join(dir, n), nil, keygen.Ed25519)
		if err != nil {
			return func(*ssh.Server) error {
				return err
			}
		}
		path = filepath.Join(dir, n+"_ed25519")
	}
	return ssh.HostKeyFile(path)
}

// WithHostKeyPEM returns an ssh.Option that sets the host key from a PEM block.
func WithHostKeyPEM(pem []byte) ssh.Option {
	return ssh.HostKeyPEM(pem)
}

// WithAuthorizedKeys allows to use a SSH authorized_keys file to allowlist users.
func WithAuthorizedKeys(path string) ssh.Option {
	return func(s *ssh.Server) error {
		keys, err := parseAuthorizedKeys(path)
		if err != nil {
			return err
		}
		return WithPublicKeyAuth(func(_ ssh.Context, key ssh.PublicKey) bool {
			for _, upk := range keys {
				if ssh.KeysEqual(upk, key) {
					return true
				}
			}
			return false
		})(s)
	}
}

// WithTrustedUserCAKeys authorize certificates that are signed with the given
// Certificate Authority public key, and are valid.
// Analogous to the TrustedUserCAKeys OpenSSH option.
func WithTrustedUserCAKeys(path string) ssh.Option {
	return func(s *ssh.Server) error {
		cas, err := parseAuthorizedKeys(path)
		if err != nil {
			return err
		}

		return WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			cert, ok := key.(*gossh.Certificate)
			if !ok {
				// not a certificate...
				return false
			}

			checker := &gossh.CertChecker{
				IsUserAuthority: func(auth gossh.PublicKey) bool {
					for _, ca := range cas {
						if bytes.Equal(auth.Marshal(), ca.Marshal()) {
							// its a cert signed by one of the CAs
							return true
						}
					}
					// it is a cert, but signed by another CA
					return false
				},
			}

			if !checker.IsUserAuthority(cert.SignatureKey) {
				return false
			}

			if err := checker.CheckCert(ctx.User(), cert); err != nil {
				return false
			}

			return true
		})(s)
	}
}

func parseAuthorizedKeys(path string) ([]ssh.PublicKey, error) {
	var keys []ssh.PublicKey

	f, err := os.Open(path)
	if err != nil {
		return keys, fmt.Errorf("failed to parse %q: %w", path, err)
	}
	defer f.Close() // nolint: errcheck

	rd := bufio.NewReader(f)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return keys, fmt.Errorf("failed to parse %q: %w", path, err)
		}
		upk, _, _, _, err := ssh.ParseAuthorizedKey(line)
		if err != nil {
			return keys, fmt.Errorf("failed to parse %q: %w", path, err)
		}
		keys = append(keys, upk)
	}
	return keys, nil
}

// WithPublicKeyAuth returns an ssh.Option that sets the public key auth handler.
func WithPublicKeyAuth(h ssh.PublicKeyHandler) ssh.Option {
	return ssh.PublicKeyAuth(h)
}

// WithPasswordAuth returns an ssh.Option that sets the password auth handler.
func WithPasswordAuth(p ssh.PasswordHandler) ssh.Option {
	return ssh.PasswordAuth(p)
}

// WithIdleTimeout returns an ssh.Option that sets the connection's idle timeout.
func WithIdleTimeout(d time.Duration) ssh.Option {
	return func(s *ssh.Server) error {
		s.IdleTimeout = d
		return nil
	}
}

// WithMaxTimeout returns an ssh.Option that sets the connection's absolute timeout.
func WithMaxTimeout(d time.Duration) ssh.Option {
	return func(s *ssh.Server) error {
		s.MaxTimeout = d
		return nil
	}
}
