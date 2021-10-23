package navigation

import (
	"regexp"
	"strings"
)

// Model is an interface for models.model, so that cycle imports are avoided
type Model interface {
	CurrentPage() int
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
	// ignore empty buffers, also ignores '*<search term>', ...
	if s.Buffer == "" {
		return
	}
	// compile pattern
	expr := s.Buffer
	if strings.HasSuffix(expr, "/i") {
		expr = "(?i)" + expr[:len(expr)-2]
	}
	pattern, err := regexp.Compile(expr)
	if err != nil {
		return
	}
	check := func(i int) bool {
		content := m.GetSlides()[i]
		if len(pattern.FindAllStringSubmatch(content, 1)) != 0 {
			m.SetPage(i)
			return true
		}
		return false
	}
	// search from next slide to end
	for i := m.CurrentPage() + 1; i < len(m.GetSlides()); i++ {
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
