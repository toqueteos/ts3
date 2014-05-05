package ts3

// Escaping table:
// NAME         CHAR    ASCII   REPLACE CHAR    REPLACE ASCII
// Backslash    \       92      \\              92 92
// Slash        /       47      \/              92 47
// Whitespace   " "     32      \s              92 115
// Pipe         |       124     \p              92 112
// Bell         \a      7       \a              92 97
// Backspace    \b      8       \b              92 98
// Formfeed     \f      12      \f              92 102
// Newline      \n      10      \n              92 110
// Car. Ret     \r      13      \r              92 114
// Hor. Tab     \t      9       \t              92 116
// Ver. Tab     \v      11      \v              92 118

import (
	"strconv"
	"strings"
)

var (
	quoteTable = map[rune][]rune{
		92:   []rune{92, 92},  // \   -> \\
		'/':  []rune{92, 47},  // /   -> \/
		' ':  []rune{92, 115}, // " " -> \s
		'|':  []rune{92, 112}, // |   -> \p
		'\a': []rune{92, 97},  // All these translate to `/` + `letter`
		'\b': []rune{92, 98},  //
		'\f': []rune{92, 102}, //
		'\n': []rune{92, 110}, //
		'\r': []rune{92, 114}, //
		'\t': []rune{92, 116}, //
		'\v': []rune{92, 118}, //
	}

	unquoteTable = map[string]string{
		`\\`: `\`, // \\ -> \
		`\/`: `/`, // \/ -> /
		`\s`: ` `, // \s -> " "
		`\p`: `|`, // \p -> |
		`\a`: "\a",
		`\b`: "\b",
		`\f`: "\f",
		`\n`: "\n",
		`\r`: "\r",
		`\t`: "\t",
		`\v`: "\v",
	}
)

// Escapes special chars
func Quote(s string) string {
	var res = make([]rune, 0)

	for _, r := range []rune(s) {
		if v, ok := quoteTable[r]; ok {
			res = append(res, v...)
		} else {
			res = append(res, r)
		}
	}

	return string(res)
}

// Unescapes special chars
func Unquote(s string) string {
	for k, v := range unquoteTable {
		s = strings.Replace(s, k, v, -1)
	}

	return s
}

// Keeps only printable ASCII runes, also cleans "\r"
func trimNet(s string) string {
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

func parseError(data string) ErrorMsg {
	var err ErrorMsg
	split := strings.Split(data, " ")

	values := make([]string, 2)
	for i, s := range split[1:] {
		values[i] = strings.Split(s, "=")[1]
	}

	err.Id, _ = strconv.Atoi(values[0])
	err.Msg = values[1]

	return err
}
