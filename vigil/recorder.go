package vigil

import (
	"context"
	"time"

	"github.com/twiglab/h2o/cache"
)

type Recorder interface {
	TabbElecty(ctx context.Context, data ElectricityMeter) error
}

type recordCache struct {
	r Recorder
	c cache.Cache[string, time.Time]
}

func WithRecorder(r Recorder) Recorder {
	return &recordCache{
		r: r,
		c: cache.NewMapCache[string, time.Time](),
	}
}

func (r *recordCache) TabbElecty(ctx context.Context, data ElectricityMeter) error {
	if t, ok, _ := r.c.Get(ctx, data.Code); ok {
		if t.Hour() == data.DataTime.Hour() {
			return nil
		}
	}
	r.c.Set(ctx, data.Code, data.DataTime)
	return r.r.TabbElecty(ctx, data)
}
