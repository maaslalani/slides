package navigation

import (
	"strings"
)

type Model interface {
	GetPage() int
	SetPage(page int)
	GetSlides() []string
}

type Search struct {
	Active bool
	Buffer string
}

func (s *Search) Begin() {
	s.Cancel() // clear buffer
	s.Active = true
}

func (s *Search) Write(key string) {
	s.Buffer += key
}

func (s *Search) Delete() {
	if len(s.Buffer) > 0 {
		s.Buffer = s.Buffer[0 : len(s.Buffer)-1]
	}
}

func (s *Search) Execute(m Model) {
	defer s.Cancel()

	check := func(i int) bool {
		content := m.GetSlides()[i]
		headers := extractHeaders(content)
		if h := hasHeader(headers, s.Buffer); h != "" {
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

func (s *Search) Cancel() {
	s.Active = false
	s.Buffer = ""
}

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
