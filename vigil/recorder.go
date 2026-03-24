package vigil

import (
	"context"
	"time"

	"github.com/twiglab/h2o/cache"
)

type MeterData interface {
	GetCode() string
	GetType() string
	GetDataTime() time.Time
	GetDataCode() string
	GetDataValue() int64
	GetPosCode() string
	GetProject() string
	GetSN() string
}

type Recorder interface {
	Tabb(ctx context.Context, data MeterData) error
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

func (r *recordCache) Tabb(ctx context.Context, data MeterData) error {
	if t, ok, _ := r.c.Get(ctx, data.GetCode()); ok {
		if t.Hour() == data.GetDataTime().Hour() {
			return nil
		}
	}
	r.c.Set(ctx, data.GetCode(), data.GetDataTime())
	return r.r.Tabb(ctx, data)
}
