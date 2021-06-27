// Package file includes utility functions
// for working with the filesystem
package file

import (
	"io/fs"
	"os"
)

// Exists is a helper to verify
// that the provided filepath exists
// on the system
func Exists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsExecutable returns whether a file has execution permissions
func IsExecutable(s fs.FileInfo) bool {
	return s.Mode().Perm()&0111 == 0111
}
