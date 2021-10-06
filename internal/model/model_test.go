package model_test

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maaslalani/slides/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func keyPresses(runes string) []tea.KeyMsg {
	var result []tea.KeyMsg

	for _, r := range runes {
		result = append(result, tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{r},
			Alt:   false,
		})
	}

	return result
}

func TestModel_UpdatePageNavigation(t *testing.T) {
	tests := []struct {
		name string
		keys []tea.KeyMsg
		want int
	}{
		{
			name: "Initial state",
			keys: nil,
			want: 0,
		},
		{
			name: "Can paginate",
			keys: keyPresses("l"),
			want: 1,
		},
		{
			name: "Can paginate to end",
			keys: keyPresses("jjjjjjjjjj"),
			want: 10,
		},
		{
			name: "Cannot paginate past end",
			keys: keyPresses("jjjjjjjjjjjjj"),
			want: 10,
		},
		{
			name: "Can move to end",
			keys: keyPresses("G"),
			want: 10,
		},
		{
			name: "Can move to start",
			keys: keyPresses("llgg"),
			want: 0,
		},
		{
			name: "Repeats",
			keys: keyPresses("2j"),
			want: 2,
		},
		{
			name: "Repeats with 0 (Vim ignores 0 and just does next command, so let's do the same)",
			keys: keyPresses("0j"),
			want: 1,
		},
		{
			name: "Direct slide addressing sub min (Vim ignores minus, so let's do the same)",
			keys: keyPresses("-11G"),
			want: 10,
		},
		{
			name: "Direct slide addressing min",
			keys: keyPresses("0G"),
			want: 0,
		},
		{
			name: "Direct slide addressing in range",
			keys: keyPresses("3G"),
			want: 2,
		},
		{
			name: "Direct slide addressing max",
			keys: keyPresses("11G"),
			want: 10,
		},
		{
			name: "Direct slide addressing supra max",
			keys: keyPresses("101G"),
			want: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m tea.Model = model.Model{
				Page: 0,
				Slides: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			}

			for _, key := range tt.keys {
				m, _ = m.Update(key)
			}

			got, ok := m.(model.Model)
			assert.True(t, ok, "model must be of correct type")

			assert.Equal(t, tt.want, got.Page)
		})
	}
}
