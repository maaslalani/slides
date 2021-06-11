package styles_test

import (
	"testing"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/maaslalani/slides/styles"
	"github.com/stretchr/testify/assert"
)

func TestCustomTheme(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     []byte
		wantErr  bool
	}{
		{name: "Select custom theme json", filepath: "./theme.json", wantErr: false},
		{name: "Use an invalid filepath", filepath: "./someinvalidfile.toml", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := glamour.NewTermRenderer(styles.CustomTheme(tt.filepath))
			if err != nil {
				assert.True(t, tt.wantErr)
			}

			want, err := glamour.NewTermRenderer(glamour.WithStylesFromJSONFile(tt.filepath))
			if err != nil {
				assert.True(t, tt.wantErr)
			}

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, got, want)
		})
	}
}

func TestSelectTheme(t *testing.T) {
	tests := []struct {
		name  string
		theme string
		want  ansi.StyleConfig
	}{
		{name: "Select dark theme", theme: "dark", want: glamour.DarkStyleConfig},
		{name: "Select light theme", theme: "light", want: glamour.LightStyleConfig},
		{name: "Select ascii theme", theme: "ascii", want: glamour.ASCIIStyleConfig},
		{name: "Select notty theme", theme: "notty", want: glamour.NoTTYStyleConfig},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, styles.SelectTheme(tt.theme), tt.want)
		})
	}
}
