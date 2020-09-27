package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key string) (interface{}, bool)  // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	capacity  int
	queue     List
	items     map[string]*cacheItem
	cacheLock sync.Mutex
}

type cacheItem struct {
	Value    interface{}
	ListItem *listItem
}

// Set.
func (c *lruCache) Set(key Key, value interface{}) bool {
	strKey := string(key)
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	elem, isExist := c.items[strKey]
	if isExist {
		c.items[strKey].Value = value
		c.queue.MoveToFront(elem.ListItem)
		return isExist
	}
	newCacheItem := &cacheItem{
		Value:    value,
		ListItem: c.queue.PushFront(value),
	}
	newCacheItem.ListItem.CacheKey = strKey
	c.items[strKey] = newCacheItem
	if c.queue.Len() > c.capacity {
		delete(c.items, c.queue.Back().CacheKey)
		c.queue.Remove(c.queue.Back())
	}
	return isExist
}

// Get.
func (c *lruCache) Get(key string) (interface{}, bool) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	elem, isExist := c.items[key]
	if isExist {
		c.queue.MoveToFront(elem.ListItem)
		return c.items[key].Value, isExist
	}
	return nil, isExist
}

// Clear.
func (c *lruCache) Clear() {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	c.queue = &list{}
	c.items = map[string]*cacheItem{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		queue:    &list{},
		capacity: capacity,
		items:    map[string]*cacheItem{},
	}
}
