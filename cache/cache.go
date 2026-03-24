package cache

import (
	"context"
	"sync"
)

type Cache[K comparable, V any] interface {
	Get(context.Context, K) (V, bool, error)
	Set(context.Context, K, V) error
}

type MapCache[K comparable, V any] struct {
	m map[K]V
}

func NewMapCache[K comparable, V any]() MapCache[K, V] {
	return MapCache[K, V]{
		m: make(map[K]V),
	}
}

func (m MapCache[K, V]) Get(ctx context.Context, k K) (v V, ok bool, err error) {
	v, ok = m.m[k]
	return
}
func (m MapCache[K, V]) Set(_ context.Context, k K, v V) error {
	m.m[k] = v
	return nil
}

type SyncMapCache[K comparable, V any] struct {
	syncMap sync.Map
}

func (l *SyncMapCache[K, V]) Get(ctx context.Context, k K) (v V, ok bool, err error) {
	var a any
	if a, ok = l.syncMap.Load(k); ok {
		v = a.(V)
		return
	}
	return
}

func (l *SyncMapCache[K, V]) Set(_ context.Context, k K, v V) error {
	l.syncMap.Store(k, v)
	return nil
}

type emptyCache[K comparable, V any] struct{}

func (e emptyCache[K, V]) Get(_ context.Context, _ K) (val V, ok bool, err error) { return }
func (e emptyCache[K, V]) Set(_ context.Context, _ K, _ V) (err error)            { return }

type tiersCache[K comparable, V any] struct {
	c Cache[K, V]
	p Cache[K, V]
}

func WithCache[K comparable, V any](p, c Cache[K, V]) Cache[K, V] {
	return &tiersCache[K, V]{
		c: &SyncMapCache[K, V]{},
		p: emptyCache[K, V]{},
	}
}

func (t tiersCache[K, V]) Get(ctx context.Context, key K) (v V, ok bool, err error) {
	if v, ok, err = t.c.Get(ctx, key); ok {
		return
	}
	if v, ok, err = t.p.Get(ctx, key); ok {
		err = t.c.Set(ctx, key, v)
	}
	return
}

func (t tiersCache[K, V]) Set(ctx context.Context, k K, v V) error {
	return nil
}
