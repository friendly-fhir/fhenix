package templatefuncs

import (
	"compress/gzip"
	"io"
	"strings"
)

type GZipModule struct {
	Reporter Reporter
}

func (m *GZipModule) Encode(s string) string {
	var sb strings.Builder
	gz := gzip.NewWriter(&sb)

	if _, err := gz.Write([]byte(s)); err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	if err := gz.Close(); err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return sb.String()
}

func (m *GZipModule) Decode(s string) string {
	r, err := gzip.NewReader(strings.NewReader(s))
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
		return ""
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return string(b)
}
