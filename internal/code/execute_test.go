package code_test

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/maaslalani/slides/internal/code"
)

func TestExecute(t *testing.T) {
	tt := map[string]struct {
		block    code.Block
		expected code.Result
	}{
		"go": {
			block: code.Block{
				Code: heredoc.Doc(`
					package main

					import "fmt"

					func main() {
						fmt.Print("Hello, go!")
					}`),
				Language: "go",
			},
			expected: code.Result{
				Out:      "Hello, go!",
				ExitCode: 0,
			},
		},
		"go error": {
			block: code.Block{
				Code: heredoc.Doc(`
					package main

					import "fmt"

					func main() {
  						mt.Print("Hello, go!")
					}`),
				Language: "go",
			},
			expected: code.Result{
				Out: heredoc.Doc(`
					# command-line-arguments
					imported and not used: "fmt"
					undefined: mt
					exit status 2`),
				ExitCode: 2,
			},
		},
		"bash": {
			block: code.Block{
				Code:     `echo "Hello, bash!"`,
				Language: "bash",
			},
			expected: code.Result{
				Out:      "Hello, bash!\n",
				ExitCode: 0,
			},
		},
		"bash invalid code": {
			block: code.Block{
				Code:     `Invalid Code`,
				Language: "bash",
			},
			expected: code.Result{
				Out: heredoc.Doc(`
					Invalid: command not found
					exit status 127`),
				ExitCode: 127,
			},
		},
		"invalid language": {
			block: code.Block{
				Code:     `Invalid Language`,
				Language: "invalid",
			},
			expected: code.Result{
				Out:      "Error: unsupported language",
				ExitCode: code.ExitCodeInternalError,
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			r := code.Execute(tc.block)
			if r.Out != tc.expected.Out {
				t.Fatalf("invalid output for lang:%s result:\n%+v\n\ngot:\n%s\nwant:\n%s", tc.block.Language, r, r.Out, tc.expected.Out)
			}

			if r.ExitCode != tc.expected.ExitCode {
				t.Fatalf("unexpected exit code, got %d, want %d", r.ExitCode, tc.expected.ExitCode)
			}
		})
	}
}
