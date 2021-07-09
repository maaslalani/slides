// Package meta implements markdown frontmatter parsing for simple
// slides configuration
package meta

import (
	"regexp"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v2"
)

var (
	comment = regexp.MustCompile(`\s*#.*`)
)

// Meta contains all of the data to be parsed
// out of a markdown file's header section
type Meta struct {
	Theme string `yaml:"theme"`
}

// New creates a new instance of the
// slideshow meta header object
func New() *Meta {
	return &Meta{}
}

// ParseHeader parses metadata from a slideshows header slide
// including theme information
//
// If no front matter is provided, it will fallback to the default theme and
// return false to acknowledge that there is no front matter in this slide
func (m *Meta) ParseHeader(header string) (*Meta, bool) {
	fallback := &Meta{Theme: "default"}
	header = comment.ReplaceAllString(header, "")
	if len(header) <= 0 {
		return fallback, false
	}
	bytes, err := frontmatter.Parse(strings.NewReader(header), &m)
	if err != nil {
		return fallback, false
	}

	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		return fallback, false
	}

	return m, true
}
