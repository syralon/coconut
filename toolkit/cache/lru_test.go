package cache

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLruCacheMaxSize(t *testing.T) {
	c := NewLRUCache(WithMaxSize(2))
	c.Put("1", 1)
	c.Put("2", 2)
	c.Put("3", 3)
	assertValue(t, c, "1", nil, false)
	assertValue(t, c, "2", 2, true)
	assertValue(t, c, "3", 3, true)
	assert.Equal(t, 2, c.list.Len())

}

func TestLruCacheMaxAge(t *testing.T) {
	c := NewLRUCache(WithMaxAge(time.Second))
	c.Put("1", 1)
	assert.Equal(t, 1, c.GetAnyway("1"))
	time.Sleep(time.Second)
	_, ok := c.Get("1")
	assert.Equal(t, false, ok)
	assert.Equal(t, 0, c.list.Len())
}

func TestLruCacheMove(t *testing.T) {
	c := NewLRUCache(WithMaxSize(2))
	c.Put("1", 1)
	c.Put("2", 2)
	c.Get("1")
	c.Put("3", 3)
	assertValue(t, c, "1", 1, true)
	assertValue(t, c, "2", nil, false)
	assertValue(t, c, "3", 3, true)
	assert.Equal(t, 2, c.list.Len())
}

func assertValue(t *testing.T, c *LRUCache, key any, val any, exist bool) {
	t.Helper()
	v, ok := c.Get(key)
	assert.Equal(t, exist, ok)
	assert.Equal(t, val, v)
}

func BenchmarkLruCachePut(b *testing.B) {
	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = uuid.New().String()
	}
	c := NewLRUCache(WithMaxSize(1024))
	b.ResetTimer()

	var i int
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keys[i%1024]
			c.Put(key, key)
			i++
		}
	})
}

func BenchmarkLruCacheGet(b *testing.B) {
	c := NewLRUCache(WithMaxSize(1024))
	keys := make([]string, 1024)
	for i := range keys {
		key := uuid.New().String()
		keys[i] = key
		c.Put(key, key)
	}
	b.ResetTimer()
	var i int
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(keys[i%1024])
			i++
		}
	})
}

func BenchmarkLruCache(b *testing.B) {
	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = uuid.New().String()
	}
	c := NewLRUCache(WithMaxSize(1024))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := keys[i%1024]
			c.Put(key, key)
			i++
		}
	})
	b.RunParallel(func(pb *testing.PB) {
		j := 0
		for pb.Next() {
			key := keys[j%1024]
			c.Get(key)
			j++
		}
	})
}
