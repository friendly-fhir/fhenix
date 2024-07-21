package templatefuncs

import "encoding/base64"

type Base64Module struct {
	Reporter Reporter
}

func (m *Base64Module) Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func (m *Base64Module) Decode(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return string(b)
}

func (m *Base64Module) RawEncode(s string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(s))
}

func (m *Base64Module) RawDecode(s string) string {
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return string(b)
}

func (m *Base64Module) URLEncode(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}

func (m *Base64Module) URLDecode(s string) string {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return string(b)
}

func (m *Base64Module) RawURLEncode(s string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

func (m *Base64Module) RawURLDecode(s string) string {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return string(b)
}
