package templatefuncs

import "gopkg.in/yaml.v3"

type YAMLModule struct {
	Reporter Reporter
}

func (m *YAMLModule) Encode(data any) string {
	b, err := yaml.Marshal(data)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return string(b)
}

func (m *YAMLModule) Decode(s string) map[string]any {
	data := map[string]any{}
	err := yaml.Unmarshal([]byte(s), &data)
	if err != nil && m.Reporter != nil {
		m.Reporter.Report(err)
	}
	return data
}
