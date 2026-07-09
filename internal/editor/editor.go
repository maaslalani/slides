package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/coral"
)

var program = func() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "vim"
}()

// Opens the current slide as a split window in tmux.
func OpenNewWindow(fileName string, slide string) tea.Cmd {
	var c *exec.Cmd

	editorName := GetEditorName(program)

	switch editorName {
	case "vim", "nvim":
		c = exec.Command(program, fmt.Sprintf("+%d", GetLineNumber(fileName, slide)), fileName)
	case "code":
		c = exec.Command(program, fileName, "--go-to", fmt.Sprintf("+%d", GetLineNumber(fileName, slide)))
	default:
		c = exec.Command(program, fileName)
	}

	return tea.ExecProcess(c, nil)
}

func linesMatch(lines []string, sLines []string, start int) bool {
	for j := range sLines {
		if strings.TrimSpace(lines[start+j]) != strings.TrimSpace(sLines[j]) {
			return false
		}
	}
	return true
}

func GetLineNumber(fileName string, slide string) int {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return -1
	}

	lines := strings.Split(string(b), "\n")
	sLines := strings.Split(slide, "\n")

	for i := range lines {
		if linesMatch(lines, sLines, i) {
			return i
		}
	}

	return -1
}

func GetEditorName(editorPath string) string {
	parts := strings.Split(editorPath, "/")
	return parts[len(parts)-1]
}

func InitEditorFlag(rootCmd *coral.Command) {
	rootCmd.PersistentFlags().StringVarP(&program, "editor", "e", program, "Specify the editor to use")
}
