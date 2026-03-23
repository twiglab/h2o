package vigil

import (
	"context"
	"time"

	"github.com/twiglab/h2o/cache"
)

type RecordCache struct {
	Recorder Recorder
	c        cache.Cache[string, time.Time]
}

func (r *RecordCache) Tabb(ctx context.Context, data ElectricityMeter) error {
	if t, ok, _ := r.c.Get(ctx, data.Code); ok {
		if t.Hour() == data.DataTime.Hour() {
			return nil
		}
	}
	r.c.Set(ctx, data.Code, data.DataTime)
	return r.Recorder.Tabb(ctx, data)
}
