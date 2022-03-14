package eval

import (
	"fmt"
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/object"
	"log"
)

func Eval(n ast.Node, environment map[string]object.Object) object.Object {
	switch v := n.(type) {
	case *ast.Program:
		return evalProgram(v, environment)
	case *ast.ExpressionStatement:
		return Eval(v.Expression, environment)

	case *ast.VariableStatement:
		value := Eval(v.Value, environment)
		if value != nil {
			if value.Type() == object.ErrorObject {
				return value
			}
		}

		environment[v.Name.Value] = value
	// expressions
	case *ast.CallExpression:

	case *ast.InfixExpression:
		left := Eval(v.Left, environment)
		if left != nil {
			if left.Type() == object.ErrorObject {
				return left
			}
		}

		right := Eval(v.Right, environment)
		if right != nil {
			if right.Type() == object.ErrorObject {
				return right
			}
		}

		return evalInfixExpression(v.Operator, left, right, environment)

	case *ast.Identifier:
		return evalIdentifier(v, environment)
	case *ast.NumberLiteral:
		return &object.NumberObject{Value: v.Value}
	case *ast.StringLiteral:
		return &object.StringObject{Value: v.Value}
	case *ast.IndexExpression:
		left := Eval(v.Left, environment)
		if left != nil {
			if left.Type() == object.ErrorObject {
				return left
			}
		}

		index := Eval(v.Index, environment)
		if index != nil {
			if index.Type() == object.ErrorObject {
				return index
			}
		}

		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		pairs := make(map[object.HashKey]object.HashPair)
		for kn, vn := range v.Pairs {
			key := Eval(kn, environment)
			if key != nil {
				if key.Type() == object.ErrorObject {
					return key
				}
			}

			hkey, ok := key.(object.Hashable)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("unushable hash key: %v", key.Type())}
			}

			value := Eval(vn, environment)
			if value != nil {
				if value.Type() == object.ErrorObject {
					return value
				}
			}

			hashed := hkey.HashKey()
			pairs[hashed] = object.HashPair{Key: key, Value: value}
		}

		return &object.Hash{Pairs: pairs}

	default:
		log.Printf("eval unhandled type %T", v)
	}

	return nil

}

func evalIndexExpression(left, index object.Object) object.Object {
	return &object.Error{Message: "skipping index for now"}
}

func evalStringInfixExpression(operator string, left, right object.Object, env map[string]object.Object) object.Object {
	leftVal := left.(*object.StringObject).Value
	rightVal := right.(*object.StringObject).Value
	return &object.StringObject{
		Value: leftVal + rightVal,
	}
}

func evalNumberInfixExpression(operator string, left, right object.Object, env map[string]object.Object) object.Object {
	leftVal := left.(*object.NumberObject).Value
	rightVal := right.(*object.NumberObject).Value
	switch operator {
	case "+":
		return &object.NumberObject{Value: leftVal + rightVal}
	case "-":
		return &object.NumberObject{Value: leftVal - rightVal}
	case "/":
		return &object.NumberObject{Value: leftVal / rightVal}
	case "*":
		return &object.NumberObject{Value: leftVal * rightVal}
	}
	return nil
}

func evalInfixExpression(operator string, left, right object.Object, env map[string]object.Object) object.Object {

	switch {
	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalStringInfixExpression(operator, left, right, env)
	case left.Type() == object.NumberObj && right.Type() == object.NumberObj:
		return evalNumberInfixExpression(operator, left, right, env)
	case left.Type() != right.Type():
		return &object.Error{Message: fmt.Sprintf("type mismatch %v %s %v", left.Type(), operator, right.Type())}
	default:
		return &object.Error{Message: "unknown operator"}
	}
}

func evalIdentifier(n *ast.Identifier, env map[string]object.Object) object.Object {

	if val, ok := env[n.Value]; ok {
		return val
	}
	if builtin, ok := builtins[n.Value]; ok {
		return builtin
	}
	return &object.Error{Message: fmt.Sprintf("identifier not found %q", n.Value)}
}

func evalProgram(n *ast.Program, environment map[string]object.Object) object.Object {
	var result object.Object

	for _, statement := range n.Statements {

		result = Eval(statement, environment)

		switch v := result.(type) {
		case *object.ReturnValue:
			return v.Value
		case *object.Error:
			fmt.Printf("%T %v\n", v, v)
			continue
		}
	}

	return result

}
