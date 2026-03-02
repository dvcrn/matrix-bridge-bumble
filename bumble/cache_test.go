package bumble

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_Set(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1", time.Second)
	cache.Set("key2", 42, time.Minute)

	assert.Equal(t, 2, len(cache.items))
	assert.Equal(t, "value1", cache.items["key1"].value)
	assert.Equal(t, 42, cache.items["key2"].value)
}

func TestCache_Get(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1", time.Second)
	cache.Set("key2", 42, time.Minute)

	value, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", value)

	value, ok = cache.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, 42, value)

	value, ok = cache.Get("key3")
	assert.False(t, ok)
	assert.Nil(t, value)
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1", time.Second)
	cache.Set("key2", 42, time.Minute)

	cache.Delete("key1")
	assert.Equal(t, 1, len(cache.items))
	assert.NotContains(t, cache.items, "key1")
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1", time.Millisecond)
	cache.Set("key2", 42, time.Minute)

	time.Sleep(time.Millisecond * 2)
	cache.Cleanup()

	assert.Equal(t, 1, len(cache.items))
	assert.NotContains(t, cache.items, "key1")
	assert.Contains(t, cache.items, "key2")
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1", time.Millisecond)

	value, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", value)

	time.Sleep(time.Millisecond * 2)

	value, ok = cache.Get("key1")
	assert.False(t, ok)
	assert.Nil(t, value)
}

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache()
	cache.Set("key1", "value1", time.Second)
	cache.Set("key2", 42, time.Minute)

	value, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", value)

	value, ok = cache.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, 42, value)
}
