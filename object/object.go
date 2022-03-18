package object

///////////////////////////////////////////////////////////////////////////////
//                              OBJECT SYSTEM
///////////////////////////////////////////////////////////////////////////////

type Object interface {
	Type() Type
	Inspect() string
}

//go:generate stringer -type ObjectType
type Type uint8

const (
	_ Type = iota
	NullType
	ErrorType
	ReturnValueType
	IntegerType
	StringType
	ArrayType
	HashType
	NumberType
	BuiltinType
	FunctionType
	BooleanType
)
