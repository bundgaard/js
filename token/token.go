package token

//go:generate stringer -type Type
type Type uint8

const (
	_ Type = iota
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
	Function
	Null
	True
	False
)

var Keywords = map[string]Type{
	"var":   Var,
	"fn":    Function,
	"null":  Null,
	"true":  True,
	"false": False,
}

type Token struct {
	Type  Type
	Value string
}

func New(tokenType Type, value string) *Token {
	return &Token{Type: tokenType, Value: value}
}
