package navigation

// Model is an interface for models.model, so that cycle imports are avoided
type Model interface {
	CurrentPage() int
	SetPage(page int)
	Pages() []string
	SetVirtualText(string)
}
