package model

import "testing"

func TestSplitSlidesWithPausedDelimiter(t *testing.T) {
	slides := []string{
		"# Title\nfirst\n+++\nsecond\n+++\nthird",
		"# Next\nonly",
	}

	parsed := splitSlides(slides)
	if len(parsed) != 2 {
		t.Fatalf("expected 2 slides, got %d", len(parsed))
	}
	if len(parsed[0].Parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parsed[0].Parts))
	}
	if parsed[0].Parts[1] != "second" {
		t.Fatalf("expected second part to be %q, got %q", "second", parsed[0].Parts[1])
	}
	if len(parsed[1].Parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(parsed[1].Parts))
	}
}

func TestAdvanceRetreatSteps(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a", "b"}},
			{Parts: []string{"c", "d", "e"}},
		},
	}

	m.SetPageAndPart(0, 0)
	m.advanceSteps(1)
	if m.Page != 0 || m.Part != 1 {
		t.Fatalf("expected page 0 part 1, got page %d part %d", m.Page, m.Part)
	}

	m.advanceSteps(1)
	if m.Page != 1 || m.Part != 0 {
		t.Fatalf("expected page 1 part 0, got page %d part %d", m.Page, m.Part)
	}

	m.retreatSteps(1)
	if m.Page != 0 || m.Part != 1 {
		t.Fatalf("expected page 0 part 1 after retreat, got page %d part %d", m.Page, m.Part)
	}
}

func TestCurrentSlideContentCumulative(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"first", "second", "third"}},
		},
	}
	m.SetPageAndPart(0, 1)

	content := m.currentSlideContent()
	if content != "first\nsecond" {
		t.Fatalf("expected cumulative content, got %q", content)
	}
}

func TestPagingFormats(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a", "b", "c"}},
			{Parts: []string{"d"}},
		},
	}
	m.SetPageAndPart(0, 1)

	tests := []struct {
		format   string
		expected string
	}{
		// No placeholder
		{"Static", "Static"},
		// 1 placeholder: slide
		{"Slide %d", "Slide 1"},
		// 2 placeholders: slide, totalSlides
		{"Slide %d / %d", "Slide 1 / 2"},
		// 3 placeholders: slide, totalSlides, part
		{"Slide %d / %d | Part %d", "Slide 1 / 2 | Part 2"},
		// 4+ placeholders: slide, totalSlides, part, totalParts
		{"Slide %d / %d (Part %d / %d)", "Slide 1 / 2 (Part 2 / 3)"},
	}

	for _, tt := range tests {
		m.Paging = tt.format
		result := m.paging()
		if result != tt.expected {
			t.Errorf("format %q: expected %q, got %q", tt.format, tt.expected, result)
		}
	}
}

func TestPagingWithLastPart(t *testing.T) {
	m := Model{
		Slides: []Slide{
			{Parts: []string{"a", "b"}},
			{Parts: []string{"c", "d", "e"}},
		},
	}
	m.SetPageAndPart(1, 2) // Last slide, last part
	m.Paging = "Slide %d / %d (Part %d / %d)"

	result := m.paging()
	expected := "Slide 2 / 2 (Part 3 / 3)"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
