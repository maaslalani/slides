package file_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/maaslalani/slides/internal/file"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     bool
	}{
		{name: "Find file exists", filepath: "file.go", want: true},
		{name: "Return false for missing file", filepath: "afilethatdoesntexist.go", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isExist := file.Exists(tt.filepath)
			if isExist {
				assert.FileExists(t, tt.filepath)
			}
			assert.Equal(t, tt.want, isExist)
		})
	}
}

func TestIsExecutable(t *testing.T) {
	tests := []struct {
		perm     fs.FileMode
		expected bool
	}{
		{0101, false},
		{0111, true},
		{0644, false},
		{0666, false},
		{0777, true},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprint(tc.perm), func(t *testing.T) {
			tmp, err := os.CreateTemp(os.TempDir(), "slides-*")
			if err != nil {
				t.Fatal("failed to create temp file")
			}
			defer os.Remove(tmp.Name())

			err = tmp.Chmod(tc.perm)
			if err != nil {
				t.Fatal(err)
			}

			s, err := tmp.Stat()
			if err != nil {
				t.Fatal("failed to stat file")
			}

			want := tc.expected
			got := file.IsExecutable(s)
			if tc.expected != got {
				t.Log(want)
				t.Log(got)
				t.Fatalf("IsExecutable returned an incorrect result, want: %t, got %t", want, got)
			}
		})
	}
}
