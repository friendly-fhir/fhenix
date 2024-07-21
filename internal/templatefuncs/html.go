package templatefuncs

import "html"

type HTMLModule struct{}

func (m *HTMLModule) Escape(s string) string {
	return html.EscapeString(s)
}

func (m *HTMLModule) Unescape(s string) string {
	return html.UnescapeString(s)
}
