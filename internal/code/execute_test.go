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
				Out:      "exit status 127",
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
