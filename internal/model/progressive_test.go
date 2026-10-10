package model_test

import (
	"testing"

	"github.com/maaslalani/slides/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestProgressive_NoMarkers(t *testing.T) {
	slide := "# Title\nJust regular content"
	reveals := model.Progressive(slide)
	assert.Equal(t, []string{slide}, reveals)
}

func TestProgressive_SingleStop(t *testing.T) {
	slide := "# Title\nFirst point\n<!-- stop -->\nSecond point"
	expected := []string{
		"# Title\nFirst point",
		"# Title\nFirst point\nSecond point",
	}
	reveals := model.Progressive(slide)
	assert.Equal(t, expected, reveals)
}

func TestProgressive_SinglePause(t *testing.T) {
	slide := "# Title\nFirst point\n<!-- pause -->\nSecond point"
	expected := []string{
		"# Title\nFirst point",
		"# Title\nFirst point\nSecond point",
	}
	reveals := model.Progressive(slide)
	assert.Equal(t, expected, reveals)
}

func TestProgressive_MultipleStopsAndPauses(t *testing.T) {
	slide := "# Progressive\nStep 1\n<!-- stop -->\nStep 2\n<!-- pause -->\nStep 3\n<!-- stop -->\nStep 4"
	expected := []string{
		"# Progressive\nStep 1",
		"# Progressive\nStep 1\nStep 2",
		"# Progressive\nStep 1\nStep 2\nStep 3",
		"# Progressive\nStep 1\nStep 2\nStep 3\nStep 4",
	}
	reveals := model.Progressive(slide)
	assert.Equal(t, expected, reveals)
}

func TestProgressive_WhitespaceAndCasingVariants(t *testing.T) {
	tests := []struct {
		name     string
		slide    string
		expected []string
	}{
		{
			name:  "uppercase STOP",
			slide: "Intro\n<!-- STOP -->\nOutro",
			expected: []string{
				"Intro",
				"Intro\nOutro",
			},
		},
		{
			name:  "uppercase PAUSE",
			slide: "Intro\n<!-- PAUSE -->\nOutro",
			expected: []string{
				"Intro",
				"Intro\nOutro",
			},
		},
		{
			name:  "no spaces around keyword",
			slide: "Intro\n<!--stop-->\nOutro",
			expected: []string{
				"Intro",
				"Intro\nOutro",
			},
		},
		{
			name:  "extra spaces around keyword",
			slide: "Intro\n<!--   pause   -->\nOutro",
			expected: []string{
				"Intro",
				"Intro\nOutro",
			},
		},
		{
			name:  "mixed case",
			slide: "Intro\n<!-- Stop -->\nOutro",
			expected: []string{
				"Intro",
				"Intro\nOutro",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reveals := model.Progressive(tt.slide)
			assert.Equal(t, tt.expected, reveals)
		})
	}
}

func TestProgressive_ListItems(t *testing.T) {
	slide := "# Key Points\n* Point 1\n<!-- stop -->\n* Point 2\n<!-- stop -->\n* Point 3"
	expected := []string{
		"# Key Points\n* Point 1",
		"# Key Points\n* Point 1\n* Point 2",
		"# Key Points\n* Point 1\n* Point 2\n* Point 3",
	}
	reveals := model.Progressive(slide)
	assert.Equal(t, expected, reveals)
}

func TestProgressive_ParagraphSpacing(t *testing.T) {
	slide := "# Title\n\nParagraph 1\n\n<!-- stop -->\n\nParagraph 2"
	expected := []string{
		"# Title\n\nParagraph 1",
		"# Title\n\nParagraph 1\n\nParagraph 2",
	}
	reveals := model.Progressive(slide)
	assert.Equal(t, expected, reveals)
}

func TestProgressive_InlineMarker(t *testing.T) {
	slide := "Hello <!-- stop --> World"
	expected := []string{
		"Hello",
		"Hello  World",
	}
	reveals := model.Progressive(slide)
	assert.Equal(t, expected, reveals)
}

func TestProgressive_TrailingMarker(t *testing.T) {
	slide := "# Title\nPoint 1\n<!-- stop -->\n"
	expected := []string{
		"# Title\nPoint 1",
	}
	reveals := model.Progressive(slide)
	assert.Equal(t, expected, reveals)
}

func TestProgressive_LeadingMarker(t *testing.T) {
	slide := "<!-- stop -->\nPoint 1"
	expected := []string{
		"Point 1",
	}
	reveals := model.Progressive(slide)
	assert.Equal(t, expected, reveals)
}

func TestProgressive_ConsecutiveMarkers(t *testing.T) {
	slide := "Point 1\n<!-- stop -->\n<!-- pause -->\nPoint 2"
	expected := []string{
		"Point 1",
		"Point 1\nPoint 2",
	}
	reveals := model.Progressive(slide)
	assert.Equal(t, expected, reveals)
}

func TestParse_FullPresentation(t *testing.T) {
	content := `---
theme: dark
author: tester
---
# Slide 1
Intro

---
# Slide 2
First item
<!-- stop -->
Second item

---
# Slide 3
Final slide`

	slides, metaData := model.Parse(content)

	assert.Equal(t, "dark", metaData.Theme)
	assert.Equal(t, "tester", metaData.Author)
	assert.Len(t, slides, 4)
	assert.Equal(t, "# Slide 1\nIntro", slides[0])
	assert.Equal(t, "# Slide 2\nFirst item", slides[1])
	assert.Equal(t, "# Slide 2\nFirst item\nSecond item", slides[2])
	assert.Equal(t, "# Slide 3\nFinal slide", slides[3])
}

func TestParse_NoFrontmatter(t *testing.T) {
	content := `# Slide 1
Only slide
<!-- stop -->
Revealed part`

	slides, _ := model.Parse(content)

	assert.Len(t, slides, 2)
	assert.Equal(t, "# Slide 1\nOnly slide", slides[0])
	assert.Equal(t, "# Slide 1\nOnly slide\nRevealed part", slides[1])
}
