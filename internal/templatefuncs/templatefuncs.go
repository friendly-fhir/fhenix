/*
Package templatefuncs provides a set of template functions that are used
in the internal template package.
*/
package templatefuncs

import (
	"compress/gzip"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var DefaultFuncs = map[string]any{
	"uppercase":  cases.Upper(language.English).String,
	"lowercase":  cases.Lower(language.English).String,
	"titlecase":  cases.Title(language.English).String,
	"pascalcase": strcase.ToCamel,
	"camelcase":  strcase.ToLowerCamel,
	"snakecase":  strcase.ToSnake,
	"kebabcase":  strcase.ToKebab,
	"shoutcase":  strcase.ToScreamingSnake,
	"pascalinitialcase": func(s string) string {
		parts := strings.Split(strings.ToLower(strcase.ToKebab(s)), "-")
		for i, part := range parts {
			if _, ok := acronyms[part]; !ok {
				parts[i] = strcase.ToCamel(part)
			} else {
				parts[i] = strings.ToUpper(part)
			}
		}
		return strings.Join(parts, "")
	},
	"camelinitialcase": func(s string) string {
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
	},
	"acronym": func(s string) string {
		parts := strings.Split(strings.ToLower(strcase.ToKebab(s)), "-")
		for i, part := range parts {
			parts[i] = strings.ToUpper(string(part[0]))
		}
		return strings.Join(parts, "")
	},
	"fold": cases.Fold().String,

	"trim":  strings.TrimSpace,
	"ltrim": strings.TrimLeft,
	"rtrim": strings.TrimRight,

	"fields":  strings.Fields,
	"split":   func(sep, text string) []string { return strings.Split(text, sep) },
	"join":    func(sep string, a []string) string { return strings.Join(a, sep) },
	"repeat":  func(n int, text string) string { return strings.Repeat(text, n) },
	"replace": func(old, new, text string) string { return strings.ReplaceAll(text, old, new) },
	"prefix": func(prefix, text string) string {
		return prefix + strings.ReplaceAll(text, "\n", "\n"+prefix)
	},
	"suffix": func(suffix, text string) string {
		return strings.ReplaceAll(text, "\n", suffix+"\n") + suffix
	},
	"append": func(suffix, content string) any {
		return content + suffix
	},
	"prepend": func(prefix, content string) any {
		return prefix + content
	},
	"indent": func(indent int, text string) string {
		space := strings.Repeat(" ", indent)
		return space + strings.ReplaceAll(text, "\n", "\n"+space)
	},
	"tabindent": func(indent int, text string) string {
		tab := strings.Repeat("\t", indent)
		return tab + strings.ReplaceAll(text, "\n", "\n"+tab)
	},
	"resize": func(columns int, text string) string {
		var sb strings.Builder
		lines := strings.Split(text, "\n")
		length := 0
		for _, line := range lines {
			// Preserve empty line breaks
			if strings.TrimSpace(line) == "" {
				sb.WriteString("\n")
				length = 0
				continue
			}
			tokens := strings.Fields(line)
			for i, token := range tokens {
				if i > 0 && length+len(token) > columns {
					sb.WriteString("\n")
					length = 0
				}
				sb.WriteString(token)
				sb.WriteString(" ")
				length += len(token) + 1
			}
		}
		return strings.TrimSpace(sb.String())
	},
	"strip": func(content string) string {
		return strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsNumber(r) {
				return r
			}
			return -1
		}, content)
	},
	"quote":    strconv.Quote,
	"unquote":  strconv.Unquote,
	"escape":   html.EscapeString,
	"unescape": html.UnescapeString,

	"cutset":    func(set, text string) string { return strings.Trim(text, set) },
	"cutprefix": func(prefix, text string) string { return strings.TrimPrefix(text, prefix) },
	"cutsuffix": func(suffix, text string) string { return strings.TrimSuffix(text, suffix) },

	"gzip": func(content string) string {
		var sb strings.Builder
		fmt.Fprintf(gzip.NewWriter(&sb), "%s", content)
		return strconv.Quote(sb.String())
	},

	"base32": func(content string) string {
		return base32.StdEncoding.EncodeToString([]byte(content))
	},
	"base64": func(content string) string {
		return base64.StdEncoding.EncodeToString([]byte(content))
	},

	"json": func(data any) string {
		b, _ := json.Marshal(data)
		return string(b)
	},

	"first": func(v any) any {
		return reflect.ValueOf(v).Index(0).Interface()
	},
	"last": func(v any) any {
		rv := reflect.ValueOf(v)
		return rv.Index(rv.Len() - 1).Interface()
	},

	"char": func(n int, text string) string { return string(text[n]) },
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
