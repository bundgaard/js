package js

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"log"
	"strings"
)

func Eval(n Node, environment map[string]Object) Object {
	switch v := n.(type) {
	case *Program:
		return evalProgram(v, environment)
	case *ExpressionStatement:
		return Eval(v.Expression, environment)

	case *VariableStatement:
		value := Eval(v.Value, environment)
		if value != nil {
			if value.Type() == ErrorObject {
				return value
			}
		}

		environment[v.Name.Value] = value
	// expressions

	case *InfixExpression:
		left := Eval(v.Left, environment)
		if left != nil {
			if left.Type() == ErrorObject {
				return left
			}
		}

		right := Eval(v.Right, environment)
		if right != nil {
			if right.Type() == ErrorObject {
				return right
			}
		}

		return evalInfixExpression(v.Operator, left, right, environment)

	case *Identifier:
		return evalIdentifier(v, environment)

	case *StringLiteral:
		return &StringObject{Value: v.Value}
	case *IndexExpression:
		left := Eval(v.Left, environment)
		if left != nil {
			if left.Type() == ErrorObject {
				return left
			}
		}

		index := Eval(v.Index, environment)
		if index != nil {
			if index.Type() == ErrorObject {
				return index
			}
		}

		return evalIndexExpression(left, index)

	case *HashLiteral:
		pairs := make(map[HashKey]HashPair)
		for kn, vn := range v.Pairs {
			key := Eval(kn, environment)
			if key != nil {
				if key.Type() == ErrorObject {
					return key
				}
			}

			hkey, ok := key.(Hashable)
			if !ok {
				return &Error{Message: fmt.Sprintf("unushable hash key: %v", key.Type())}
			}

			value := Eval(vn, environment)
			if value != nil {
				if value.Type() == ErrorObject {
					return value
				}
			}

			hashed := hkey.HashKey()
			pairs[hashed] = HashPair{Key: key, Value: value}
		}

		return &Hash{Pairs: pairs}

	default:
		log.Printf("eval unhandled type %T", v)
	}

	return nil

}

func evalIndexExpression(left, index Object) Object {
	return &Error{Message: "skipping index for now"}
}

func evalStringInfixExpression(operator string, left, right Object, env map[string]Object) Object {
	leftVal := left.(*StringObject).Value
	rightVal := right.(*StringObject).Value
	return &StringObject{
		Value: leftVal + rightVal,
	}
}

func evalInfixExpression(operator string, left, right Object, env map[string]Object) Object {

	switch {
	case left.Type() == StringObj && right.Type() == StringObj:
		return evalStringInfixExpression(operator, left, right, env)
	case left.Type() != right.Type():
		return &Error{Message: fmt.Sprintf("type mismatch %v %s %v", left.Type(), operator, right.Type())}
	default:
		return &Error{Message: "unknown operator"}
	}
}

func evalIdentifier(n *Identifier, env map[string]Object) Object {

	if val, ok := env[n.Value]; ok {
		return val
	}
	return &Error{Message: fmt.Sprintf("identifier not found %q", n.Value)}
}

func evalProgram(n *Program, environment map[string]Object) Object {
	var result Object

	for _, statement := range n.Statements {

		result = Eval(statement, environment)

		switch v := result.(type) {
		case *ReturnValue:
			return v.Value
		case *Error:
			fmt.Printf("%T %v\n", v, v)
			continue
		}
	}

	return result

}

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
	NullObject
	ErrorObject
	ReturnValueObject
	IntegerObject
	StringObj
	ArrayObject
	HashObject
)

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
