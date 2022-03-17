package eval

import (
	"fmt"
	"github.com/bundgaard/js/ast"
	"github.com/bundgaard/js/object"
	"log"
)

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.BuiltinObject:
		return fn.Fn(args...)
	default:
		return newError("not function: %q", fn.Type())
	}
}

func evalExpressions(exps []ast.Expression, environment *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, environment)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result

}
func Eval(n ast.Node, environment *object.Environment) object.Object {
	// log.Printf("Eval %T %v", n, n)
	switch v := n.(type) {
	case *ast.Program:
		return evalProgram(v, environment)
	case *ast.ExpressionStatement:
		return Eval(v.Expression, environment)

	case *ast.VariableStatement:
		value := Eval(v.Value, environment)
		if value != nil && isError(value) {
			return value
		}

		environment.Set(v.Name.Value, value)

	// expressions
	case *ast.CallExpression:
		fn := Eval(v.Function, environment)
		if isError(fn) {
			return fn
		}

		args := evalExpressions(v.Arguments, environment)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)

	case *ast.InfixExpression:
		left := Eval(v.Left, environment)
		if left != nil {
			if isError(left) {
				return left
			}
		}

		right := Eval(v.Right, environment)
		if right != nil {
			if isError(right) {
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
		return evalIndexExpression(v, environment)
	case *ast.HashLiteral:
		return evalHashLiteral(v, environment)
	case *ast.ArrayLiteral:
		elements := evalExpressions(v.Elements, environment)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	default:
		log.Printf("eval unhandled type %T %v", v, v)
	}

	return nil

}

func evalHashLiteral(hashLiteral *ast.HashLiteral, environment *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for kn, vn := range hashLiteral.Pairs {
		key := Eval(kn, environment)
		if key != nil && isError(key) {
			return key
		}

		hkey, ok := key.(object.Hashable)
		if !ok {
			return &object.Error{Message: fmt.Sprintf("unushable hash key: %q", key.Type())}
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
}

func evalIndexExpression(v *ast.IndexExpression, environment *object.Environment) object.Object {
	left := Eval(v.Left, environment)
	if left != nil && isError(left) {
		return left
	}

	index := Eval(v.Index, environment)
	if index != nil && isError(index) {
		return index
	}

	return &object.Error{Message: "skipping index for now"}
}

func evalStringInfixExpression(operator string, left, right object.Object, env *object.Environment) object.Object {
	leftVal := left.(*object.StringObject).Value
	rightVal := right.(*object.StringObject).Value
	return &object.StringObject{
		Value: leftVal + rightVal,
	}
}

func evalNumberInfixExpression(operator string, left, right object.Object, env *object.Environment) object.Object {
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

func evalInfixExpression(operator string, left, right object.Object, env *object.Environment) object.Object {

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

func evalIdentifier(n *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(n.Value); ok {
		return val
	}
	if builtin, ok := builtins[n.Value]; ok {
		return builtin
	}
	return newError("identifier %q not found", n.Value)
}

func evalProgram(n *ast.Program, environment *object.Environment) object.Object {
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
