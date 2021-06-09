package styles

import (
	_ "embed"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	salmon = lipgloss.Color("#E8B4BC")
)

var (
	Author = lipgloss.NewStyle().Foreground(salmon).Align(lipgloss.Left).MarginLeft(2)
	Date   = lipgloss.NewStyle().Faint(true).Align(lipgloss.Left).Margin(0, 1)
	Page   = lipgloss.NewStyle().Foreground(salmon).Align(lipgloss.Right).MarginRight(3)
	Slide  = lipgloss.NewStyle().Padding(1)
	Status = lipgloss.NewStyle().Padding(1)
)

var (
	//go:embed theme.json
	DefaultTheme []byte

	//go:embed theme_dark.json
	DarkTheme []byte

	//go:embed theme_light.json
	LightTheme []byte
)

func JoinHorizontal(left, right string, width int) string {
	length := lipgloss.Width(left + right)
	if width < length {
		return left + " " + right
	}
	padding := strings.Repeat(" ", width-length)
	return left + padding + right
}

func JoinVertical(top, bottom string, height int) string {
	h := lipgloss.Height(top) + lipgloss.Height(bottom)
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

func SelectTheme(theme string) ([]byte, error) {
	switch theme {
	case "default":
		return DefaultTheme, nil
	case "dark":
		return DarkTheme, nil
	case "light":
		return LightTheme, nil
	default:
		return nil, errors.New("could not apply custom theme")
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
