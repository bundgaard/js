package js

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"unicode"
)

///////////////////////////////////////////////////////////////////////////////

type Pos struct {
	Line   int
	Column int
}

const EofRune = -1

type scanner struct {
	rd           io.RuneReader
	buf          bytes.Buffer
	peeking      bool
	peekRune     rune
	last         rune
	line, column uint
}

func (s *scanner) read() rune {
	if s.peeking {
		s.peeking = false
		return s.peekRune
	}
	return s.readChar()
}
func (s *scanner) readChar() rune {
	r, _, err := s.rd.ReadRune()
	if err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr)
		}
		r = EofRune
	}
	s.last = r
	return r
}

func (s *scanner) peek() rune {
	if s.peeking {
		return s.peekRune
	}
	r := s.read()
	s.peeking = true
	s.peekRune = r
	return r
}

func (s *scanner) back(r rune) {
	s.peeking = true
	s.peekRune = r
}

func (s *scanner) accum(r rune, valid func(rune) bool) {
	s.buf.Reset()
	for {
		s.buf.WriteRune(r)
		r = s.read()
		if r == EofRune {
			return
		}

		if !valid(r) {
			s.back(r)
			return
		}
	}
}

func (s *scanner) NextToken() *Token {
	for {
		r := s.read()

		switch {
		case isSpace(r):
		case r == '=':
			return newToken(Assign, "=")
		case r == EofRune:
			return newToken(EOF, "EOF")
		case r == ';':
			return newToken(Semi, ";")
		case r == '.':
			return newToken(Dot, ".")
		case r == ',':
			return newToken(Comma, ",")
		case r == '"' || r == '\'':
			return newToken(String, s.readString(r))
		case r == '+':
			return newToken(Add, "+")
		case r == '-':
			return newToken(Sub, "-")
		case r == '/':
			pr := s.peek()
			if pr == '/' {
				// read to newline or EofRune
				for s.last != '\n' && s.last != EofRune {
					s.read()
					s.peeking = false
				}
				continue
			} else if pr == '*' {
				// read to */

				for {
					if s.last == '*' && s.peek() == '/' {
						s.read()
						break
					}
					s.read()
				}
				continue

			}
			return newToken(Div, "/")
		case r == '*':
			return newToken(Mul, "*")
		case r == '(':
			return newToken(OpenParen, "(")
		case r == ')':
			return newToken(OpenParen, ")")
		case r == '{':
			return newToken(OpenCurly, "{")
		case r == '}':
			return newToken(CloseCurly, "}")
		case r == '[':
			return newToken(OpenBracket, "[")
		case r == ']':
			return newToken(CloseBracket, "]")
		default:
			token := new(Token)
			if isLetter(r) {
				name := s.readName()
				v, ok := keywords[name]
				if ok {
					token.Type = v
					token.Value = name
				} else {

					token.Type = Ident
					token.Value = name
				}

				return token
			} else if isDigit(r) {
				token.Type = Number
				token.Value = s.readLiteral()
			} else {
				token.Type = Illegal
				token.Value = string(s.last)

			}
			return token
		}
	}

}

func (s *scanner) readString(quote rune) string {
	s.accum(quote, isAlphanum)

	// check if last is quote
	if r := s.peek(); r != quote || r != EofRune {
		fmt.Fprintf(os.Stderr, "invalid token after %s", &s.buf)
		os.Exit(1)
	}

	return s.buf.String()
}

func (s *scanner) readLiteral() string {
	s.accum(s.last, isDigit)
	return s.buf.String()
}
func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func (s *scanner) readName() string {
	s.accum(s.last, isAlphanum)
	return s.buf.String()
}
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_'
}

func NewScanner(rd io.RuneReader) *scanner {
	s := &scanner{rd: rd, line: 1, column: 1}
	s.readChar()
	return s
}

func NewScannerFromFile(fp string) *scanner {

	buf, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	return NewScanner(bytes.NewReader(buf))
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAlphanum(r rune) bool {
	return r == '_' || isNumber(r) || unicode.IsLetter(r)
}
