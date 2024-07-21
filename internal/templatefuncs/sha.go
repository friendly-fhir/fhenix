package templatefuncs

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

type SHAModule struct{}

func (m *SHAModule) Sum1(s string) string {
	result := sha1.Sum([]byte(s))
	return fmt.Sprintf("%x", result)
}

func (m *SHAModule) Sum224(s string) string {
	result := sha256.Sum224([]byte(s))
	return fmt.Sprintf("%x", result)
}

func (m *SHAModule) Sum256(s string) string {
	result := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", result)
}

func (m *SHAModule) Sum384(s string) string {
	result := sha512.Sum384([]byte(s))
	return fmt.Sprintf("%x", result)
}

func (m *SHAModule) Sum512(s string) string {
	result := sha512.Sum512([]byte(s))
	return fmt.Sprintf("%x", result)
}

func (m *SHAModule) Sum512_224(s string) string {
	result := sha512.Sum512_224([]byte(s))
	return fmt.Sprintf("%x", result)
}

func (m *SHAModule) Sum512_256(s string) string {
	result := sha512.Sum512_256([]byte(s))
	return fmt.Sprintf("%x", result)
}
