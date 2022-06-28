// Package keygen handles the creation of new SSH key pairs.
package keygen

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/caarlos0/sshmarshal"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

// KeyType represents a type of SSH key.
type KeyType string

// Supported key types.
const (
	RSA     KeyType = "rsa"
	Ed25519 KeyType = "ed25519"
	ECDSA   KeyType = "ecdsa"
)

const rsaDefaultBits = 4096

// ErrMissingSSHKeys indicates we're missing some keys that we expected to
// have after generating. This should be an extreme edge case.
var ErrMissingSSHKeys = errors.New("missing one or more keys; did something happen to them after they were generated?")

// ErrUnsupportedKeyType indicates an unsupported key type.
type ErrUnsupportedKeyType struct {
	keyType string
}

// Error implements the error interface for ErrUnsupportedKeyType
func (e ErrUnsupportedKeyType) Error() string {
	err := "unsupported key type"
	if e.keyType != "" {
		err += fmt.Sprintf(": %s", e.keyType)
	}
	return err
}

// FilesystemErr is used to signal there was a problem creating keys at the
// filesystem-level. For example, when we're unable to create a directory to
// store new SSH keys in.
type FilesystemErr struct {
	Err error
}

// Error returns a human-readable string for the error. It implements the error
// interface.
func (e FilesystemErr) Error() string {
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e FilesystemErr) Unwrap() error {
	return e.Err
}

// SSHKeysAlreadyExistErr indicates that files already exist at the location at
// which we're attempting to create SSH keys.
type SSHKeysAlreadyExistErr struct {
	Path string
}

// SSHKeyPair holds a pair of SSH keys and associated methods.
type SSHKeyPair struct {
	path       string // private key filename path; public key will have .pub appended
	passphrase []byte
	keyType    KeyType
	privateKey crypto.PrivateKey
}

func (s SSHKeyPair) privateKeyPath() string {
	p := fmt.Sprintf("%s_%s", s.path, s.keyType)
	return p
}

func (s SSHKeyPair) publicKeyPath() string {
	return s.privateKeyPath() + ".pub"
}

// New generates an SSHKeyPair, which contains a pair of SSH keys.
func New(path string, passphrase []byte, keyType KeyType) (*SSHKeyPair, error) {
	var err error
	s := &SSHKeyPair{
		path:       path,
		keyType:    keyType,
		passphrase: passphrase,
	}
	if s.KeyPairExists() {
		privData, err := ioutil.ReadFile(s.privateKeyPath())
		if err != nil {
			return nil, err
		}
		var k interface{}
		if len(passphrase) > 0 {
			k, err = ssh.ParseRawPrivateKeyWithPassphrase(privData, passphrase)
		} else {
			k, err = ssh.ParseRawPrivateKey(privData)
		}
		if err != nil {
			return nil, err
		}
		switch k := k.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, *ed25519.PrivateKey:
			s.privateKey = k
		default:
			return nil, ErrUnsupportedKeyType{fmt.Sprintf("%T", k)}
		}
		return s, nil
	}
	switch keyType {
	case Ed25519:
		err = s.generateEd25519Keys()
	case RSA:
		err = s.generateRSAKeys(rsaDefaultBits)
	case ECDSA:
		err = s.generateECDSAKeys(elliptic.P384())
	default:
		return nil, ErrUnsupportedKeyType{string(keyType)}
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// NewWithWrite generates an SSHKeyPair and writes it to disk if not exist.
func NewWithWrite(path string, passphrase []byte, keyType KeyType) (*SSHKeyPair, error) {
	s, err := New(path, passphrase, keyType)
	if err != nil {
		return nil, err
	}
	if !s.KeyPairExists() {
		if err = s.WriteKeys(); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// PrivateKey returns the unencrypted private key.
func (s *SSHKeyPair) PrivateKey() crypto.PrivateKey {
	switch s.keyType {
	case RSA, Ed25519, ECDSA:
		return s.privateKey
	default:
		return nil
	}
}

// PrivateKeyPEM returns the unencrypted private key in OPENSSH PEM format.
func (s *SSHKeyPair) PrivateKeyPEM() []byte {
	block, err := s.pemBlock(nil)
	if err != nil {
		return nil
	}
	return pem.EncodeToMemory(block)
}

// PublicKey returns the SSH public key (RFC 4253). Ready to be used in an
// OpenSSH authorized_keys file.
func (s *SSHKeyPair) PublicKey() []byte {
	var pk crypto.PublicKey
	// Prepare public key
	switch s.keyType {
	case RSA:
		key, ok := s.privateKey.(*rsa.PrivateKey)
		if !ok {
			return nil
		}
		pk = key.Public()
	case Ed25519:
		key, ok := s.privateKey.(*ed25519.PrivateKey)
		if !ok {
			return nil
		}
		pk = key.Public()
	case ECDSA:
		key, ok := s.privateKey.(*ecdsa.PrivateKey)
		if !ok {
			return nil
		}
		pk = key.Public()
	default:
		return nil
	}
	p, err := ssh.NewPublicKey(pk)
	if err != nil {
		return nil
	}
	// serialize public key
	ak := ssh.MarshalAuthorizedKey(p)
	return pubKeyWithMemo(ak)
}

func (s *SSHKeyPair) pemBlock(passphrase []byte) (*pem.Block, error) {
	key := s.PrivateKey()
	if key == nil {
		return nil, ErrMissingSSHKeys
	}
	switch s.keyType {
	case RSA, Ed25519, ECDSA:
		if len(passphrase) > 0 {
			return sshmarshal.MarshalPrivateKeyWithPassphrase(key, "", passphrase)
		}
		return sshmarshal.MarshalPrivateKey(key, "")
	default:
		return nil, ErrUnsupportedKeyType{string(s.keyType)}
	}
}

// generateEd25519Keys creates a pair of EdD25519 keys for SSH auth.
func (s *SSHKeyPair) generateEd25519Keys() error {
	// Generate keys
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	s.privateKey = &privateKey

	return nil
}

// generateEd25519Keys creates a pair of EdD25519 keys for SSH auth.
func (s *SSHKeyPair) generateECDSAKeys(curve elliptic.Curve) error {
	// Generate keys
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}
	s.privateKey = privateKey
	return nil
}

// generateRSAKeys creates a pair for RSA keys for SSH auth.
func (s *SSHKeyPair) generateRSAKeys(bitSize int) error {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return err
	}
	// Validate private key
	err = privateKey.Validate()
	if err != nil {
		return err
	}
	s.privateKey = privateKey
	return nil
}

// prepFilesystem makes sure the state of the filesystem is as it needs to be
// in order to write our keys to disk. It will create and/or set permissions on
// the SSH directory we're going to write our keys to (for example, ~/.ssh) as
// well as make sure that no files exist at the location in which we're going
// to write out keys.
func (s *SSHKeyPair) prepFilesystem() error {
	var err error

	keyDir := filepath.Dir(s.path)
	if keyDir != "" {
		keyDir, err = homedir.Expand(keyDir)
		if err != nil {
			return err
		}

		info, err := os.Stat(keyDir)
		if os.IsNotExist(err) {
			// Directory doesn't exist: create it
			return os.MkdirAll(keyDir, 0700)
		}
		if err != nil {
			// There was another error statting the directory; something is awry
			return FilesystemErr{Err: err}
		}
		if !info.IsDir() {
			// It exists but it's not a directory
			return FilesystemErr{Err: fmt.Errorf("%s is not a directory", keyDir)}
		}
		if info.Mode().Perm() != 0700 {
			// Permissions are wrong: fix 'em
			if err := os.Chmod(keyDir, 0700); err != nil {
				return FilesystemErr{Err: err}
			}
		}
	}

	// Make sure the files we're going to write to don't already exist
	if fileExists(s.privateKeyPath()) {
		return SSHKeysAlreadyExistErr{Path: s.privateKeyPath()}
	}
	if fileExists(s.publicKeyPath()) {
		return SSHKeysAlreadyExistErr{Path: s.publicKeyPath()}
	}

	// The directory looks good as-is
	return nil
}

// WriteKeys writes the SSH key pair to disk.
func (s *SSHKeyPair) WriteKeys() error {
	var err error
	priv := s.PrivateKeyPEM()
	pub := s.PublicKey()
	if priv == nil || pub == nil {
		return ErrMissingSSHKeys
	}

	// Encrypt private key with passphrase
	if len(s.passphrase) > 0 {
		block, err := s.pemBlock(s.passphrase)
		if err != nil {
			return err
		}
		priv = pem.EncodeToMemory(block)
	}
	if err = s.prepFilesystem(); err != nil {
		return err
	}

	if err := writeKeyToFile(priv, s.privateKeyPath()); err != nil {
		return err
	}
	if err := writeKeyToFile(pub, s.publicKeyPath()); err != nil {
		return err
	}

	return nil
}

// KeyPairExists checks if the SSH key pair exists on disk.
func (s *SSHKeyPair) KeyPairExists() bool {
	return fileExists(s.privateKeyPath()) && fileExists(s.publicKeyPath())
}

func writeKeyToFile(keyBytes []byte, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ioutil.WriteFile(path, keyBytes, 0600)
	}
	return FilesystemErr{Err: fmt.Errorf("file %s already exists", path)}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return true
}

// attaches a user@host suffix to a serialized public key. returns the original
// pubkey if we can't get the username or host.
func pubKeyWithMemo(pubKey []byte) []byte {
	u, err := user.Current()
	if err != nil {
		return pubKey
	}
	hostname, err := os.Hostname()
	if err != nil {
		return pubKey
	}

	return append(bytes.TrimRight(pubKey, "\n"), []byte(fmt.Sprintf(" %s@%s\n", u.Username, hostname))...)
}

// Error returns the a human-readable error message for SSHKeysAlreadyExistErr.
// It satisfies the error interface.
func (e SSHKeysAlreadyExistErr) Error() string {
	return fmt.Sprintf("ssh key %s already exists", e.Path)
}
