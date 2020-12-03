package zcache

import (
	"context"
	"sync"
	"time"
)

type cachemap struct {
	ma sync.Map
}

func NewCache() *cachemap {
	m := cachemap{}
	return &m
}

//
func (m *cachemap) store(ctx context.Context, key, value interface{}) {
	valueCtx := context.WithValue(ctx, key, value)
	m.ma.Store(key, valueCtx)
}

func (m *cachemap) load(key interface{}) (interface{}, bool) {
	value, ok := m.ma.Load(key)
	if ok {
		ctx := value.(context.Context)
		select {
		case <-ctx.Done():
			m.ma.Delete(key)
			return nil, false
		default:
			return ctx.Value(key), true
		}
	}
	return nil, false
}

func (m *cachemap) delete(key interface{}) {
	m.ma.Delete(key)
}

//
func (m *cachemap) StoreWithTimeout(key, value interface{}, duration time.Duration) {
	ctxTimeout, _ := context.WithTimeout(context.Background(), duration)
	m.store(ctxTimeout, key, value)
}
func (m *cachemap) Store(key, value interface{}) {
	m.store(context.TODO(), key, value)
}

func (m *cachemap) NotFoundStore(key, value interface{}) bool {
	_, exist := m.load(key)
	if exist {
		return false
	}

	m.store(context.TODO(), key, value)
	return true
}

func (m *cachemap) StoreWithContext(ctx context.Context, key, value interface{}) {
	m.store(ctx, key, value)
}

func (m *cachemap) Load(key interface{}) (interface{}, bool) {
	return m.load(key)
}

func (m *cachemap) LoadAndDelete(key interface{}) (interface{}, bool) {
	value, loaded := m.load(key)
	if !loaded {
		return nil, false
	}
	m.delete(key)
	return value, true
}

func (m *cachemap) Delete(key interface{}) {
	m.delete(key)
}
