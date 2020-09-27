package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(6)

		c.Set("4", "chetyre")
		c.Set("3", "tri")
		c.Set("2", "dva")
		c.Set("1", "odin")
		c.Set("2", "dva")
		c.Set("3", "tri")
		c.Set("5", "piat")
		c.Set("6", "shest")
		c.Set("1", "odin")
		c.Set("1", "odin")
		c.Set("7", "sem")

		val, ok := c.Get("4")
		require.False(t, ok)
		require.Nil(t, val)

		c.Clear()
		c.Set("1", "odin")
		c.Set("2", "dva")
		c.Set("3", "tri")

		val, ok = c.Get("3")
		require.True(t, ok)
		require.Equal(t, "tri", val)

		val, ok = c.Get("2")
		require.True(t, ok)
		require.Equal(t, "dva", val)

		val, ok = c.Get("1")
		require.True(t, ok)
		require.Equal(t, "odin", val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip() // Remove if task with asterisk completed

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			c.Get(strconv.Itoa(rand.Intn(1000000)))
		}
	}()

	wg.Wait()
}
