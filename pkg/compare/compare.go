package compare

import "strings"

func Lower(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

func ContainsLower(a, b string) bool {
	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}
