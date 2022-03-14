// Code generated by "stringer -type TokenType"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-1]
	_ = x[Illegal-2]
	_ = x[Assign-3]
	_ = x[Semi-4]
	_ = x[Dot-5]
	_ = x[Comma-6]
	_ = x[Colon-7]
	_ = x[Quote-8]
	_ = x[SQuote-9]
	_ = x[Ident-10]
	_ = x[Literal-11]
	_ = x[String-12]
	_ = x[Add-13]
	_ = x[Sub-14]
	_ = x[Mul-15]
	_ = x[Div-16]
	_ = x[OpenParen-17]
	_ = x[CloseParen-18]
	_ = x[OpenBracket-19]
	_ = x[CloseBracket-20]
	_ = x[OpenCurly-21]
	_ = x[CloseCurly-22]
	_ = x[CommentLine-23]
	_ = x[CommentBlock-24]
	_ = x[Var-25]
	_ = x[Number-26]
}

const _TokenType_name = "EOFIllegalAssignSemiDotCommaColonQuoteSQuoteIdentLiteralStringAddSubMulDivOpenParenCloseParentOpenBracketCloseBracketOpenCurlyCloseCurlyCommentLineCommentBlockVarNumber"

var _TokenType_index = [...]uint8{0, 3, 10, 16, 20, 23, 28, 33, 38, 44, 49, 56, 62, 65, 68, 71, 74, 83, 94, 105, 117, 126, 136, 147, 159, 162, 168}

func (i TokenType) String() string {
	i -= 1
	if i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}