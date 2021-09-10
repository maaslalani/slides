package code_test

import (
	"testing"

	"github.com/maaslalani/slides/internal/code"
)

func TestParse(t *testing.T) {
	tt := []struct {
		markdown string
		expected []code.Block
	}{
		// We can't put backticks ```
		// in multi-line strings, ~~~ instead
		{
			markdown: `
~~~ruby
puts "Hello, world!"
~~~
`,
			expected: []code.Block{
				{
					Code:     `puts "Hello, world!"`,
					Language: "ruby",
				},
			},
		},
		{
			markdown: `
~~~go
fmt.Println("Hello, world!")
~~~
`,
			expected: []code.Block{
				{
					Code:     `fmt.Println("Hello, world!")`,
					Language: "go",
				},
			},
		},
		{
			markdown: `
~~~python
print("Hello, world!")
~~~`,
			expected: []code.Block{
				{
					Code:     `print("Hello, world!")`,
					Language: "python",
				},
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
			expected: []code.Block{
				{
					Code: `package main

import "fmt"

func main() {
  fmt.Println("Written in Go!")
}`,
					Language: "go",
				},
			},
		},
		{
			markdown: `
# Slide 1
Just a regular slide, no code block
`,
			expected: nil,
		},
		{
			markdown: ``,
			expected: nil,
		},
		{
			markdown: `
~~~ruby
puts "Hello, world!"
~~~

~~~go
fmt.Println("Hello, world!")
~~~
`,
			expected: []code.Block{
				{
					Code:     `puts "Hello, world!"`,
					Language: "ruby",
				},
				{
					Code:     `fmt.Println("Hello, world!")`,
					Language: "go",
				},
			},
		},
	}

	for _, tc := range tt {
		blocks, _ := code.Parse(tc.markdown)
		if len(blocks) != len(tc.expected) {
			t.Errorf("parse fail: incorrect size of blocks")
		}
		for i, block := range blocks {
			expected := tc.expected[i]
			if block.Code != expected.Code {
				t.Log(block.Code)
				t.Log(expected.Code)
				t.Fatal("parse failed: incorrect code")
			}
			if block.Language != expected.Language {
				t.Fatalf("incorrect language, got %s, want %s", block.Language, expected.Language)
			}
		}
	}
}
