package equlib

import "errors"

var libmap = make(map[string]any)

func MustHas() {
	if len(libmap) <= 0 {
		panic("equ lib len is 0")
	}
}

func Register(s string, a any) {
	libmap[s] = a
}

func From[T any](key string) (v T, err error) {
	var (
		a  any
		ok bool
	)

	a, ok = libmap[key]
	if !ok {
		err = errors.New("key no found: " + key)
		return
	}

	v = a.(T)
	return
}
