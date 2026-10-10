package model

import (
	"regexp"
	"strings"

	"github.com/maaslalani/slides/internal/meta"
)

var (
	stopMarkerRegex       = regexp.MustCompile(`(?i)<!--\s*(?:stop|pause)\s*-->`)
	stopLineRegex         = regexp.MustCompile(`(?mi)^[ \t]*<!--\s*(?:stop|pause)\s*-->[ \t]*\r?\n?`)
	multipleNewlinesRegex = regexp.MustCompile(`\n{3,}`)
)

// cleanStopMarkers removes any stop or pause comment markers from the content.
func cleanStopMarkers(s string) string {
	s = stopLineRegex.ReplaceAllString(s, "")
	s = stopMarkerRegex.ReplaceAllString(s, "")
	s = multipleNewlinesRegex.ReplaceAllString(s, "\n\n")
	return s
}

// Progressive splits a single slide containing progressive reveal comment markers
// (<!-- stop --> or <!-- pause -->) into progressive reveal sub-slides.
// If the slide does not contain any markers, a slice containing the slide is returned.
func Progressive(slide string) []string {
	indices := stopMarkerRegex.FindAllStringIndex(slide, -1)
	if len(indices) == 0 {
		return []string{strings.TrimRight(slide, "\r\n")}
	}

	var reveals []string
	for _, loc := range indices {
		step := cleanStopMarkers(slide[:loc[0]])
		step = strings.TrimRight(step, " \t\r\n")
		if strings.TrimSpace(step) == "" {
			continue
		}
		if len(reveals) > 0 && strings.TrimSpace(reveals[len(reveals)-1]) == strings.TrimSpace(step) {
			continue
		}
		reveals = append(reveals, step)
	}

	finalStep := cleanStopMarkers(slide)
	finalStep = strings.TrimRight(finalStep, " \t\r\n")
	if strings.TrimSpace(finalStep) != "" {
		if len(reveals) == 0 || strings.TrimSpace(reveals[len(reveals)-1]) != strings.TrimSpace(finalStep) {
			reveals = append(reveals, finalStep)
		}
	}

	if len(reveals) == 0 {
		return []string{strings.TrimRight(cleanStopMarkers(slide), " \t\r\n")}
	}

	return reveals
}

// Parse parses presentation markdown content into individual slides,
// expanding any progressive reveal comments (<!-- stop --> or <!-- pause -->)
// and extracting presentation metadata frontmatter.
func Parse(content string) ([]string, *meta.Meta) {
	content = strings.ReplaceAll(content, "\r", "")

	content = strings.TrimPrefix(content, strings.TrimPrefix(delimiter, "\n"))
	rawSlides := strings.Split(content, delimiter)

	metaData, exists := meta.New().Parse(rawSlides[0])
	if exists && len(rawSlides) > 1 {
		rawSlides = rawSlides[1:]
	}

	var slides []string
	for _, slide := range rawSlides {
		slides = append(slides, Progressive(slide)...)
	}

	return slides, metaData
}
