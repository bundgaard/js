package object

import "hash/fnv"

type StringObject struct {
	Value string
}

func (s *StringObject) Type() Type      { return StringType }
func (s *StringObject) Inspect() string { return s.Value }
func (s *StringObject) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
