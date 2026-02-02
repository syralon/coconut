package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type Option func(c *LRUCache)

func WithMaxSize(size int) Option {
	return func(c *LRUCache) {
		c.size = size
	}
}

func WithMaxAge(age time.Duration) Option {
	return func(c *LRUCache) {
		c.age = age
	}
}

type element struct {
	key, data any
	exp       time.Time
}

type LRUCache struct {
	list  list.List
	index map[any]*list.Element
	lock  sync.RWMutex
	age   time.Duration
	size  int
}

func NewLRUCache(opts ...Option) *LRUCache {
	c := &LRUCache{index: make(map[any]*list.Element)}
	for _, opt := range opts {
		opt(c)
	}
	if c.age == 0 && c.size == 0 {
		panic(errors.New("an unlimited lru cache is created with no expiration and no max size"))
	}
	return c
}

func (c *LRUCache) Put(key any, value any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	must := c.size > 0 && c.list.Len() >= c.size
	c.pop(c.list.Front(), must)

	var exp time.Time
	if c.age > 0 {
		exp = time.Now().Add(c.age)
	}
	data := &element{key: key, data: value, exp: exp}
	if e, ok := c.index[key]; ok {
		e.Value = data
		c.list.MoveToBack(e)
		return
	}
	node := c.list.PushBack(data)
	c.index[key] = node
	return
}

func (c *LRUCache) Get(key any) (any, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	e, ok := c.index[key]
	if !ok {
		return nil, false
	}
	v := e.Value.(*element)

	// remove expired element
	if !v.exp.IsZero() && v.exp.Before(time.Now()) {
		c.pop(e, true)
		return nil, false
	}

	// reset expiration
	if c.age > 0 {
		v.exp = time.Now().Add(c.age)
	}
	c.list.MoveToBack(e)

	return v.data, true
}

func (c *LRUCache) GetAnyway(key any) any {
	val, _ := c.Get(key)
	return val
}

func (c *LRUCache) Remove(key any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	e := c.index[key]
	delete(c.index, key)
	if e != nil {
		c.list.Remove(e)
	}
}

func (c *LRUCache) pop(e *list.Element, must bool) {
	if e == nil {
		return
	}
	v := e.Value.(*element)
	if !must && (v.exp.IsZero() || v.exp.After(time.Now())) {
		return
	}
	c.list.Remove(e)
	delete(c.index, v.key)
}

type LRUCacheT[K, V any] struct {
	*LRUCache
}

func NewLRUCacheT[K, V any](opts ...Option) *LRUCacheT[K, V] {
	return &LRUCacheT[K, V]{LRUCache: NewLRUCache(opts...)}
}

func (c *LRUCacheT[K, V]) Put(key K, value V) {
	c.LRUCache.Put(key, value)
}

func (c *LRUCacheT[K, V]) Get(key K) (val V, ok bool) {
	v, ok := c.LRUCache.Get(key)
	if !ok {
		return val, false
	}
	return v.(V), true
}

func (c *LRUCacheT[K, V]) GetAnyway(key K) V {
	v, _ := c.Get(key)
	return v
}

func (c *LRUCacheT[K, V]) Remove(key K) {
	c.LRUCache.Remove(key)
}
