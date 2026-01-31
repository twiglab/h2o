package hank

import (
	"context"
	"sync"
)

type Cache[K comparable, V any] interface {
	Get(context.Context, K) (V, bool, error)
	Set(context.Context, K, V) error
	Clear(context.Context) error
	Forget(context.Context, K) (V, bool, error)
}

type SyncMapCache[K comparable, V any] struct {
	Cache[K, V]
	syncMap sync.Map
}

func (l *SyncMapCache[K, V]) Get(ctx context.Context, k K) (v V, ok bool, err error) {
	var a any
	if a, ok = l.syncMap.Load(k); ok {
		v = a.(V)
		return
	}

	v, ok, err = l.Cache.Get(ctx, k)
	if ok {
		err = l.Set(ctx, k, v)
	}

	return
}

func (l *SyncMapCache[K, V]) Set(_ context.Context, k K, v V) error {
	l.syncMap.Store(k, v)
	return nil
}

func (l *SyncMapCache[K, V]) Clear(_ context.Context) error {
	l.syncMap.Clear()
	return nil
}

func (l *SyncMapCache[K, V]) Forget(_ context.Context, k K) (v V, ok bool, err error) {
	var a any
	if a, ok = l.syncMap.LoadAndDelete(k); ok {
		v = a.(V)
	}
	return
}

type emptyCache[K comparable, V any] struct{}

func (e emptyCache[K, V]) Get(_ context.Context, _ K) (val V, ok bool, err error)    { return }
func (e emptyCache[K, V]) Set(_ context.Context, _ K, _ V) (err error)               { return }
func (e emptyCache[K, V]) Clear(_ context.Context) (err error)                       { return }
func (e emptyCache[K, V]) Forget(_ context.Context, _ K) (val V, ok bool, err error) { return }

type TiersCache[K comparable, V any] struct {
	local  Cache[K, V]
	second Cache[K, V]
}

func NewTiersCache[K comparable, V any]() *TiersCache[K, V] {
	return &TiersCache[K, V]{
		local:  &SyncMapCache[K, V]{},
		second: emptyCache[K, V]{},
	}
}

func (t *TiersCache[K, V]) WithLocal(local Cache[K, V]) *TiersCache[K, V] {
	t.local = local
	return t
}

func (t *TiersCache[K, V]) WithSecond(second Cache[K, V]) *TiersCache[K, V] {
	t.second = second
	return t
}

func (t *TiersCache[K, V]) Get(ctx context.Context, key K) (v V, ok bool, err error) {
	if v, ok, err = t.local.Get(ctx, key); ok {
		return
	}
	if v, ok, err = t.second.Get(ctx, key); ok {
		err = t.local.Set(ctx, key, v)
	}
	return
}

func (t *TiersCache[K, V]) Set(ctx context.Context, k K, v V) error {
	return t.local.Set(ctx, k, v)
}
