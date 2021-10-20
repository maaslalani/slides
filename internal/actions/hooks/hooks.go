package hooks

var Hooks = []HookFunc{
	searchHook,
	helpHook,
	gotoHook,
}

type Ctx struct {
	Prefix  string // ':', '/' or '?'
	Command string
	Args    []string
	Model   Model
}

type Model interface {
	GetPage() int
	SetPage(page int)
	GetSlides() []string
}

type HookFunc func(ctx *Ctx) (message string, accepted bool)
