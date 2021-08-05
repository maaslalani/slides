// Package meta implements markdown frontmatter parsing for simple
// slides configuration
package meta

import (
	"gopkg.in/yaml.v2"
)

const defaultTheme = "default"

// Temporary structure to differentiate values not present in the YAML header
// from values set to empty strings in the YAML header. We replace values not
// set by defaults values when parsing a header.
type parsedMeta struct {
	Theme *string `yaml:"theme"`
}

// Meta contains all of the data to be parsed
// out of a markdown file's header section
type Meta struct {
	Theme string
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
	fallback := &Meta{Theme: defaultTheme}

	var tmp parsedMeta
	err := yaml.Unmarshal([]byte(header), &tmp)
	if err != nil {
		return fallback, false
	}

	if tmp.Theme != nil {
		m.Theme = *tmp.Theme
	} else {
		m.Theme = fallback.Theme
	}

	return m, true
}
