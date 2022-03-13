package js

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
//                              OBJECT SYSTEM
///////////////////////////////////////////////////////////////////////////////

type Object interface {
	Type() ObjectType
	Inspect() string
}

type ObjectType uint8

const (
	_ ObjectType = iota
	NullObj
	ErrorObject
	ReturnValueObject
	IntegerObject
	StringObj
	ArrayObject
	HashObject
	NumberObj
	BuiltinObj
)

type NullObject struct {
}

func (no *NullObject) Type() ObjectType { return NullObj }
func (no *NullObject) Inspect() string  { return "null" }

type BuiltinFunction func(args ...Object) Object
type BuiltinObject struct {
	Fn BuiltinFunction
}

func (b *BuiltinObject) Type() ObjectType { return BuiltinObj }
func (b *BuiltinObject) Inspect() string  { return "builtin function" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return ReturnValueObject }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return IntegerObject }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type StringObject struct {
	Value string
}

func (s *StringObject) Type() ObjectType { return StringObj }
func (s *StringObject) Inspect() string  { return s.Value }
func (s *StringObject) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type NumberObject struct {
	Value int64
}

func (n *NumberObject) Type() ObjectType { return NumberObj }
func (n *NumberObject) Inspect() string  { return fmt.Sprintf("%d", n.Value) }
func (n *NumberObject) HashKey() HashKey {
	h := fnv.New64a()
	fmt.Fprintf(h, "%d", n.Value)
	return HashKey{Type: n.Type(), Value: h.Sum64()}
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ErrorObject }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type HashPair struct {
	Key   Object
	Value Object
}

type Hashable interface {
	HashKey() HashKey
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HashObject }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
