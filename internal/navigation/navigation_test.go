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

func TestModel_UpdatePageNavigation(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want navigation.State
	}{
		{
			name: "Initial state",
			keys: nil,
			want: navigation.State{
				Buffer: "",
				Slide: 0,
				ClearVirtualText: false,
			},
		},
		{
			name: "Can paginate",
			keys: keyPresses("l"),
			want: navigation.State{
				Buffer: "",
				Slide: 1,
				ClearVirtualText: true,
			},
		},
		{
			name: "Can paginate to end",
			keys: keyPresses("jjjjjjjjjj"),
			want: navigation.State{
				Buffer: "",
				Slide: 10,
				ClearVirtualText: true,
			},
		},
		{
			name: "Cannot paginate past end",
			keys: keyPresses("jjjjjjjjjjjjj"),
			want: navigation.State{
				Buffer: "",
				Slide: 10,
				ClearVirtualText: true,
			},
		},
		{
			name: "Can move to end",
			keys: keyPresses("G"),
			want: navigation.State{
				Buffer: "",
				Slide: 10,
				ClearVirtualText: false,
			},
		},
		{
			name: "Can move to start",
			keys: keyPresses("llgg"),
			want: navigation.State{
				Buffer: "",
				Slide: 0,
				ClearVirtualText: false,
			},
		},
		{
			name: "Repeats",
			keys: keyPresses("2j"),
			want: navigation.State{
				Buffer: "",
				Slide: 2,
				ClearVirtualText: true,
			},
		},
		{
			name: "Repeats with 0 (Vim ignores 0 and just does next command, so let's do the same)",
			keys: keyPresses("0j"),
			want: navigation.State{
				Buffer: "",
				Slide: 1,
				ClearVirtualText: true,
			},
		},
		{
			name: "Direct slide addressing sub min (Vim ignores minus, so let's do the same)",
			keys: keyPresses("-11G"),
			want: navigation.State{
				Buffer: "",
				Slide: 10,
				ClearVirtualText: false,
			},
		},
		{
			name: "Direct slide addressing min",
			keys: keyPresses("0G"),
			want: navigation.State{
				Buffer: "",
				Slide: 0,
				ClearVirtualText: false,
			},
		},
		{
			name: "Direct slide addressing in range",
			keys: keyPresses("3G"),
			want: navigation.State{
				Buffer: "",
				Slide: 2,
				ClearVirtualText: false,
			},
		},
		{
			name: "Direct slide addressing max",
			keys: keyPresses("11G"),
			want: navigation.State{
				Buffer: "",
				Slide: 10,
				ClearVirtualText: false,
			},
		},
		{
			name: "Direct slide addressing supra max",
			keys: keyPresses("101G"),
			want: navigation.State{
				Buffer: "",
				Slide: 10,
				ClearVirtualText: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentState := navigation.State{
				Buffer:           "",
				Slide:            0,
				ClearVirtualText: false,
			}
			numSlides := 11

			for _, key := range tt.keys {
				currentState = navigation.Navigate(currentState, key, numSlides)
			}

			assert.Equal(t, tt.want, currentState)
		})
	}
}
