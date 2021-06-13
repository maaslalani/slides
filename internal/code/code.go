package code

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"time"
)

type Block struct {
	Code     string
	Language string
}

type Result struct {
	Out           string
	ExitCode      int
	ExecutionTime time.Duration
}

// ?: means non-capture group
var re = regexp.MustCompile("(?s)(?:```|~~~)(\\w+)\n(.*)\n(?:```|~~~)\n")

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
	defer os.Remove(f.Name())

	_, err = f.WriteString(code.Code)
	if err != nil {
		return Result{
			Out:      "Error: could not write to file",
			ExitCode: ExitCodeInternalError,
		}
	}

	cmd := exec.Command(language.Command[0], append(language.Command[1:], f.Name())...)

	// For accuracy of program execution speed, we can't put anything after
	// recording the start time or before recording the end time.
	start := time.Now()
	out, err := cmd.Output()
	end := time.Now()

	exitCode := 0
	if err != nil {
		exitCode = 1
	}

	return Result{
		Out:           string(out),
		ExitCode:      exitCode,
		ExecutionTime: end.Sub(start),
	}
}
