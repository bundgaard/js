package scanner

import (
	"bytes"
	"fmt"
	"github.com/bundgaard/js/token"
	"io"
	"io/ioutil"
	"log"
	"os"
)

///////////////////////////////////////////////////////////////////////////////

type Pos struct {
	Line   int
	Column int
}

const EofRune = -1

type Scanner struct {
	rd           io.RuneReader
	Buf          bytes.Buffer
	peeking      bool
	peekRune     rune
	last         rune
	Line, Column uint
}

func (s *Scanner) read() rune {
	if s.peeking {
		s.peeking = false
		return s.peekRune
	}
	return s.readChar()
}
func (s *Scanner) readChar() rune {
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

func (s *Scanner) peek() rune {
	if s.peeking {
		return s.peekRune
	}
	r := s.read()
	s.peeking = true
	s.peekRune = r
	return r
}

func (s *Scanner) back(r rune) {
	s.peeking = true
	s.peekRune = r
}

func (s *Scanner) accum(r rune, valid func(rune) bool) {
	s.Buf.Reset()
	for {
		s.Buf.WriteRune(r)
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

func (s *Scanner) NextToken() *token.Token {
	for {
		r := s.read()
		switch {
		case isSpace(r):
		case r == ':':
			return token.New(token.Colon, ":")
		case r == '=':
			return token.New(token.Assign, "=")
		case r == EofRune:
			return token.New(token.EOF, "EOF")
		case r == ';':
			return token.New(token.Semi, ";")
		case r == '.':
			return token.New(token.Dot, ".")
		case r == ',':
			return token.New(token.Comma, ",")
		case r == '"' || r == '\'':
			return token.New(token.String, s.readString(r))
		case r == '+':
			return token.New(token.Add, "+")
		case r == '-':
			return token.New(token.Sub, "-")
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
			return token.New(token.Div, "/")
		case r == '*':
			return token.New(token.Mul, "*")
		case r == '(':
			return token.New(token.OpenParen, "(")
		case r == ')':
			return token.New(token.CloseParen, ")")
		case r == '{':
			return token.New(token.OpenCurly, "{")
		case r == '}':
			return token.New(token.CloseCurly, "}")
		case r == '[':
			return token.New(token.OpenBracket, "[")
		case r == ']':
			return token.New(token.CloseBracket, "]")
		default:
			tk := new(token.Token)
			if isLetter(r) {
				name := s.readName()
				v, ok := token.Keywords[name]
				if ok {
					tk.Type = v
					tk.Value = name
				} else {
					tk.Type = token.Ident
					tk.Value = name
				}
				return tk
			} else if isDigit(r) {
				tk.Type = token.Number
				tk.Value = s.readLiteral()
			} else {
				tk.Type = token.Illegal
				tk.Value = string(s.last)

			}
			return tk
		}
	}

}

func (s *Scanner) readString(quote rune) string {
	s.Buf.Reset()

	var escaped bool
	// "foobar"
	// "\"foobar\""
	for {
		// TODO: fix so we read escape characters \<?>
		if !escaped && s.peek() == quote {
			break
		}

		// Are we escaped ?
		// Escaped characters would be \\, \' \", \/
		if s.peek() == '\\' {
			r := s.read() // consume slash
			if s.peek() == quote {
				s.Buf.WriteRune(r)

				escaped = !escaped
			}
		}
		r := s.read()
		s.Buf.WriteRune(r)
	}
	s.read() // eat quote

	return s.Buf.String()
}

func (s *Scanner) readLiteral() string {
	s.accum(s.last, isDigit)
	return s.Buf.String()
}

func (s *Scanner) readName() string {
	s.accum(s.last, isAlphaNum)
	return s.Buf.String()
}

func New(rd io.RuneReader) *Scanner {
	s := &Scanner{rd: rd, Line: 1, Column: 1}
	return s
}

func NewScannerFromFile(fp string) *Scanner {

	buf, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	return New(bytes.NewReader(buf))
}
