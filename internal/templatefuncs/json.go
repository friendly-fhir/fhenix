package templatefuncs

import "encoding/json"

type JSONModule struct {
	Reporter Reporter
}

func (m *JSONModule) Encode(data any) string {
	b, err := json.Marshal(data)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
		return StringOnError
	}
	return string(b)
}

func (m *JSONModule) EncodeIndent(data any) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
		return StringOnError
	}
	return string(b)
}

func (m *JSONModule) Decode(s string) map[string]any {
	data := map[string]any{}
	err := json.Unmarshal([]byte(s), &data)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return data
}
