// Package meta implements markdown frontmatter parsing for simple
// slides configuration
package meta

import (
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"
)

// Temporary structure to differentiate values not present in the YAML header
// from values set to empty strings in the YAML header. We replace values not
// set by defaults values when parsing a header.
type parsedMeta struct {
	Theme  *string `yaml:"theme"`
	Author *string `yaml:"author"`
	Date   *string `yaml:"date"`
	Paging *string `yaml:"paging"`
}

// Meta contains all of the data to be parsed
// out of a markdown file's header section
type Meta struct {
	Theme  string
	Author string
	Date   string
	Paging string
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
		Theme:  defaultTheme(),
		Author: defaultAuthor(),
		Date:   defaultDate(),
		Paging: defaultPaging(),
	}

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

	if tmp.Author != nil {
		m.Author = *tmp.Author
	} else {
		m.Author = fallback.Author
	}

	if tmp.Date != nil {
		m.Date = parseDate(*tmp.Date)
	} else {
		m.Date = fallback.Date
	}

	if tmp.Paging != nil {
		m.Paging = *tmp.Paging
	} else {
		m.Paging = fallback.Paging
	}

	return m, true
}

func defaultTheme() string {
	return "default"
}

func defaultAuthor() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}

	return user.Name
}

func defaultDate() string {
	return "2006-01-02"
}

func defaultPaging() string {
	return "Slide %d / %d"
}

func parseDate(value string) string {
	pairs := [][]string{
		{"YYYY", "2006"},
		{"YY", "06"},
		{"MMMM", "January"},
		{"MMM", "Jan"},
		{"MM", "01"},
		{"mm", "1"},
		{"DD", "02"},
		{"dd", "2"},
	}

	for _, p := range pairs {
		value = strings.ReplaceAll(value, p[0], p[1])
	}
	return value
}
