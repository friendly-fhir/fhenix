package templatefuncs

import (
	"fmt"
	"reflect"
)

type ListModule struct {
	Reporter Reporter
}

func (m *ListModule) First(v any) any {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		if rv.Len() > 0 {
			return rv.Index(0).Interface()
		}
	case reflect.String:
		s := rv.String()
		if len(s) > 0 {
			return s[:1]
		}
		return ""
	}
	if m.Reporter != nil {
		m.Reporter.Report(fmt.Errorf("%w %T for list.First", ErrInvalidType, v))
	}
	return nil
}

func (m *ListModule) Last(v any) any {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		if rv.Len() > 0 {
			return rv.Index(rv.Len() - 1).Interface()
		}
	case reflect.String:
		s := rv.String()
		if len(s) > 0 {
			return s[len(s)-1:]
		}
		return ""
	}
	if m.Reporter != nil {
		m.Reporter.Report(fmt.Errorf("%w %T for list.Last", ErrInvalidType, v))
	}
	return nil
}

func (m *ListModule) Length(v any) int {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		return rv.Len()
	}
	if m.Reporter != nil {
		m.Reporter.Report(fmt.Errorf("%w %T for list.Length", ErrInvalidType, v))
	}
	return 0
}

func (m *ListModule) Index(i int, v any) any {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		if i >= 0 && i < rv.Len() {
			return rv.Index(i).Interface()
		}
		if m.Reporter != nil {
			m.Reporter.Report(fmt.Errorf("%w %d for list.Index", ErrIndexOutOfRange, i))
		}
		return nil
	case reflect.String:
		s := rv.String()
		if i >= 0 && i < len(s) {
			return s[i : i+1]
		}
		if m.Reporter != nil {
			m.Reporter.Report(fmt.Errorf("%w %d for list.Index", ErrIndexOutOfRange, i))
		}
		return ""
	}
	if m.Reporter != nil {
		m.Reporter.Report(fmt.Errorf("%w %T for list.Index", ErrInvalidType, v))
	}
	return nil
}

func (m *ListModule) Slice(start, end int, v any) any {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		if start < 0 {
			start = rv.Len() + start
		}
		if end < 0 {
			end = rv.Len() + end
		}
		if start < 0 {
			start = 0
		}
		if end > rv.Len() {
			end = rv.Len()
		}
		if start >= end {
			return reflect.Zero(rv.Type()).Interface()
		}
		return rv.Slice(start, end).Interface()
	case reflect.String:
		s := rv.String()
		if start < 0 {
			start = len(s) + start
		}
		if end < 0 {
			end = len(s) + end
		}
		if start < 0 {
			start = 0
		}
		if end > len(s) {
			end = len(s)
		}
		if start >= end {
			return ""
		}
		return s[start:end]
	}
	if m.Reporter != nil {
		m.Reporter.Report(fmt.Errorf("%w %T for list.Slice", ErrInvalidType, v))
	}
	return nil
}

func (m *ListModule) Reverse(v any) any {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		length := rv.Len()
		result := reflect.MakeSlice(rv.Type(), length, length)
		for i := 0; i < length; i++ {
			result.Index(i).Set(rv.Index(length - i - 1))
		}
		return result.Interface()
	case reflect.String:
		runes := []rune(rv.String())
		for i, ch := range runes[:len(runes)/2] {
			idx := len(runes) - 1 - i
			runes[i] = runes[idx]
			runes[idx] = ch
		}
		return string(runes)
	}
	if m.Reporter != nil {
		m.Reporter.Report(fmt.Errorf("%w %T for list.Reverse", ErrInvalidType, v))
	}
	return nil
}
