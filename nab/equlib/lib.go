package equlib

import (
	"fmt"
)

var libmap = make(map[string]any)

func MustHas() {
	if len(libmap) <= 0 {
		panic("equ lib len is 0")
	}
}

func Register(s string, a any) { libmap[s] = a }

func From[T any](key string) (v T, err error) {
	a, ok := libmap[key]
	if !ok {
		return v, fmt.Errorf("key no found: %s", key)
	}
	v, ok = a.(T)
	if !ok {
		return v, fmt.Errorf("key %s: registered type %T does not implement %T", key, a, v)
	}
	return v, nil
}
