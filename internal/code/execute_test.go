package code_test

import (
	"testing"

	"github.com/maaslalani/slides/internal/code"
)

func TestExecute(t *testing.T) {
	tt := []struct {
		block    code.Block
		expected code.Result
	}{
		{
			block: code.Block{
				Code:     `fn main() { println!("Hello, world!"); }`,
				Language: code.Rust,
			},
			expected: code.Result{
				Out:      "Hello, world!\n",
				ExitCode: 0,
			},
		},
		{
			block: code.Block{
				Code:     `puts "Hello, world!"`,
				Language: "ruby",
			},
			expected: code.Result{
				Out:      "Hello, world!\n",
				ExitCode: 0,
			},
		},
		{
			block: code.Block{
				Code:     `puts "Hi, there!"`,
				Language: "ruby",
			},
			expected: code.Result{
				Out:      "Hi, there!\n",
				ExitCode: 0,
			},
		},
		{
			block: code.Block{
				Code:     `print "No new line"`,
				Language: "ruby",
			},
			expected: code.Result{
				Out:      "No new line",
				ExitCode: 0,
			},
		},
		{
			block: code.Block{
				Code: `
package main

import "fmt"

func main() {
  fmt.Print("Hello, go!")
}
        `,
				Language: "go",
			},
			expected: code.Result{
				Out:      "Hello, go!",
				ExitCode: 0,
			},
		},
		{
			block: code.Block{
				Code:     `print("Hello, python!")`,
				Language: "python",
			},
			expected: code.Result{
				Out:      "Hello, python!\n",
				ExitCode: 0,
			},
		},
		{
			block: code.Block{
				Code:     `echo "Hello, bash!"`,
				Language: "bash",
			},
			expected: code.Result{
				Out:      "Hello, bash!\n",
				ExitCode: 0,
			},
		},
		{
			block: code.Block{
				Code:     `Invalid Code`,
				Language: "bash",
			},
			expected: code.Result{
				Out:      "",
				ExitCode: 1,
			},
		},
		{
			block: code.Block{
				Code:     `Invalid Code`,
				Language: "invalid",
			},
			expected: code.Result{
				Out:      "Error: unsupported language",
				ExitCode: code.ExitCodeInternalError,
			},
		},
	}

	for _, tc := range tt {
		if testing.Short() {
			t.SkipNow()
		}
		r := code.Execute(tc.block)
		if r.Out != tc.expected.Out {
			t.Fatalf("invalid output for lang %s, got %s, want %s | %+v",
				tc.block.Language, r.Out, tc.expected.Out, r)
		}

		if r.ExitCode != tc.expected.ExitCode {
			t.Fatalf("unexpected exit code, got %d, want %d", r.ExitCode, tc.expected.ExitCode)
		}
	}
}
