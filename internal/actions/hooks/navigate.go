package hooks

import (
	"strconv"
	"strings"
)

// searchHook - jump to a slide containing a that contains the <search term>
// /<search term> - forward search
// ?<search term> - backward search
var searchHook HookFunc = func(ctx *Ctx) (msg string, accept bool) {
	if ctx.Prefix != "/" && ctx.Prefix != "?" {
		return "", false
	}
	check := func(ctx *Ctx, i int) bool {
		content := ctx.Model.GetSlides()[i]
		headers := extractHeaders(content)
		if h := hasHeader(headers, ctx.Command); h != "" {
			ctx.Model.SetPage(i)
			accept = true
			return true
		}
		return false
	}
	// forward search
	if ctx.Prefix == "/" {
		// search from next slide to end
		for i := ctx.Model.GetPage() + 1; i < len(ctx.Model.GetSlides()); i++ {
			if check(ctx, i) {
				return
			}
		}
		// search from first slide to previous
		for i := 0; i < ctx.Model.GetPage(); i++ {
			if check(ctx, i) {
				return
			}
		}
	} else if ctx.Prefix == "?" {
		// search from previous slide to start
		for i := ctx.Model.GetPage() - 1; i >= 0; i-- {
			if check(ctx, i) {
				return
			}
		}
		// search from end to next slide
		for i := len(ctx.Model.GetSlides()) - 1; i > ctx.Model.GetPage(); i-- {
			if check(ctx, i) {
				return
			}
		}
	}
	return "", true
}

// gotoHook - :<slide>
// jump to a slide
var gotoHook HookFunc = func(ctx *Ctx) (msg string, accept bool) {
	// check if command is a number
	if n, err := strconv.Atoi(ctx.Command); err == nil {
		accept = true
		if n > 0 && len(ctx.Model.GetSlides()) >= n {
			ctx.Model.SetPage(n - 1)
		}
		return
	}
	return "", false
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
