package object

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return ReturnValueObject }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
