package ts3

import (
	"strings"
)

// Escapes special chars
func ts3Quote(s string) string {
	s = strings.Replace(s, "/", "\\/", -1)
	s = strings.Replace(s, " ", "\\s", -1)
	s = strings.Replace(s, "|", "\\p", -1)
	return strings.Trim(s, "\r")
}

// Unescapes special chars
func ts3Unquote(s string) string {
	s = strings.Replace(s, "\\/", "/", -1)
	s = strings.Replace(s, "\\s", " ", -1)
	s = strings.Replace(s, "\\p", "|", -1)
	return strings.Trim(s, "\r")
}

// Keeps only printable ASCII runes, also cleans "\r"
func StringsTrimNet(s string) string {
	var res []rune

	s = strings.Trim(s, "\r")

	// Just pretty ASCII runes
	for _, r := range s {
		switch {
		case 32 >= r || r <= 127:
			res = append(res, r)
		}
	}

	return string(res)
}
