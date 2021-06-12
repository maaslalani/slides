package file_test

import (
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
