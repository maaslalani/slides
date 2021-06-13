package code_test

import (
	"testing"

	"github.com/maaslalani/slides/internal/code"
)

func TestParse(t *testing.T) {
	tt := []struct {
		markdown string
		expected code.Block
	}{
		// We can't put backticks ```
		// in multi-line strings, ~~~ instead
		{
			markdown: `
~~~ruby
puts "Hello, world!"
~~~
`,
			expected: code.Block{
				Code:     `puts "Hello, world!"`,
				Language: "ruby",
			},
		},
		{
			markdown: `
~~~go
fmt.Println("Hello, world!")
~~~
`,
			expected: code.Block{
				Code:     `fmt.Println("Hello, world!")`,
				Language: "go",
			},
		},
		{
			markdown: `
# Welcome to Slides

A terminal based presentation tool

~~~go
package main

import "fmt"

func main() {
  fmt.Println("Written in Go!")
}
~~~
`,
			expected: code.Block{
				Code: `package main

import "fmt"

func main() {
  fmt.Println("Written in Go!")
}`,
				Language: "go",
			},
		},
		{
			markdown: `
# Slide 1
Just a regular slide, no code block
`,
			expected: code.Block{},
		},
		{
			markdown: ``,
			expected: code.Block{},
		},
	}

	for _, tc := range tt {
		b, _ := code.Parse(tc.markdown)
		if b.Code != tc.expected.Code {
			t.Log(b.Code)
			t.Log(tc.expected.Code)
			t.Fatal("parse failed: incorrect code")
		}
		if b.Language != tc.expected.Language {
			t.Fatalf("incorrect language, got %s, want %s", b.Language, tc.expected.Language)
		}
	}
}
