package templatefuncs

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type StringModule struct {
	Reporter Reporter
}

func (m *StringModule) Upper(s string) string {
	return cases.Upper(language.English).String(s)
}

func (m *StringModule) Lower(s string) string {
	return cases.Lower(language.English).String(s)
}

func (m *StringModule) Title(s string) string {
	return cases.Title(language.English).String(s)
}

func (m *StringModule) Pascal(s string) string {
	return strcase.ToCamel(s)
}

func (m *StringModule) Camel(s string) string {
	return strcase.ToLowerCamel(s)
}

func (m *StringModule) Snake(s string) string {
	return strcase.ToSnake(s)
}

func (m *StringModule) Kebab(s string) string {
	return strcase.ToKebab(s)
}

func (m *StringModule) Shout(s string) string {
	return strcase.ToScreamingSnake(s)
}

var acronyms = map[string]struct{}{
	"id":    {},
	"url":   {},
	"uri":   {},
	"uuid":  {},
	"oid":   {},
	"json":  {},
	"xml":   {},
	"html":  {},
	"http":  {},
	"https": {},
	"xhtml": {},
}

func (m *StringModule) PascalInitialism(s string) string {
	parts := strings.Split(strings.ToLower(strcase.ToKebab(s)), "-")
	for i, part := range parts {
		if _, ok := acronyms[part]; !ok {
			parts[i] = strcase.ToCamel(part)
		} else {
			parts[i] = strings.ToUpper(part)
		}
	}
	return strings.Join(parts, "")
}

func (m *StringModule) CamelInitialism(s string) string {
	parts := strings.Split(strings.ToLower(strcase.ToKebab(s)), "-")
	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
			continue
		}
		if _, ok := acronyms[part]; !ok {
			parts[i] = strcase.ToCamel(part)
		} else {
			parts[i] = strings.ToUpper(part)
		}
	}
	return strings.Join(parts, "")
}

func (m *StringModule) Acronym(s string) string {
	parts := strings.Split(strings.ToLower(strcase.ToKebab(s)), "-")
	for i, part := range parts {
		if len(part) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(string(part[0]))
	}
	return strings.Join(parts, "")
}

func (m *StringModule) TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func (m *StringModule) Trim(cutset, s string) string {
	return strings.Trim(s, cutset)
}

func (m *StringModule) TrimLeft(cutset, s string) string {
	return strings.TrimLeft(s, cutset)
}

func (m *StringModule) TrimRight(cutset, s string) string {
	return strings.TrimRight(s, cutset)
}

func (m *StringModule) CutPrefix(prefix, s string) string {
	result, _ := strings.CutPrefix(s, prefix)
	return result
}

func (m *StringModule) CutSuffix(suffix, s string) string {
	result, _ := strings.CutSuffix(s, suffix)
	return result
}

func (m *StringModule) Fields(s string) []string {
	return strings.Fields(s)
}

func (m *StringModule) Split(sep, text string) []string {
	return strings.Split(text, sep)
}

func (m *StringModule) Join(sep string, a []string) string {
	return strings.Join(a, sep)
}

func (m *StringModule) Repeat(n int, text string) string {
	return strings.Repeat(text, n)
}

func (m *StringModule) Replace(old, new, text string) string {
	return strings.ReplaceAll(text, old, new)
}

func (m *StringModule) Suffix(suffix, content string) string {
	return content + suffix
}

func (m *StringModule) Prefix(prefix, content string) string {
	return prefix + content
}

func (m *StringModule) Strip(content string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, content)
}

func (m *StringModule) Quote(content string) string {
	return strconv.Quote(content)
}

func (m *StringModule) Unquote(content string) (string, error) {
	return strconv.Unquote(content)
}

func (m *StringModule) Char(n int, s string) string {
	if n < 0 || n >= len(s) {
		if m.Reporter != nil {
			m.Reporter.Report(fmt.Errorf("%w: %d", ErrIndexOutOfRange, n))
		}
		return ""
	}
	return string(s[n])
}

func (m *StringModule) First(s string) string {
	return m.Char(0, s)
}

func (m *StringModule) Last(s string) string {
	return m.Char(len(s)-1, s)
}

func (m *StringModule) Reverse(s string) string {
	runes := []rune(s)
	for i, ch := range runes[:len(runes)/2] {
		idx := len(runes) - 1 - i
		runes[i] = runes[idx]
		runes[idx] = ch
	}
	return string(runes)
}

func (m *StringModule) Substring(start, length int, s string) string {
	runes := []rune(s)
	if start < 0 || length <= 0 {
		return ""
	}
	end := min(len(s), start+length)
	return string(runes[start:end])
}
