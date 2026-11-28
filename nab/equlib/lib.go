package equlib

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go/v5"
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
		return v, fmt.Errorf("key %s: registered type %T does not implement required interface", key, a)
	}
	return v, nil
}

func TryDo(ctx context.Context, attempts uint, delay time.Duration, dofn func() error) error {
	return retry.New(
		retry.Attempts(attempts),
		retry.Delay(delay),
		retry.DelayType(retry.FixedDelay), //  固定间隔，不指数放大
		retry.Context(ctx),                // 支持取消
		retry.LastErrorOnly(true),         // 只返回最后一次错误，日志更干净
	).Do(dofn)
}
