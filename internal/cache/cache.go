package cache

import (
	"container/list"
	"sync"
	"time"
)

type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

type Cache[K comparable, V any] struct {
	mu       sync.RWMutex
	items    map[K]*list.Element
	order    *list.List
	maxSize  int
	ttl      time.Duration
}

func New[K comparable, V any](maxSize int, ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		items:   make(map[K]*list.Element),
		order:   list.New(),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	e := elem.Value.(*entry[K, V])
	if time.Now().After(e.expiresAt) {
		c.removeElement(elem)
		var zero V
		return zero, false
	}

	c.order.MoveToFront(elem)
	return e.value, true
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		e := elem.Value.(*entry[K, V])
		e.value = value
		e.expiresAt = time.Now().Add(c.ttl)
		return
	}

	e := &entry[K, V]{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	elem := c.order.PushFront(e)
	c.items[key] = elem

	if c.order.Len() > c.maxSize {
		c.removeOldest()
	}
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
	}
}

func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]*list.Element)
	c.order.Init()
}

func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.order.Len()
}

func (c *Cache[K, V]) removeElement(elem *list.Element) {
	c.order.Remove(elem)
	e := elem.Value.(*entry[K, V])
	delete(c.items, e.key)
}

func (c *Cache[K, V]) removeOldest() {
	elem := c.order.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}
