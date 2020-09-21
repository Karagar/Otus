package hw04_lru_cache //nolint:golint,stylecheck

import (
	"errors"
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
		newCache := NewCache(6)
		err := error(nil)
		slice := newCache.(*lruCache).Queue.ToSlice()
		targetSlice := []interface{}{}

		for i, elem := range targetSlice {
			if elem != slice[i] {
				err = errors.New("Expected and recieved list are different")
			}
		}

		if err != nil {
			t.Errorf("\n\t%s", err)
		} else {
			newCache.Set("1", "odin")
			newCache.Set("1", "odin")
			newCache.Set("1", "odin")
			newCache.Set("1", "odin")

			slice := newCache.(*lruCache).Queue.ToSlice()
			targetSlice = []interface{}{"odin"}
			for i, elem := range targetSlice {
				if elem != slice[i] {
					err = errors.New("Expected and recieved list are different")
				}
			}
		}

		if err != nil {
			t.Errorf("\n\t%s", err)
		} else {
			newCache.Set("2", "dva")
			newCache.Set("3", "tri")
			newCache.Set("4", "chetyre")
			newCache.Set("2", "dva")
			newCache.Set("3", "tri")
			newCache.Set("5", "piat")
			newCache.Set("6", "shest")
			newCache.Set("1", "odin")
			newCache.Set("1", "odin")
			newCache.Set("1", "odin")
			newCache.Set("7", "sem")

			slice := newCache.(*lruCache).Queue.ToSlice()
			targetSlice = []interface{}{"sem", "odin", "shest", "piat", "tri", "dva"}
			for i, elem := range targetSlice {
				if elem != slice[i] {
					err = errors.New("Expected and recieved list are different")
				}
			}
		}

		if err != nil {
			t.Errorf("\n\t%s", err)
		} else {
			newCache.Get("1")
			newCache.Set("7", "sem")
			newCache.Get("1")

			slice := newCache.(*lruCache).Queue.ToSlice()
			targetSlice = []interface{}{"odin", "sem", "shest", "piat", "tri", "dva"}
			for i, elem := range targetSlice {
				if elem != slice[i] {
					err = errors.New("Expected and recieved list are different")
				}
			}
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip() // Remove if task with asterisk completed

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			c.Set(strconv.Itoa(i), i)
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
