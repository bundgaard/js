package js

import (
	"github.com/bundgaard/js/object"
	"testing"
)

func TestLibraryNew(t *testing.T) {
	data := `var x = 50;
var y = 100;
var z = x + y;
println(z);`

	_, environment := New(data)

	if _, ok := environment.Get("z"); !ok {
		t.Errorf("expected %q. got none", "z")
	}

	if x, ok := environment.Get("x"); ok && x.(*object.NumberObject).Value != 50 {
		t.Errorf("expected %q to be %q. got %q", "x", 50, x.Inspect())
	}
}
