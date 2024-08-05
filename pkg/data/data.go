package data

import "fmt"

type Quantity int64

func (q Quantity) String() string {
	return toDataUnit(int64(q))
}

func toDataUnit(units int64) string {
	if units < 1024 {
		return fmt.Sprintf("%d B", units)
	}
	if units < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(units)/1024)
	}
	if units < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(units)/1024/1024)
	}
	if units < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(units)/1024/1024/1024)
	}
	return fmt.Sprintf("%.2f TB", float64(units)/1024/1024/1024/1024)
}
