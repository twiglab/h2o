package hkv

import (
	"context"
	"time"

	"github.com/twiglab/h2o/hank"
)

type Item struct {
	Data hank.MetaData
	Time time.Time
}

func (i Item) IsSince(d time.Duration) bool {
	return time.Since(i.Time) < d
}

type Cache struct {
	m        map[string]Item
	Duration time.Duration
}

func NewCache(d time.Duration) *Cache {
	return &Cache{
		m:        make(map[string]Item),
		Duration: d,
	}
}

func (c *Cache) Get(ctx context.Context, code string) (hank.MetaData, bool, error) {
	if i, ok := c.m[code]; ok {
		if i.IsSince(c.Duration) {
			return i.Data, ok, nil
		}
	}
	return hank.MetaData{}, false, nil
}

func (c *Cache) Set(ctx context.Context, code string, data hank.MetaData) error {
	c.m[code] = Item{Data: data, Time: time.Now()}
	return nil
}
