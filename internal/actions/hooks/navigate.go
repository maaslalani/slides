package hooks

import (
	"fmt"
	"strconv"
	"strings"
)

var searchHook HookFunc = func(ctx *Ctx) (msg string, accept bool) {
	check := func(ctx *Ctx, i int) bool {
		content := ctx.Model.GetSlides()[i]
		headers := extractHeaders(content)
		if h := hasHeader(headers, ctx.Command); h != "" {
			ctx.Model.SetPage(i)
			msg = "found: " + h
			accept = true
			return true
		}
		return false
	}

	// forward search
	if ctx.Prefix == "/" {
		for i := ctx.Model.GetPage() + 1; i < len(ctx.Model.GetSlides()); i++ {
			if check(ctx, i) {
				return
			}
		}
		for i := 0; i < ctx.Model.GetPage(); i++ {
			if check(ctx, i) {
				return
			}
		}
		return "cannot forward-find: " + ctx.Command, true
	} else if ctx.Prefix == "?" {
		for i := ctx.Model.GetPage() - 1; i >= 0; i-- {
			if check(ctx, i) {
				return
			}
		}
		for i := len(ctx.Model.GetSlides()) - 1; i > ctx.Model.GetPage(); i-- {
			if check(ctx, i) {
				return
			}
		}
		return "cannot backward-find: " + ctx.Command, true
	}
	return "", false
}

var gotoHook HookFunc = func(ctx *Ctx) (msg string, accept bool) {
	// check if command is a number
	if n, err := strconv.Atoi(ctx.Command); err == nil {
		accept = true
		if n <= 0 {
			msg = fmt.Sprintf("error: cannot go to slide 0")
		} else if l := len(ctx.Model.GetSlides()); l < n {
			msg = fmt.Sprintf("error: input %d > %d", n, l)
		} else {
			msg = fmt.Sprintf("ok: switched to slide %d", n)
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
