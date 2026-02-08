package navigation

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/maaslalani/slides/styles"
)

// Model is an interface for models.model, so that cycle imports are avoided
type Model interface {
	CurrentPage() int
	CurrentPart() int
	SetPage(page int)
	SetPageAndPart(page int, part int)
	Pages() []string
	SlideParts(page int) []string
}

// Search represents the current search
type Search struct {
	// Active - Show search bar instead of author and date?
	// Store keystrokes in Query?
	Active bool
	// Query stores the current "search term"
	SearchTextInput textinput.Model
}

// NewSearch creates and returns a new search model with the default settings.
func NewSearch() Search {
	ti := textinput.New()
	ti.Placeholder = "search"
	ti.Prompt = "/"
	ti.PromptStyle = styles.Search
	ti.TextStyle = styles.Search
	return Search{SearchTextInput: ti}
}

// Query returns the text input's value.
func (s *Search) Query() string {
	return s.SearchTextInput.Value()
}

// SetQuery sets the text input's value
func (s *Search) SetQuery(query string) {
	s.SearchTextInput.SetValue(query)
}

// Done marks the search as done, but does not delete the search buffer. This
// is useful if, for example, you want to jump to the next result and you
// therefore still need the buffer.
func (s *Search) Done() {
	s.Active = false
}

// Begin a new search (deletes old buffer)
func (s *Search) Begin() {
	s.Active = true
	s.SetQuery("")
}

// Execute search
func (s *Search) Execute(m Model) {
	defer s.Done()
	expr := s.Query()
	if expr == "" {
		return
	}
	if strings.HasSuffix(expr, "/i") {
		expr = "(?i)" + expr[:len(expr)-2]
	}
	pattern, err := regexp.Compile(expr)
	if err != nil {
		return
	}
	check := func(i int) bool {
		parts := m.SlideParts(i)
		for idx, part := range parts {
			if len(pattern.FindAllStringSubmatch(part, 1)) != 0 {
				m.SetPageAndPart(i, idx)
				return true
			}
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
