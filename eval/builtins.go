package eval

import (
	"fmt"
	"github.com/bundgaard/js/object"
	"os"
)

var builtins = map[string]*object.BuiltinObject{
	"println": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 1 {
				fmt.Fprintf(os.Stdout, "%v\n", args[0].Inspect())
			}

			return &object.NullObject{}
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d, want 1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.StringObject:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to %q not supported, got %q", "len", args[0].Type())
			}
		},
	},
}
