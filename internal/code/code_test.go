package code_test

import (
	"testing"

	"github.com/maaslalani/slides/internal/code"
)

func TestParse(t *testing.T) {
	tt := []struct {
		markdown  string
		expecteds []code.Block
	}{
		// We can't put backticks ```
		// in multi-line strings, ~~~ instead
		{
			markdown: `
~~~ruby
puts "Hello, world!"
~~~
`,
			expecteds: []code.Block{
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
			expecteds: []code.Block{
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
			expecteds: []code.Block{
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
			expecteds: []code.Block{
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
			expecteds: nil,
		},
		{
			markdown:  ``,
			expecteds: nil,
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
			expecteds: []code.Block{
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
		bs, _ := code.Parse(tc.markdown)
		if len(bs) != len(tc.expecteds) {
			t.Errorf("parse fail: incorrect size of blocks")
		}
		for idx, b := range bs {
			expected := tc.expecteds[idx]
			if b.Code != expected.Code {
				t.Log(b.Code)
				t.Log(expected.Code)
				t.Fatal("parse failed: incorrect code")
			}
			if b.Language != expected.Language {
				t.Fatalf("incorrect language, got %s, want %s", b.Language, expected.Language)
			}
		}
	}
}
