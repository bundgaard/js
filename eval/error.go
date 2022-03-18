package eval

import (
	"fmt"
	"github.com/bundgaard/js/object"
)

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorType
	}
	return false
}
func newError(format string, v ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, v...)}
}
