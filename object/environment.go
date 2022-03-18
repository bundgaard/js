package object

import "fmt"

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

func (e *Environment) GetString(name string) (string, error) {
	obj, ok := e.store[name]
	if !ok {
		return "", fmt.Errorf("%q is not found", name)
	}

	switch v := obj.(type) {
	case *StringObject:
		return v.Value, nil
	default:
		return "", fmt.Errorf("%q is not of type string", name)
	}
}

func (e *Environment) GetBool(name string) (bool, error) {
	obj, ok := e.store[name]
	if !ok {
		return false, fmt.Errorf("%q is not in environment", name)
	}

	switch v := obj.(type) {
	case *Boolean:
		return v.Value, nil
	default:
		return false, fmt.Errorf("%q was not a bool (%T)", name, v)
	}
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

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outer
	return env
}
