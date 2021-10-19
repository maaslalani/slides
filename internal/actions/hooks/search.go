package hooks

var searchHook HookFunc = func(ctx *Ctx) (string, bool) {
	// forward search
	if ctx.Prefix == "/" {
		ctx.Model.SetPage(1)
		return "searching ...", true
	}
	return "", false
}
