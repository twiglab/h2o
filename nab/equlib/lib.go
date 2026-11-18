package equlib

import "fmt"

var libmap = make(map[string]any)

func MustHas() {
	if len(libmap) <= 0 {
		panic("equ lib len is 0")
	}
}

func Register(s string, a any) { libmap[s] = a }

// RegisterAs 注册设备实现
// 配合 `var _ nab.DeviceCollector = DTZY{}` 之类的断言，让遗漏在编译期暴露。
// 后期用此方法取代Register
func RegisterAs[T any](s string, a T) { libmap[s] = a }

func From[T any](key string) (v T, err error) {
	a, ok := libmap[key]
	if !ok {
		return v, fmt.Errorf("key no found: %s", key)
	}
	v, ok = a.(T)
	if !ok {
		return v, fmt.Errorf("key %s does not implement %T", key, v)
	}
	return v, nil
}
