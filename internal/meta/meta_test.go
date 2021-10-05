package meta_test

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/maaslalani/slides/internal/meta"
	"github.com/stretchr/testify/assert"
)

func TestMeta_ParseHeader(t *testing.T) {
	user, _ := user.Current()
	date := "2006-01-02"

	tests := []struct {
		name      string
		slideshow string
		want      *meta.Meta
	}{
		{
			name:      "Parse theme from header",
			slideshow: fmt.Sprintf("---\ntheme: %q\n", "dark"),
			want: &meta.Meta{
				Theme:  "dark",
				Author: user.Name,
				Date:   date,
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Fallback to default if no theme provided",
			slideshow: "\n# Header Slide\n > Subtitle\n",
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   date,
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Parse author from header",
			slideshow: fmt.Sprintf("---\nauthor: %q\n", "gopher"),
			want: &meta.Meta{
				Theme:  "default",
				Author: "gopher",
				Date:   date,
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Fallback to default if no author provided",
			slideshow: "\n# Header Slide\n > Subtitle\n",
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   date,
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Parse static date from header",
			slideshow: fmt.Sprintf("---\ndate: %q\n", "31/01/1970"),
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   "31/01/1970",
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Parse go-styled date from header",
			slideshow: fmt.Sprintf("---\ndate: %q\n", "Jan 2, 2006"),
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   "Jan 2, 2006",
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Parse YYYY-MM-DD date from header",
			slideshow: fmt.Sprintf("---\ndate: %q\n", "YYYY-MM-DD"),
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   "2006-01-02",
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Parse dd/mm/YY date from header",
			slideshow: fmt.Sprintf("---\ndate: %q\n", "dd/mm/YY"),
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   "2/1/06",
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Parse MMM dd, YYYY date from header",
			slideshow: fmt.Sprintf("---\ndate: %q\n", "MMM dd, YYYY"),
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   "Jan 2, 2006",
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Parse MMMM DD, YYYY date from header",
			slideshow: fmt.Sprintf("---\ndate: %q\n", "MMMM DD, YYYY"),
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   "January 02, 2006",
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Fallback to default if no date provided",
			slideshow: "\n# Header Slide\n > Subtitle\n",
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   date,
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Parse paging from header",
			slideshow: fmt.Sprintf("---\npaging: %q\n", "%d of %d"),
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   date,
				Paging: "%d of %d",
			},
		},
		{
			name:      "Fallback to default if no numebring provided",
			slideshow: "\n# Header Slide\n > Subtitle\n",
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   date,
				Paging: "Slide %d / %d",
			},
		},
		{
			name:      "Fallback if first slide is valid yaml",
			slideshow: "---\n# Header Slide---\nContent\n",
			want: &meta.Meta{
				Theme:  "default",
				Author: user.Name,
				Date:   date,
				Paging: "Slide %d / %d",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &meta.Meta{}
			got, hasMeta := m.Parse(tt.slideshow)
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

func ExampleMeta_Parse() {
	header := `
---
theme: "dark"
author: "Gopher"
date: "Apr. 4, 2021"
paging: "%d"
---
`
	// Parse the header from the markdown
	// file
	m, _ := meta.New().Parse(header)

	// Print the return theme
	// meta
	fmt.Println(m.Theme)
	fmt.Println(m.Author)
	fmt.Println(m.Date)
	fmt.Println(m.Paging)
}
