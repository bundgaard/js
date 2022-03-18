package object

import (
	"bytes"
	"fmt"
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/token"
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

//go:generate stringer -type ObjectType
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
	FunctionObj
)

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
	var (
		out   bytes.Buffer
		pairs []string
	)

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ArrayObject }
func (ao *Array) Inspect() string {
	var (
		out      strings.Builder
		elements []string
	)

	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// fn <name> (<parameters>) {
// <body>
// }

type Function struct {
	Token       *token.Token
	Name        string
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (fl *Function) Type() ObjectType { return FunctionObj }
func (fl *Function) Inspect() string {
	var out strings.Builder
	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.Token.Value)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outer
	return env
}
