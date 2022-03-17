package object

type Environment struct {
	store map[string]Object
	Outer *Environment
}

func NewEnvironment() *Environment {
	e := make(map[string]Object)
	return &Environment{store: e}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) ForEach(iterator func(key string, value Object)) {
	for k, v := range e.store {
		iterator(k, v)
	}
}
