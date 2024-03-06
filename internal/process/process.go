package process

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"github.com/kballard/go-shellquote"
)

// Block represents a pre-processable block which looks like the following: It
// is delimited by ~~~ and contains a command to be run along with the input to
// be passed, the entire block should be replaced with its command output
//
// ~~~sd block process
// block
// ~~~
type Block struct {
	Command string
	Input   string
	Output  string
	Raw     string
}

// String implements the Stringer interface.
func (b Block) String() string {
	return fmt.Sprintf("===\n%s\n%s\n%s\n===", b.Raw, b.Command, b.Input)
}

// ?: means non-capture group
var reng = regexp.MustCompile("~~~(.+)\n(?:.|\n)*?\n~~~\\s?")
var reg = regexp.MustCompile("(?s)~~~(.+?)\n(.*?)\n~~~\\s?")

// Parse takes some markdown and returns blocks to be pre-processed
func Parse(markdown string) []Block {
	var blocks []Block
	matches := reng.FindAllString(markdown, -1)
	for _, match := range matches {
		m := reg.FindStringSubmatch(match)
		blocks = append(blocks, Block{
			Command: m[1],
			Input:   m[2],
			Raw:     strings.TrimSuffix(m[0], "\n"),
		})
	}
	return blocks
}

// Execute takes performs the execution of the block's command
// by passing in the block's input as stdin and sets the block output
func (b *Block) Execute() {
	c, _ := shellquote.Split(b.Command)
	cmd := exec.Command(c[0], c[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	go func() {
		defer stdin.Close()
		_, _ = io.WriteString(stdin, b.Input)
	}()

	out, err := cmd.Output()
	if err != nil {
		return
	}

	b.Output = string(out)
}

// Pre processes the markdown content by executing the commands necessary and
// returns the new processed content
func Pre(content string) string {
	blocks := Parse(content)

	if len(blocks) <= 0 {
		return content
	}

	for _, block := range blocks {
		// TODO: Use goroutines, if possible
		block.Execute()

		// If multiple blocks have the same Raw value The will _likely_ have the
		// same Output value so we can probably optimize this
		// There may be edge cases, though, since block execution is not deterministic.
		content = strings.Replace(content, block.Raw, block.Output, 1)
	}
	return content
}
