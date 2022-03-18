package object

type NullObject struct {
}

func (no *NullObject) Type() Type      { return NullType }
func (no *NullObject) Inspect() string { return "null" }
