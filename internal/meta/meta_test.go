package meta_test

import (
	"fmt"
	"testing"

	"github.com/maaslalani/slides/internal/meta"
	"github.com/stretchr/testify/assert"
)

func TestMeta_ParseHeader(t *testing.T) {
	tests := []struct {
		name      string
		slideshow string
		want      *meta.Meta
	}{
		{
			name:      "Parse theme from header",
			slideshow: fmt.Sprintf("---\ntheme: %q\n", "dark"),
			want: &meta.Meta{
				Theme: "dark",
			},
		},
		{
			name:      "Fallback to default if no theme provided",
			slideshow: "\n# Header Slide\n > Subtitle\n",
			want: &meta.Meta{
				Theme: "default",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &meta.Meta{}
			got, hasMeta := m.ParseHeader(tt.slideshow)
			if !hasMeta {
				assert.NotNil(t, got)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *meta.Meta
	}{
		{name: "Create meta struct", want: &meta.Meta{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, meta.New(), tt.want)
		})
	}
}

func ExampleMeta_ParseHeader() {
	header := `
---
theme: "dark"
---
`
	// Parse the header from the markdown
	// file
	m, _ := meta.New().ParseHeader(header)

	// Print the return theme
	// meta
	fmt.Println(m.Theme)
}
