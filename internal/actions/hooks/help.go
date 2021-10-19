package hooks

var helpHook HookFunc = func(ctx *Ctx) (string, bool) {
	if ctx.Command == "h" || ctx.Command == "help" {
		return "help: use /<term> to search.", true
	}
	return "", false
}
