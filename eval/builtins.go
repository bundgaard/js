package eval

import "github.com/bundgaard/js/object"

var builtins = map[string]*object.BuiltinObject{
	"println": &object.BuiltinObject{
		Fn: func(args ...object.Object) object.Object {
			return &object.NullObject{}
		},
	},
}
