package navigation

import (
	"regexp"
	"strings"
)

// Model is an interface for models.model, so that cycle imports are avoided
type Model interface {
	CurrentPage() int
	SetPage(page int)
	Pages() []string
}

// Search represents the current search
type Search struct {
	// Active - Show search bar instead of author and date?
	// Store keystrokes in Query?
	Active bool
	// Query stores the current "search term"
	Query string
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
	s.Active = true
	s.Query = ""
}

// Write a keystroke to the buffer
func (s *Search) Write(key string) {
	s.Query += key
}

// Delete the last keystroke from the buffer
func (s *Search) Delete() {
	if len(s.Query) > 0 {
		s.Query = s.Query[0 : len(s.Query)-1]
	}
}

// Execute search
func (s *Search) Execute(m Model) {
	defer s.Done()
	if s.Query == "" {
		return
	}
	expr := s.Query
	if strings.HasSuffix(expr, "/i") {
		expr = "(?i)" + expr[:len(expr)-2]
	}
	pattern, err := regexp.Compile(expr)
	if err != nil {
		return
	}
	check := func(i int) bool {
		content := m.Pages()[i]
		if len(pattern.FindAllStringSubmatch(content, 1)) != 0 {
			m.SetPage(i)
			return true
		}
		return false
	}
	// search from next slide to end
	for i := m.CurrentPage() + 1; i < len(m.Pages()); i++ {
		if check(i) {
			return
		}
	}
	// search from first slide to previous
	for i := 0; i < m.CurrentPage(); i++ {
		if check(i) {
			return
		}
	}
}
