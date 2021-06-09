package styles

import (
	_ "embed"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/charmbracelet/lipgloss"
)

const (
	salmon = Color("#E8B4BC")
)

var (
	Author = NewStyle().Foreground(salmon).Align(Left).MarginLeft(2)
	Date   = NewStyle().Faint(true).Align(Left).Margin(0, 1)
	Page   = NewStyle().Foreground(salmon).Align(Right).MarginRight(3)
	Slide  = NewStyle().Padding(1)
	Status = NewStyle().Padding(1)
)

func JoinHorizontal(left, right string, width int) string {
	length := Width(left + right)
	if width < length {
		return left + " " + right
	}
	padding := strings.Repeat(" ", width-length)
	return left + padding + right
}

func JoinVertical(top, bottom string, height int) string {
	h := Height(top) + Height(bottom)
	if height < h {
		return top + "\n" + bottom
	}
	fill := strings.Repeat("\n", height-h)
	return top + fill + bottom
}

// CustomTheme reads a json theme file from the entered
// path and returns the bytes
func CustomTheme(filename string) ([]byte, error) {
	if fileExists(filename) {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		return b, err
	}
	return nil, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//go:embed theme.json
var Theme []byte
