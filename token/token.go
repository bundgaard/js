package token

//go:generate stringer -type TokenType
type TokenType uint8

const (
	_ TokenType = iota
	EOF
	Illegal

	Assign  // =
	Semi    // ;
	Dot     // .
	Comma   // ,
	Colon   // :
	Quote   // "
	SQuote  // '
	Ident   // identifier
	Literal // 1 1.0 "foo"
	String  // "foo"

	Add // +
	Sub // -
	Mul // *
	Div // /

	OpenParen    // (
	CloseParen   // )
	OpenBracket  // [
	CloseBracket // ]
	OpenCurly    // {
	CloseCurly   //	}
	CommentLine  // //
	CommentBlock // /* */

	Var
	Number
)

var Keywords = map[string]TokenType{
	"var": Var,
}

type Token struct {
	Type  TokenType
	Value string
}

func New(tokenType TokenType, value string) *Token {
	return &Token{Type: tokenType, Value: value}
}
