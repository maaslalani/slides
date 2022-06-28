package termenv

import (
	"text/template"
)

// TemplateFuncs contains a few useful template helpers.
func TemplateFuncs(p Profile) template.FuncMap {
	return template.FuncMap{
		"Color": func(values ...interface{}) string {
			s := p.String(values[len(values)-1].(string))
			switch len(values) {
			case 2:
				s = s.Foreground(p.Color(values[0].(string)))
			case 3:
				s = s.
					Foreground(p.Color(values[0].(string))).
					Background(p.Color(values[1].(string)))
			}

			return s.String()
		},
		"Foreground": func(values ...interface{}) string {
			s := p.String(values[len(values)-1].(string))
			if len(values) == 2 {
				s = s.Foreground(p.Color(values[0].(string)))
			}

			return s.String()
		},
		"Background": func(values ...interface{}) string {
			s := p.String(values[len(values)-1].(string))
			if len(values) == 2 {
				s = s.Background(p.Color(values[0].(string)))
			}

			return s.String()
		},
		"Bold":      styleFunc(p, Style.Bold),
		"Faint":     styleFunc(p, Style.Faint),
		"Italic":    styleFunc(p, Style.Italic),
		"Underline": styleFunc(p, Style.Underline),
		"Overline":  styleFunc(p, Style.Overline),
		"Blink":     styleFunc(p, Style.Blink),
		"Reverse":   styleFunc(p, Style.Reverse),
		"CrossOut":  styleFunc(p, Style.CrossOut),
	}
}

func styleFunc(p Profile, f func(Style) Style) func(...interface{}) string {
	return func(values ...interface{}) string {
		s := p.String(values[0].(string))
		return f(s).String()
	}
}
