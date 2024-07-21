package templatefuncs

import (
	"encoding/base32"
)

type Base32Module struct {
	Reporter Reporter
}

func (m *Base32Module) Encode(s string) string {
	return base32.StdEncoding.EncodeToString([]byte(s))
}

func (m *Base32Module) Decode(s string) string {
	b, err := base32.StdEncoding.DecodeString(s)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return string(b)
}

func (m *Base32Module) HexEncode(s string) string {
	return base32.HexEncoding.EncodeToString([]byte(s))
}

func (m *Base32Module) HexDecode(s string) string {
	b, err := base32.HexEncoding.DecodeString(s)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return string(b)
}
