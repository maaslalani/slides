package meta

import (
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v2"
)

type Meta struct {
	Theme string `yaml:"theme"`
}

// New creates a new instance of the
// slideshow meta object
func New() *Meta {
	return &Meta{}
}

// ParseHeader parses metadata from a slideshows header slide
// including theme information
func (m *Meta) ParseHeader(header string) (*Meta, error) {
	bytes, err := frontmatter.Parse(strings.NewReader(header), &m)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
