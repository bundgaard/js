package object

type NullObject struct {
}

func (no *NullObject) Type() ObjectType { return NullObj }
func (no *NullObject) Inspect() string  { return "null" }
