package linepay

import "strconv"

func Bool(v bool) *bool { return &v }

func Int(v int) *int { return &v }

func Int64(v int64) *int64 { return &v }

func String(v string) *string { return &v }

func ParseInt64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

func MustParseInt64(v string) int64 {
	i, _ := strconv.ParseInt(v, 10, 64)
	return i
}
