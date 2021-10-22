package navigation

import (
	"regexp"
	"strings"
)

type searchFunc func(content string) bool

type searchType struct {
	Desc string
	fn   func(buf string) searchFunc
}

var SearchTypes = map[rune]*searchType{
	// header search
	// Search only in headings (case insensitive)
	0: {
		Desc: "header",
		fn: func(buf string) searchFunc {
			return func(content string) bool {
				return hasHeader(extractHeaders(content), buf) != ""
			}
		},
	},
	// full-text search
	// Search for any text on the presentation, including paragraphs, but also code blocks.
	// (case insensitive)
	'*': {
		Desc: "full-text-ci",
		fn: func(buf string) searchFunc {
			return func(content string) bool {
				return strings.Contains(strings.ToLower(content), strings.ToLower(buf))
			}
		},
	},
	// full-text search
	// Search for any text on the presentation, including paragraphs, but also code blocks.
	// (case sensitive)
	'^': {
		Desc: "full-text-cs",
		fn: func(buf string) searchFunc {
			return func(content string) bool {
				return strings.Contains(content, buf)
			}
		},
	},
	// regex search
	// Search for any text on the presentation, including paragraphs, but also code blocks.
	// (regex specific)
	'$': {
		Desc: "full-text-regex",
		fn: func(buf string) searchFunc {
			pattern, err := regexp.Compile(buf)
			if err != nil {
				return nil
			}
			return func(content string) bool {
				return pattern.MatchString(content)
			}
		},
	},
}

// Model is an interface for models.model, so that cycle imports are avoided
type Model interface {
	GetPage() int
	SetPage(page int)
	GetSlides() []string
}

// Search represents the current search
type Search struct {
	// Active - Show search bar instead of author and date?
	// Store keystrokes in Buffer?
	Active bool
	// Buffer stores the current "search term"
	Buffer string
}

// getSearchType returns the search type and the search-type identifier (*, $, ...).
// The search type should never be nil.
func (s *Search) getSearchType() (key rune, typ *searchType) {
	var ok bool
	if s.Buffer != "" {
		key = rune(s.Buffer[0])
		if typ, ok = SearchTypes[key]; !ok {
			key = 0
		} else {
			return
		}
	}
	typ = SearchTypes[key]
	return
}

// GetBufType returns the search type and the current search-buffer.
// If the search type is a non-standard type (such as full text search),
// the buffer is returned only after the second character, otherwise completely.
// The search type should never be nil.
func (s *Search) GetBufType() (buf string, typ *searchType) {
	var k rune
	if k, typ = s.getSearchType(); k == 0 {
		buf = s.Buffer
	} else {
		buf = s.Buffer[1:]
	}
	return
}

// Cancel the current search and delete the search buffer
func (s *Search) Cancel() {
	s.Buffer = ""
	s.Done()
}

// Mark Search as
// Done - Do not delete search buffer
// This is useful if, for example, you want to jump to the next result
// and you therefore still need the buffer
func (s *Search) Done() {
	s.Active = false
}

// Begin a new search (deletes old buffer)
func (s *Search) Begin() {
	s.Cancel() // clear buffer
	s.Active = true
}

// Write a keystroke to the buffer
func (s *Search) Write(key string) {
	s.Buffer += key
}

// Delete the last keystroke from the buffer
func (s *Search) Delete() {
	if len(s.Buffer) > 0 {
		s.Buffer = s.Buffer[0 : len(s.Buffer)-1]
	}
}

// Execute search
func (s *Search) Execute(m Model) {
	defer s.Done()
	buf, typ := s.GetBufType()
	// ignore empty buffers, also ignores '*<search term>', ...
	if buf == "" {
		return
	}
	check := func(i int) bool {
		content := m.GetSlides()[i]
		if typ.fn(buf)(content) {
			m.SetPage(i)
			return true
		}
		return false
	}
	// search from next slide to end
	for i := m.GetPage() + 1; i < len(m.GetSlides()); i++ {
		if check(i) {
			return
		}
	}
	// search from first slide to previous
	for i := 0; i < m.GetPage(); i++ {
		if check(i) {
			return
		}
	}
}

// extractHeaders is a hacky method to extract and return all headings of a slide
// Is probably much easier via RegEx: ^#+\s+([\w]+)$
func extractHeaders(slide string) (r []string) {
	code := false
	for _, l := range strings.Split(slide, "\n") {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "```") {
			code = !code
		}
		if code {
			continue
		}
		if !strings.HasPrefix(l, "#") || !strings.Contains(l, " ") {
			continue
		}
		header := strings.TrimSpace(l[strings.Index(l, " ")+1:])
		if header != "" {
			r = append(r, header)
		}
	}
	return
}

func hasHeader(headers []string, needle string) string {
	if len(headers) > 0 {
		for _, h := range headers {
			if strings.Contains(strings.ToLower(h), strings.ToLower(needle)) {
				return h
			}
		}
	}
	return ""
}
