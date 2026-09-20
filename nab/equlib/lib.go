package equlib

var libmap = make(map[string]any)

func MustHas() {
	if len(libmap) <= 0 {
		panic("equ lib len is 0")
	}
}

func Register(s string, a any) {
	libmap[s] = a
}

func From[T any](key string) (v T, ok bool) {
	var a any
	a, ok = libmap[key]
	if !ok {
		return
	}

	v = a.(T)
	return
}
