package vigil

import (
	"context"
	"time"

	"github.com/twiglab/h2o/cache"
)

type Recorder interface {
	Tabb(ctx context.Context, code string, data Meter) error
}

type recordCache struct {
	r Recorder
	c cache.SyncMapCache[string, time.Time]
}

func WithRecorder(r Recorder) Recorder {
	return &recordCache{
		r: r,
	}
}

func (r *recordCache) Tabb(ctx context.Context, code string, data Meter) error {
	if t, ok, _ := r.c.Get(ctx, data.Code); ok {
		if t.Hour() == data.DataTime.Hour() {
			return nil
		}
	}
	if err := r.c.Set(ctx, data.Code, data.DataTime); err != nil {
		return err
	}
	return r.r.Tabb(ctx, data.Code, data)
}
