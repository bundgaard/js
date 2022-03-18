package object

type BuiltinFunction func(args ...Object) Object
type BuiltinObject struct {
	Fn BuiltinFunction
}

func (b *BuiltinObject) Type() Type      { return BuiltinType }
func (b *BuiltinObject) Inspect() string { return "builtin function" }
