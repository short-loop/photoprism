package buffer

import (
	"sync"
	"sync/atomic"
)

type BufferMap struct {
	m   sync.Map
	len int64
	sync.RWMutex
}

func NewBufferMap() *BufferMap {
	return &BufferMap{
		m: sync.Map{},
	}
}

func (bm *BufferMap) Get(key ApiBufferKey) (value interface{}, ok bool) {
	bm.RLock()
	defer bm.RUnlock()
	bm.m.Range(func(k, v interface{}) bool {
		if key.Equals(k) {
			value = v
			ok = true
			return false
		}
		return true
	})
	return value, ok
}

func (bm *BufferMap) GetOrUpdate(key ApiBufferKey, updateWith func() interface{}) (actual interface{}, loaded bool) {
	v, loaded := bm.Get(key)
	if loaded {
		return v, loaded
	}
	value := updateWith()
	bm.Lock()
	defer bm.Unlock()
	bm.m.Store(key, value)
	bm.len++
	return value, loaded
}

func (bm *BufferMap) Delete(key interface{}) {
	bm.Lock()
	defer bm.Unlock()
	bm.len--
	bm.m.Delete(key)
}

func (bm *BufferMap) Range(f func(key interface{}, value interface{}) bool) {
	bm.m.Range(func(k, v interface{}) bool {
		return f(k, v)
	})
}

func (bm *BufferMap) Len() int64 {
	return atomic.LoadInt64(&bm.len)
}

func (bm *BufferMap) getKeys() []interface{} {
	bm.RLock()
	defer bm.RUnlock()
	keys := make([]interface{}, 0)
	bm.Range(func(key interface{}, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}
