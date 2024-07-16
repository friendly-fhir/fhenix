package ansi

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	// enabled is used to detect if
	enabled = true

	createFormatFunc func(...byte) string
	formatFunc       func(string, ...any) string

	stripCodes = regexp.MustCompile("\033\\[([0-9;]+)m")
)

// DisableColors will disable all color formatting globally
func DisableColors() {
	createFormatFunc = noFormat
	formatFunc = stripSprintf
	enabled = false
}

func noFormat(_ ...byte) string {
	return ""
}

func ansiFormat(codes ...byte) string {
	sb := strings.Builder{}
	sb.WriteString(ControlSequenceIntroducer)
	for _, f := range codes[:len(codes)-1] {
		sb.WriteString(strconv.Itoa(int(f)))
		sb.WriteRune(Separator)
	}
	sb.WriteString(strconv.Itoa(int(codes[len(codes)-1])))
	sb.WriteRune(SGRSuffix)
	return sb.String()
}

func stripSprintf(format string, args ...any) string {
	content := fmt.Sprintf(format, args...)
	if !strings.ContainsRune(format, ControlCode) {
		return content
	}
	return string(stripCodes.ReplaceAll([]byte(format), nil))
}

func init() {
	enabled = true
	createFormatFunc = ansiFormat
	formatFunc = fmt.Sprintf
	unsetIf("NOCOLOR")
	unsetIf("NO_COLOR")
}

func unsetIf(key string) {
	if got, ok := os.LookupEnv(key); ok {
		if got, err := strconv.ParseBool(got); err == nil && got {
			createFormatFunc = noFormat
			formatFunc = stripSprintf
			enabled = false
		}
	}
}
