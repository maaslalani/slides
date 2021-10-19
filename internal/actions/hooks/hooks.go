package hooks

var Hooks = []HookFunc{
	searchHook,
	helpHook,
}

type Ctx struct {
	Prefix  string // ':', '/' or '?'
	Command string
	Args    []string
	Model   Model
}

type Model interface {
	SetPage(page int)
}

type HookFunc func(ctx *Ctx) (message string, accepted bool)
