package templatefuncs

import (
	"go/build"
	"runtime/debug"
)

type BuildModule struct{}

func (m *BuildModule) Context() *build.Context {
	return &build.Default
}

func (m *BuildModule) VCS() string {
	for _, setting := range m.settings() {
		if setting.Key == "vcs" {
			return setting.Value
		}
	}
	return "unknown"
}

func (m *BuildModule) VCSRevision() string {
	for _, setting := range m.settings() {
		if setting.Key == "vcs.revision" {
			return setting.Value
		}
	}
	return "unknown"
}

func (m *BuildModule) VCSTime() string {
	for _, setting := range m.settings() {
		if setting.Key == "vcs.time" {
			return setting.Value
		}
	}
	return "unknown"
}

func (m *BuildModule) settings() []debug.BuildSetting {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Settings
	}
	return nil
}
