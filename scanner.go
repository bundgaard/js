package js

import (
	"io"
	"io/ioutil"
	"log"
	"unicode"
)

type Scanner struct {
	in           io.Reader
	buf          []byte
	position     int
	readPosition int
	ch           byte

	errh         func(line, col uint, msg string) // error handler
	ioerr        error                            // pending IO error
	line, column uint
}

func (s *Scanner) init(src io.Reader, errh func(line, col uint, msg string)) {
	s.in = src
	s.errh = errh

	if s.buf == nil {
		s.buf = make([]byte, nextSize(0))
	}

	s.ioerr = nil
	s.line = 0
	s.column = 0

}

func (s *Scanner) readChar() {
	if s.readPosition >= len(s.buf) {
		s.ch = 0
	} else {
		s.ch = s.buf[s.readPosition]
	}
	s.column++
	s.position = s.readPosition
	s.readPosition++
}

func nextSize(size int) int {
	const min = 4 << 10 // 4K minimum buffer size
	const max = 1 << 20 // 1M maximum buffer size which will be doubled

	if size < min {
		return min
	}
	if size <= max {
		return size << 1
	}

	return size + max
}

func (s *Scanner) peek() byte {
	var c byte
	if s.readPosition >= len(s.buf) {
		c = 0
	} else {
		c = s.buf[s.readPosition]
	}
	return c
}

func (s *Scanner) NextToken() *Token {
	token := new(Token)

	s.skipWhitespace()

	switch s.ch {
	case '=':
		token = newToken(Assign, s.ch)
	case ';':
		token = newToken(Semi, s.ch)
	case '.':

		token = newToken(Dot, s.ch)
	case ',':
		token = newToken(Comma, s.ch)
	case '"', '\'':
		token.Type = String
		token.Value = s.readString(s.ch)
	case '+':
		token = newToken(Add, s.ch) // TODO(dbundgaard) turn into bitset of operators...
	case '-':
		token = newToken(Sub, s.ch) // TODO(dbundgaard) turn into bitset of operators...
	case '*':
		token = newToken(Mul, s.ch) // TODO(dbundgaard) turn into bitset of operators...
	case '/':
		peek := s.peek()
		if peek == '/' {
			position := s.position
			s.readChar()
			// line comment ignore to end of line

			for s.ch != '\n' && s.ch != 0 {
				s.readChar()
			}
			token.Type = CommentLine
			token.Value = string(s.buf[position:s.position])
		} else if peek == '*' {
			position := s.position
			s.readChar()
			for {

				if s.ch == '*' && s.buf[s.readPosition] == '/' {
					s.readChar()

					break
				}
				s.readChar()
			}

			token.Type = CommentBlock
			v := s.buf[position : s.position+1]
			token.Value = string(v)
		} else {
			token = newToken(Div, s.ch) // TODO(dbundgaard) turn into bitset of operators...
		}

	case '(':
		token = newToken(OpenParen, s.ch)
	case ')':
		token = newToken(CloseParent, s.ch)
	case '[':
		token = newToken(OpenBracket, s.ch)
	case ']':
		token = newToken(CloseBracket, s.ch)
	case '{':
		token = newToken(OpenCurly, s.ch)
	case '}':
		token = newToken(CloseCurly, s.ch)
	case ':':
		token = newToken(Colon, s.ch)
	case 0:
		token.Value = ""
		token.Type = EOF
	default:
		if isLetter(s.ch) {
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
		} else if isDigit(s.ch) {
			token.Type = Number
			token.Value = s.readLiteral()
		} else {
			token.Type = Illegal
			token.Value = string(s.ch)

		}
		return token
	}

	s.readChar()
	return token
}

func (s *Scanner) readString(quote byte) string {
	position := s.position + 1
	for {

		if s.ch == '\\' && s.buf[s.readPosition] == quote {
			s.readChar()
			continue
		} /* else if s.ch == '\\' && s.input[s.readPosition] == 'u' {
			s.readChar()
			s.readChar()
		} */
		s.readChar()
		if s.ch == quote || s.ch == 0 {
			break
		}

	}

	return string(s.buf[position:s.position])

}
func (s *Scanner) readLiteral() string {
	position := s.position
	for isDigit(s.ch) {
		s.readChar()
	}
	v := s.buf[position:s.position]
	return string(v)

}
func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
func (s *Scanner) readName() string {
	position := s.position
	for isLetter(s.ch) || isDigit(s.ch) {
		s.readChar()
	}
	v := s.buf[position:s.position]
	return string(v)
}
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_'
}
func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(rune(s.ch)) {
		if s.ch == '\n' {
			s.column = 1
			s.line++
		}
		s.readChar()
	}
}

func NewScanner(input string) *Scanner {
	s := &Scanner{buf: []byte(input), line: 1, column: 1}
	s.readChar()
	return s
}

func NewScannerFromFile(fp string) *Scanner {
	buf, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	return NewScanner(string(buf))
}
