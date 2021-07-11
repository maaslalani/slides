// Package meta implements markdown frontmatter parsing for simple
// slides configuration
package meta

import (
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

// Parse parses metadata from a slideshows header slide
// including theme information
//
// If no front matter is provided, it will fallback to the default theme and
// return false to acknowledge that there is no front matter in this slide
func (m *Meta) Parse(header string) (*Meta, bool) {
	fallback := &Meta{
		Theme: "default",
	}

	err := yaml.Unmarshal([]byte(header), &m)
	if err != nil {
		return fallback, false
	}

	// This fixes a bug where the first slide of a presentation won't show up if
	// the first slide is valid YAML (i.e. "# Header")
	// FIXME: This only works because we currently only have one option (theme),
	if m.Theme == "" {
		return fallback, false
	}

	return m, true
}
