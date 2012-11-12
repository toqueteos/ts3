package ts3

import (
	"testing"
)

var testQuote = []struct {
	in, expected string
}{
	{"", ""},
	{" ", `\s`},
	{"Hello\a", "Hello\\a"},
	{"Hello buddy!", `Hello\sbuddy!`},
	{"id|1 id|2", `id\p1\sid\p2`},
	{"heyooo\n\r", "heyooo\\n\\r"},
	{"Lots\tof\vtabs", "Lots\\tof\\vtabs"},
	{"I've never used these: \b, \f", `I've\snever\sused\sthese:\s\b,\s\f`},
}

var testUnquote = []struct {
	in, expected string
}{
	{"", ""},
	{`\s`, " "},
	{"Hello\\a", "Hello\a"},
	{`Hello\sbuddy!`, "Hello buddy!"},
	{`id\p1\sid\p2`, "id|1 id|2"},
	{"heyooo\\n\\r", "heyooo\n\r"},
	{"Lots\\tof\\vtabs", "Lots\tof\vtabs"},
	{`I've\snever\sused\sthese:\s\b,\s\f`, "I've never used these: \b, \f"},
}

func TestQuote(t *testing.T) {
	for index, tt := range testQuote {
		output := Quote(tt.in)
		if output != tt.expected {
			t.Errorf("%d. Got %q, want %q", index+1, output, tt.expected)
		}
	}
}

func TestUnquote(t *testing.T) {
	for index, tt := range testUnquote {
		output := Unquote(tt.in)
		if output != tt.expected {
			t.Errorf("%d. Got %q, want %q", index+1, output, tt.expected)
		}
	}
}

func TestQuoteUnquote(t *testing.T) {
	for index, tt := range testQuote {
		output := Unquote(Quote(tt.in))
		if output != tt.in {
			t.Errorf("%d. Got %q, want %q", index+1, output, tt.in)
		}
	}
}
