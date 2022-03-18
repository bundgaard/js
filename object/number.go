package object

import (
	"fmt"
	"hash/fnv"
)

type NumberObject struct {
	Value int64
}

func (n *NumberObject) Type() Type      { return NumberType }
func (n *NumberObject) Inspect() string { return fmt.Sprintf("%d", n.Value) }
func (n *NumberObject) HashKey() HashKey {
	h := fnv.New64a()
	fmt.Fprintf(h, "%d", n.Value)
	return HashKey{Type: n.Type(), Value: h.Sum64()}
}
