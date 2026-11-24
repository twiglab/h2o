package equlib

import (
	"context"
	"time"

	"github.com/avast/retry-go/v5"
)

func TryDo(ctx context.Context, attempts uint, delay time.Duration, dofn func() error) error {
	return retry.New(
		retry.Attempts(attempts),
		retry.Delay(delay),
		retry.DelayType(retry.FixedDelay), //  固定间隔，不指数放大
		retry.Context(ctx),                // 支持取消
		retry.LastErrorOnly(true),         // 只返回最后一次错误，日志更干净
	).Do(dofn)
}
