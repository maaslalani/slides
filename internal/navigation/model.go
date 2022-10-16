package navigation

// Model is an interface for models.model, so that cycle imports are avoided
type Model interface {
	CurrentPage() int
	SetPage(page int)
	Pages() []string
	SetVirtualText(string)
}

type mockModel struct {
	slides      []string
	page        int
	virtualText string
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

func (m *mockModel) SetVirtualText(text string) {
	m.virtualText = text
}
