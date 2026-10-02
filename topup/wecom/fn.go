package wecom

import (
	"fmt"
	"time"
)

func du(i int64) string {
	x := float64(i) * 0.01
	return fmt.Sprintf("%.2f", x)
}

func fmtTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.DateTime)
}

func fmtTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func sub(a, b int64) int64 {
	return a - b
}

func add(a, b int64) int64 {
	return a + b
}
