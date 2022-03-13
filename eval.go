package js

import (
	"fmt"
	"log"
)

var builtins = map[string]*BuiltinObject{
	"println": &BuiltinObject{
		Fn: func(args ...Object) Object {
			return &NullObject{}
		},
	},
}

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
	case *CallExpression:

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
	case *NumberLiteral:
		return &NumberObject{Value: v.Value}
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

func evalNumberInfixExpression(operator string, left, right Object, env map[string]Object) Object {
	leftVal := left.(*NumberObject).Value
	rightVal := right.(*NumberObject).Value
	switch operator {
	case "+":
		return &NumberObject{Value: leftVal + rightVal}
	case "-":
		return &NumberObject{Value: leftVal - rightVal}
	case "/":
		return &NumberObject{Value: leftVal / rightVal}
	case "*":
		return &NumberObject{Value: leftVal * rightVal}
	}
	return nil
}

func evalInfixExpression(operator string, left, right Object, env map[string]Object) Object {

	switch {
	case left.Type() == StringObj && right.Type() == StringObj:
		return evalStringInfixExpression(operator, left, right, env)
	case left.Type() == NumberObj && right.Type() == NumberObj:
		return evalNumberInfixExpression(operator, left, right, env)
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
	if builtin, ok := builtins[n.Value]; ok {
		return builtin
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
