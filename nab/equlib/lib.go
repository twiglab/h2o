package equlib

var libmap = make(map[string]any)

func MustHas() {
	if Size() <= 0 {
		panic("equ lib is 0")
	}
}

func Size() int {
	return len(libmap)
}

func Register(s string, a any) {
	libmap[s] = a
}

func From[T any](key string) T {
	v := libmap[key]
	return v.(T)
}
