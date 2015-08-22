package util

import (
	"regexp"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	slugizeSep = "-"
)

var (
	rControl   = regexp.MustCompile("[\u0000-\u001f]")
	rSpecial   = regexp.MustCompile("[\\s~`!@#\\$%\\^&\\*\\(\\)\\-_\\+=\\[\\]\\{\\}\\|\\;:\"'<>,\\.\\?\\/]+")
	rRepeatSep = regexp.MustCompile(slugizeSep + "{2,}")
	rEdgeSep   = regexp.MustCompile("^" + slugizeSep + "+|" + slugizeSep + "+$")
)

func Slugize(str string) string {
	str = escapeDiacritic(str)
	str = rControl.ReplaceAllString(str, slugizeSep)
	str = rSpecial.ReplaceAllString(str, slugizeSep)
	str = rRepeatSep.ReplaceAllString(str, slugizeSep)
	str = rEdgeSep.ReplaceAllString(str, "")

	return str
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func escapeDiacritic(str string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, str)
	return result
}
