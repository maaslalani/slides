package code

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Block represents a code block.
type Block struct {
	Code     string
	Language string
}

// Result represents the output for an executed code block.
type Result struct {
	Out           string
	ExitCode      int
	ExecutionTime time.Duration
}

// ?: means non-capture group
var re = regexp.MustCompile("(?s)(?:```|~~~)(\\w+)\n(.*?)\n(?:```|~~~)\\s?")

var (
	// ErrParse is the returned error when we cannot parse the code block (i.e.
	// there is no code block on the current slide) or the code block is
	// incorrectly written.
	ErrParse = errors.New("Error: could not parse code block")
)

// Parse takes a block of markdown and returns an array of Block's with code
// and associated languages
func Parse(markdown string) ([]Block, error) {
	matches := re.FindAllStringSubmatch(markdown, -1)

	var rv []Block
	for _, match := range matches {
		// There was either no language specified or no code block
		// Either way, we cannot execute the expression
		if len(match) < 3 {
			continue
		}
		rv = append(rv, Block{
			Language: match[1],
			Code:     RemoveComments(match[2]),
		})

	}

	if len(rv) == 0 {
		return nil, ErrParse
	}

	return rv, nil
}

const (
	// ExitCodeInternalError represents the exit code in which the code
	// executing the code didn't work.
	ExitCodeInternalError = -1
)

// Execute takes a code.Block and returns the output of the executed code
func Execute(code Block) Result {
	// Check supported language
	language, ok := Languages[code.Language]
	if !ok {
		return Result{
			Out:      "Error: unsupported language",
			ExitCode: ExitCodeInternalError,
		}
	}

	// Write the code block to a temporary file
	f, err := ioutil.TempFile(os.TempDir(), "slides-*."+Languages[code.Language].Extension)
	if err != nil {
		return Result{
			Out:      "Error: could not create file",
			ExitCode: ExitCodeInternalError,
		}
	}

	defer f.Close()
	defer os.Remove(f.Name())

	_, err = f.WriteString(code.Code)
	if err != nil {
		return Result{
			Out:      "Error: could not write to file",
			ExitCode: ExitCodeInternalError,
		}
	}

	var (
		output   strings.Builder
		exitCode int
	)

	// replacer for commands
	repl := strings.NewReplacer(
		"<file>", f.Name(),
		// <name>: file name without extension and without path
		"<name>", filepath.Base(strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))),
		"<path>", filepath.Dir(f.Name()),
	)

	// For accuracy of program execution speed, we can't put anything after
	// recording the start time or before recording the end time.
	start := time.Now()

	for _, c := range language.Commands {
		var command []string
		// replace <file>, <name> and <path> in commands
		for _, v := range c {
			command = append(command, repl.Replace(v))
		}
		// execute and write output
		cmd := exec.Command(command[0], command[1:]...)
		out, err := cmd.Output()
		if err != nil {
			output.Write([]byte(err.Error()))
		} else {
			output.Write(out)
		}

		// update status code
		if err != nil {
			if cmd.ProcessState != nil {
				exitCode = cmd.ProcessState.ExitCode()
			} else {
				exitCode = 1 // non-zero
			}
		}
	}

	end := time.Now()

	return Result{
		Out:           output.String(),
		ExitCode:      exitCode,
		ExecutionTime: end.Sub(start),
	}
}
