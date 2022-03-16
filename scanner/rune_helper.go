package scanner

import "unicode"

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_'
}
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAlphaNum(r rune) bool {
	return r == '_' || isNumber(r) || unicode.IsLetter(r)
}
