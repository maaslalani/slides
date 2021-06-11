// Package meta implements markdown frontmatter parsing for simple
// slides configuration
package meta

import (
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v2"
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
// If no front matter is provided, it will fallback to the default theme
func (m *Meta) ParseHeader(header string) (*Meta, error) {
	fallback := &Meta{Theme: "default"}
	bytes, err := frontmatter.Parse(strings.NewReader(header), &m)
	if err != nil {
		return fallback, nil
	}

	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		return fallback, err
	}

	return m, nil
}
