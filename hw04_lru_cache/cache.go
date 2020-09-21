package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type Key string

type Cache interface {
	Set(key string, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key string) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                                 // Очистить кэш
}

type lruCache struct {
	Capacity  int
	Queue     List
	Items     map[string]*cacheItem
	cacheLock sync.Mutex
}

// Set.
func (c *lruCache) Set(key string, value interface{}) bool {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	if len(c.Items) == 0 {
		c.Items = map[string]*cacheItem{}
	}
	elem, isExist := c.Items[key]
	if isExist {
		c.Items[key].Value = value
		c.Queue.MoveToFront(elem.ListItem)
		return isExist
	}
	newCacheItem := &cacheItem{
		Value:    value,
		ListItem: c.Queue.PushFront(value),
	}
	newCacheItem.ListItem.CacheKey = key
	c.Items[key] = newCacheItem
	if c.Queue.Len() > c.Capacity {
		delete(c.Items, c.Queue.Back().CacheKey)
		c.Queue.Remove(c.Queue.Back())
	}
	return isExist
}

// Get.
func (c *lruCache) Get(key string) (interface{}, bool) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	elem, isExist := c.Items[key]
	if isExist {
		c.Queue.MoveToFront(elem.ListItem)
		return c.Items[key].Value, isExist
	}
	return nil, isExist
}

// Clear.
func (c *lruCache) Clear() {
	c.cacheLock.Lock()
	c.Capacity = 0
	c.Queue = *(new(List))
	c.Items = map[string]*cacheItem{}
	c.cacheLock.Unlock()
}

type cacheItem struct {
	Value    interface{}
	ListItem *listItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		Queue:    &list{},
		Capacity: capacity,
	}
}
