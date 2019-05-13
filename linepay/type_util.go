package linepay

import "strconv"

// Bool bool to pointer function
func Bool(v bool) *bool { return &v }

// Int int to pointer function
func Int(v int) *int { return &v }

// Int64 int64 to pointer function
func Int64(v int64) *int64 { return &v }

// String string to pointer function
func String(v string) *string { return &v }

// ParseInt64 string to int64 function
func ParseInt64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

// MustParseInt64 string to int64 without error function
func MustParseInt64(v string) int64 {
	i, _ := strconv.ParseInt(v, 10, 64)
	return i
}
