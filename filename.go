package rec

import "strings"

var (
	ifr = map[rune]bool{
		'/':  true,
		'<':  true,
		'>':  true,
		':':  true,
		'"':  true,
		'\\': true,
		'|':  true,
		'?':  true,
		'*':  true,
	}
)

func legalFilename(s string) string {
	rs := []rune(s)
	sb := strings.Builder{}
	for k := range rs {
		if !ifr[rs[k]] {
			_, _ = sb.WriteRune(rs[k])
		}
	}
	s = sb.String()
	if s == "" {
		s = "Untitled"
	}

	return s
}
