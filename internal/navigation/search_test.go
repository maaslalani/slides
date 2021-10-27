package navigation

import (
	"testing"
)

type mockModel struct {
	slides []string
	page   int
}

func (m *mockModel) CurrentPage() int {
	return m.page
}

func (m *mockModel) SetPage(page int) {
	m.page = page
}

func (m *mockModel) Pages() []string {
	return m.slides
}

func TestSearch(t *testing.T) {
	data := []string{
		"hi",
		"first",
		"second",
		"third",
		"AbCdEfG",
		"abcdefg",
		"seconds",
	}

	type query struct {
		desc     string
		query    string
		expected int
	}

	// query -> expected page
	queries := []query{
		{"basic 'first'", "first", 1},
		{"basic 'abc'", "abc", 5},
		{"basic 'abc' next occurrence", "abc", 5},
		{"'abc' ignore case", "abc/i", 4},
		{"'abc' ignore case", "abc/i", 5},
		{"'abc' ignore case", "abc/i", 4},
		{"next occurrence 1/2", "sec", 6},
		{"next occurrence 2/2", "sec", 2},
		{"regex", "a.c", 5},
		{"regex next occurrence", "a.c", 5},
		{"regex ignore case", "a.c/i", 4},
		{"regex ignore case next occurrence", "a.c/i", 5},
	}

	m := &mockModel{
		slides: data,
		page:   0,
	}

	s := &Search{}
	for _, query := range queries {
		s.Query = query.query
		s.Execute(m)
		if m.CurrentPage() != query.expected {
			t.Errorf("[%s] expected page %d, got %d", query.desc, query.expected, m.CurrentPage())
		}
	}

}
