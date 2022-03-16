package eval

import (
	"fmt"
	"github.com/bundgaard/js/object"
	"os"
)

var builtins = map[string]*object.BuiltinObject{
	"println": &object.BuiltinObject{
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 1 {
				fmt.Fprintf(os.Stdout, "%v\n", args[0].Inspect())
			}

			return &object.NullObject{}
		},
	},
}
