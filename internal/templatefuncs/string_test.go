package templatefuncs_test

import (
	"testing"

	"github.com/friendly-fhir/fhenix/internal/templatefuncs"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestString_Upper(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "HELLO"},
		{"upper", "HELLO", "HELLO"},
		{"mixed", "HeLLo", "HELLO"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Upper(tc.s); got != tc.want {
				t.Errorf("StringModule.Upper() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Lower(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "hello"},
		{"upper", "HELLO", "hello"},
		{"mixed", "HeLLo", "hello"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Lower(tc.s); got != tc.want {
				t.Errorf("StringModule.Lower() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Title(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "Hello"},
		{"upper", "HELLO", "Hello"},
		{"mixed", "HeLLo", "Hello"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Title(tc.s); got != tc.want {
				t.Errorf("StringModule.Title() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Pascal(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "Hello"},
		{"upper", "HELLO", "Hello"},
		{"camel", "helloWorld", "HelloWorld"},
		{"kebab", "hello-world", "HelloWorld"},
		{"screaming", "HELLO_WORLD", "HelloWorld"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Pascal(tc.s); got != tc.want {
				t.Errorf("StringModule.Pascal() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Camel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "hello"},
		{"upper", "HELLO", "hello"},
		{"camel", "helloWorld", "helloWorld"},
		{"kebab", "hello-world", "helloWorld"},
		{"screaming", "HELLO_WORLD", "helloWorld"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Camel(tc.s); got != tc.want {
				t.Errorf("StringModule.Camel() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Snake(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "hello"},
		{"upper", "HELLO", "hello"},
		{"camel", "helloWorld", "hello_world"},
		{"kebab", "hello-world", "hello_world"},
		{"screaming", "HELLO_WORLD", "hello_world"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Snake(tc.s); got != tc.want {
				t.Errorf("StringModule.Snake() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Kebab(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "hello"},
		{"upper", "HELLO", "hello"},
		{"camel", "helloWorld", "hello-world"},
		{"kebab", "hello-world", "hello-world"},
		{"screaming", "HELLO_WORLD", "hello-world"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Kebab(tc.s); got != tc.want {
				t.Errorf("StringModule.Kebab() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Shout(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "HELLO"},
		{"upper", "HELLO", "HELLO"},
		{"camel", "helloWorld", "HELLO_WORLD"},
		{"kebab", "hello-world", "HELLO_WORLD"},
		{"screaming", "HELLO_WORLD", "HELLO_WORLD"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Shout(tc.s); got != tc.want {
				t.Errorf("StringModule.Shout() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_PascalInitialism(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "Hello"},
		{"upper", "HELLO", "Hello"},
		{"camel", "helloWorld", "HelloWorld"},
		{"kebab", "hello-world", "HelloWorld"},
		{"screaming", "HELLO_WORLD", "HelloWorld"},
		{"acronym", "http", "HTTP"},
		{"starts with acronym", "http-server", "HTTPServer"},
		{"ends with acronym", "http-server-2-0", "HTTPServer20"},
		{"acronym in middle", "serve-http-now", "ServeHTTPNow"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.PascalInitialism(tc.s); got != tc.want {
				t.Errorf("StringModule.PascalInitialism() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_CamelInitialism(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "hello"},
		{"upper", "HELLO", "hello"},
		{"camel", "helloWorld", "helloWorld"},
		{"kebab", "hello-world", "helloWorld"},
		{"screaming", "HELLO_WORLD", "helloWorld"},
		{"acronym", "http", "http"},
		{"starts with acronym", "http-server", "httpServer"},
		{"ends in acronym", "server-http", "serverHTTP"},
		{"acronym in middle", "serve-http-now", "serveHTTPNow"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.CamelInitialism(tc.s); got != tc.want {
				t.Errorf("StringModule.CamelInitialism() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Acronym(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"lower", "hello", "H"},
		{"upper", "HELLO", "H"},
		{"camel", "helloWorld", "HW"},
		{"kebab", "hello-world", "HW"},
		{"screaming", "HELLO_WORLD", "HW"},
		{"acronym", "http", "H"},
		{"starts with acronym", "http-server", "HS"},
		{"ends with acronym", "server-http", "SH"},
		{"acronym in middle", "serve-http-now", "SHN"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Acronym(tc.s); got != tc.want {
				t.Errorf("StringModule.Acronym() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_TrimSpace(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"space", " hello ", "hello"},
		{"tab", "\thello\t", "hello"},
		{"newline", "\nhello\n", "hello"},
		{"mixed", " \t\nhello \t\n", "hello"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.TrimSpace(tc.s); got != tc.want {
				t.Errorf("StringModule.Trim() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Trim(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		s      string
		cutset string
		want   string
	}{
		{"empty", "", "", ""},
		{"space", " hello ", " ", "hello"},
		{"tab", "\thello\t", "\t", "hello"},
		{"newline", "\nhello\n", "\n", "hello"},
		{"mixed", " \t\nhello \t\n", " \t\n", "hello"},
		{"cutset", "hello world", "hed ", "llo worl"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Trim(tc.cutset, tc.s); got != tc.want {
				t.Errorf("StringModule.Trim() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_TrimLeft(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		s      string
		cutset string
		want   string
	}{
		{"empty", "", "", ""},
		{"space", " hello ", " ", "hello "},
		{"tab", "\thello\t", "\t", "hello\t"},
		{"newline", "\nhello\n", "\n", "hello\n"},
		{"mixed", " \t\nhello \t\n", " \t\n", "hello \t\n"},
		{"cutset", "hello world", "hello ", "world"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.TrimLeft(tc.cutset, tc.s); got != tc.want {
				t.Errorf("StringModule.LTrim() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_TrimRight(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		s      string
		cutset string
		want   string
	}{
		{"empty", "", "", ""},
		{"space", " hello ", " ", " hello"},
		{"tab", "\thello\t", "\t", "\thello"},
		{"newline", "\nhello\n", "\n", "\nhello"},
		{"mixed", " \t\nhello \t\n", " \t\n", " \t\nhello"},
		{"cutset", "hello world", " world", "he"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.TrimRight(tc.cutset, tc.s); got != tc.want {
				t.Errorf("StringModule.TrimRight(%q, %q) = %q, want %q", tc.cutset, tc.s, got, tc.want)
			}
		})
	}
}

func TestString_CutPrefix(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		prefix string
		s      string
		want   string
	}{
		{"empty", "", "", ""},
		{"no prefix", "prefix", "text", "text"},
		{"single line", "prefix", "prefixtext", "text"},
		{"multi line", "prefix", "prefixline1\nprefixline2", "line1\nprefixline2"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.CutPrefix(tc.prefix, tc.s); got != tc.want {
				t.Errorf("StringModule.CutPrefix() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_CutSuffix(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		suffix string
		s      string
		want   string
	}{
		{"empty", "", "", ""},
		{"no suffix", "suffix", "text", "text"},
		{"single line", "suffix", "textsuffix", "text"},
		{"multi line", "suffix", "line1suffix\nline2suffix", "line1suffix\nline2"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.CutSuffix(tc.suffix, tc.s); got != tc.want {
				t.Errorf("StringModule.CutSuffix() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Fields(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want []string
	}{
		{"empty", "", []string{}},
		{"space", "hello world", []string{"hello", "world"}},
		{"tab", "hello\tworld", []string{"hello", "world"}},
		{"newline", "hello\nworld", []string{"hello", "world"}},
		{"mixed", "hello \t\nworld", []string{"hello", "world"}},
		{"multiple", "hello world  ", []string{"hello", "world"}},
		{"empty fields", "hello	world", []string{"hello", "world"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			got := m.Fields(tc.s)
			if !cmp.Equal(got, tc.want) {
				t.Errorf("StringModule.Fields() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Split(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		sep  string
		text string
		want []string
	}{
		{"empty", "", "", []string{}},
		{"space", " ", "hello world", []string{"hello", "world"}},
		{"tab", "\t", "hello\tworld", []string{"hello", "world"}},
		{"newline", "\n", "hello\nworld", []string{"hello", "world"}},
		{"mixed", " \t\n", "hello \t\nworld", []string{"hello", "world"}},
		{"empty fields", " ", "hello	world", []string{"hello	world"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			got := m.Split(tc.sep, tc.text)
			if !cmp.Equal(got, tc.want, cmpopts.EquateEmpty()) {
				t.Errorf("StringModule.Split() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Join(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		sep  string
		a    []string
		want string
	}{
		{"empty", "", []string{}, ""},
		{"space", " ", []string{"hello", "world"}, "hello world"},
		{"tab", "\t", []string{"hello", "world"}, "hello\tworld"},
		{"newline", "\n", []string{"hello", "world"}, "hello\nworld"},
		{"mixed", " \t\n", []string{"hello", "world"}, "hello \t\nworld"},
		{"empty fields", " ", []string{"hello	world"}, "hello	world"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.Join(tc.sep, tc.a); got != tc.want {
				t.Errorf("StringModule.Join() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Repeat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		n    int
		text string
		want string
	}{
		{"empty", 0, "", ""},
		{"zero", 0, "hello", ""},
		{"one", 1, "hello", "hello"},
		{"two", 2, "hello", "hellohello"},
		{"three", 3, "hello", "hellohellohello"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.Repeat(tc.n, tc.text); got != tc.want {
				t.Errorf("StringModule.Repeat() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Replace(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		old  string
		new  string
		text string
		want string
	}{
		{"empty", "", "", "", ""},
		{"no change", "hello", "world", "hello", "world"},
		{"single", "hello", "world", "hello world", "world world"},
		{"multiple", "hello", "world", "hello hello", "world world"},
		{"case sensitive", "Hello", "world", "hello Hello", "hello world"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.Replace(tc.old, tc.new, tc.text); got != tc.want {
				t.Errorf("StringModule.Replace() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Suffix(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		suffix  string
		content string
		want    string
	}{
		{"empty", "", "", ""},
		{"no suffix", "suffix", "text", "textsuffix"},
		{"single line", "suffix", "text", "textsuffix"},
		{"multi line", "suffix", "line1\nline2", "line1\nline2suffix"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.Suffix(tc.suffix, tc.content); got != tc.want {
				t.Errorf("StringModule.Suffix() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Prefix(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		prefix string
		text   string
		want   string
	}{
		{"empty", "", "", ""},
		{"no prefix", "prefix", "text", "prefixtext"},
		{"single line", "prefix", "text", "prefixtext"},
		{"multi line", "prefix", "line1\nline2", "prefixline1\nline2"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.Prefix(tc.prefix, tc.text); got != tc.want {
				t.Errorf("StringModule.Prefix() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Strip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"space", " hello ", "hello"},
		{"tab", "\thello\t", "hello"},
		{"newline", "\nhello\n", "hello"},
		{"mixed", " \t\nhello \t\n", "hello"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Strip(tc.s); got != tc.want {
				t.Errorf("StringModule.Strip() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Quote(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", `""`},
		{"simple", "hello", `"hello"`},
		{"escape", "hello\nworld", `"hello\nworld"`},
		{"quote", `"hello"`, `"\"hello\""`},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got := m.Quote(tc.s); got != tc.want {
				t.Errorf("StringModule.Quote() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Unquote(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", `""`, ""},
		{"simple", `"hello"`, "hello"},
		{"escape", `"hello\nworld"`, "hello\nworld"},
		{"quote", `"\"hello\""`, `"hello"`},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}
			if got, _ := m.Unquote(tc.s); got != tc.want {
				t.Errorf("StringModule.Unquote() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Char(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		n       int
		s       string
		want    string
		wantErr error
	}{
		{
			name:    "empty",
			n:       0,
			s:       "",
			want:    "",
			wantErr: templatefuncs.ErrIndexOutOfRange,
		}, {
			name: "zero",
			n:    0,
			s:    "hello",
			want: "h",
		}, {
			name: "one",
			n:    1,
			s:    "hello",
			want: "e",
		}, {
			name: "two",
			n:    2,
			s:    "hello",
			want: "l",
		}, {
			name: "three",
			n:    3,
			s:    "hello",
			want: "l",
		}, {
			name:    "out of bounds",
			n:       5,
			s:       "hello",
			want:    "",
			wantErr: templatefuncs.ErrIndexOutOfRange,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotErr error
			reporter := templatefuncs.ReporterFunc(func(err error) {
				gotErr = err
			})
			m := &templatefuncs.StringModule{
				Reporter: reporter,
			}

			got := m.Char(tc.n, tc.s)

			if want := tc.want; got != want {
				t.Errorf("StringModule.Char(%d, %q) = %q, want %q", tc.n, tc.s, got, want)
			}
			if got, want := gotErr, tc.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("StringModule.Char(%d, %q) error = %v, want %v", tc.n, tc.s, got, want)
			}
		})
	}
}

func TestString_First(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"simple", "hello", "h"},
		{"complex", "hello world", "h"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.First(tc.s); got != tc.want {
				t.Errorf("StringModule.First() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Last(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"simple", "hello", "o"},
		{"complex", "hello world", "d"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.Last(tc.s); got != tc.want {
				t.Errorf("StringModule.Last() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Reverse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"single", "a", "a"},
		{"double", "ab", "ba"},
		{"simple", "hello", "olleh"},
		{"complex", "hello world", "dlrow olleh"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.Reverse(tc.s); got != tc.want {
				t.Errorf("StringModule.Reverse() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestString_Substring(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		s      string
		start  int
		length int
		want   string
	}{
		{
			name:   "Empty string returns empty",
			s:      "",
			start:  0,
			length: 0,
			want:   "",
		}, {
			name:   "Substring at start",
			s:      "hello",
			start:  0,
			length: 2,
			want:   "he",
		}, {
			name:   "Substring in middle",
			s:      "hello",
			start:  1,
			length: 2,
			want:   "el",
		}, {
			name:   "Substring at end",
			s:      "hello",
			start:  3,
			length: 2,
			want:   "lo",
		}, {
			name:   "Substring longer than string",
			s:      "hello",
			start:  0,
			length: 10,
			want:   "hello",
		}, {
			name:   "Substring with negative start",
			s:      "hello",
			start:  -2,
			length: 2,
			want:   "",
		}, {
			name:   "Substring with negative length",
			s:      "hello",
			start:  2,
			length: -2,
			want:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &templatefuncs.StringModule{}

			if got := m.Substring(tc.start, tc.length, tc.s); got != tc.want {
				t.Errorf("StringModule.Substring(%d, %d, %q) = %q, want %q", tc.start, tc.length, tc.s, got, tc.want)
			}
		})
	}
}
