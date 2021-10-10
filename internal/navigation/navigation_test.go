package navigation_test

import (
	"github.com/maaslalani/slides/internal/navigation"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func keyPresses(keys string) []string {
	return strings.Split(keys, "")
}

func emptyBufferStateWithSlide(slide int) navigation.State {
	return navigation.State{
		Buffer:    "",
		Slide:     slide,
		NumSlides: 11,
	}
}

func TestModel_UpdatePageNavigation(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want navigation.State
	}{
		{
			name: "Initial state",
			keys: nil,
			want: emptyBufferStateWithSlide(0),
		},
		{
			name: "Can paginate",
			keys: keyPresses("l"),
			want: emptyBufferStateWithSlide(1),
		},
		{
			name: "Can paginate to end",
			keys: keyPresses("jjjjjjjjjj"),
			want: emptyBufferStateWithSlide(10),
		},
		{
			name: "Cannot paginate past end",
			keys: keyPresses("jjjjjjjjjjjjj"),
			want: emptyBufferStateWithSlide(10),
		},
		{
			name: "Can move to end",
			keys: keyPresses("G"),
			want: emptyBufferStateWithSlide(10),
		},
		{
			name: "Can move to start",
			keys: keyPresses("llgg"),
			want: emptyBufferStateWithSlide(0),
		},
		{
			name: "Repeats",
			keys: keyPresses("2j"),
			want: emptyBufferStateWithSlide(2),
		},
		{
			name: "Repeats with 0 (Vim ignores 0 and just does next command, so let's do the same)",
			keys: keyPresses("0j"),
			want: emptyBufferStateWithSlide(1),
		},
		{
			name: "Direct slide addressing sub min (Vim ignores minus, so let's do the same)",
			keys: keyPresses("-11G"),
			want: emptyBufferStateWithSlide(10),
		},
		{
			name: "Direct slide addressing min",
			keys: keyPresses("0G"),
			want: emptyBufferStateWithSlide(0),
		},
		{
			name: "Direct slide addressing in range",
			keys: keyPresses("3G"),
			want: emptyBufferStateWithSlide(2),
		},
		{
			name: "Direct slide addressing max",
			keys: keyPresses("11G"),
			want: emptyBufferStateWithSlide(10),
		},
		{
			name: "Direct slide addressing supra max",
			keys: keyPresses("101G"),
			want: emptyBufferStateWithSlide(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentState := navigation.State{
				Buffer:    "",
				Slide:     0,
				NumSlides: 11,
			}

			for _, key := range tt.keys {
				currentState = navigation.Navigate(currentState, key)
			}

			assert.Equal(t, tt.want, currentState)
		})
	}
}
