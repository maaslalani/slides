package styles

import (
	_ "embed"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
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

// CustomTheme reads in a custom glamour theme configuration
// from a filepath
func CustomTheme(filepath string) glamour.TermRendererOption {
	if fileExists(filepath) {
		return glamour.WithStylesFromJSONFile(filepath)
	}
	if filepath == "default" {
		return glamour.WithStylesFromJSONBytes(DefaultTheme)
	}
	return nil
}

// SelectTheme picks a glamour style config based
// on the theme provided in the markdown header
func SelectTheme(theme string) ansi.StyleConfig {
	switch theme {
	case "ascii":
		return glamour.ASCIIStyleConfig
	case "light":
		return glamour.LightStyleConfig
	case "notty":
		return glamour.NoTTYStyleConfig
	default:
		return glamour.DarkStyleConfig
	}
}

// fileExists is a helper to verify
// that the provided filepath exists
// on the system
func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
