package object

import "fmt"

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BooleanType
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
