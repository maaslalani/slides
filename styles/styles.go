package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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

func SpreadHorizontal(left, right string, width int) string {
	length := lipgloss.Width(left + right)
	if width < length {
		return ""
	}
	padding := strings.Repeat(" ", width-length)
	return left + padding + right
}
