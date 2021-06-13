package code

import (
	"errors"
	"regexp"
)

type Block struct {
	Code     string
	Language string
}

// ?: means non-capture group
var re = regexp.MustCompile("(?:```|~~~)(.*)\n(.*)\n(?:```|~~~)")

var (
	ErrParse = errors.New("Error: could not parse code block")
)

// Parse takes a block of markdown and returns an array of Block's with code
// and associated languages
func Parse(markdown string) (Block, error) {
	match := re.FindStringSubmatch(markdown)

	// There was either no language specified or no code block
	// Either way, we cannot execute the expression
	if len(match) < 3 {
		return Block{}, ErrParse
	}

	return Block{
		Language: match[1],
		Code:     match[2],
	}, nil
}
