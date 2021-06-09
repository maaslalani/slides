package styles

import (
	_ "embed"
	"io/ioutil"
	"os"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

const (
	salmon = lg.Color("#E8B4BC")
)

var (
	Author = lg.NewStyle().Foreground(salmon).Align(lg.Left).MarginLeft(2)
	Date   = lg.NewStyle().Faint(true).Align(lg.Left).Margin(0, 1)
	Page   = lg.NewStyle().Foreground(salmon).Align(lg.Right).MarginRight(3)
	Slide  = lg.NewStyle().Padding(1)
	Status = lg.NewStyle().Padding(1)
)

func JoinHorizontal(left, right string, width int) string {
	length := lg.Width(left + right)
	if width < length {
		return left + " " + right
	}
	padding := strings.Repeat(" ", width-length)
	return left + padding + right
}

func JoinVertical(top, bottom string, height int) string {
	h := lg.Height(top) + lg.Height(bottom)
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
