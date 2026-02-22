package cache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type ByteLRU struct {
	mu         sync.Mutex
	maxEntries int
	defaultTTL time.Duration
	ll         *list.List
	items      map[string]*list.Element
}

type cacheEntry struct {
	key       string
	value     []byte
	expiresAt time.Time
}

func NewByteLRU(maxEntries int, defaultTTL time.Duration) *ByteLRU {
	if maxEntries < 1 {
		maxEntries = 1
	}
	if defaultTTL < 0 {
		defaultTTL = 0
	}
	return &ByteLRU{
		maxEntries: maxEntries,
		defaultTTL: defaultTTL,
		ll:         list.New(),
		items:      make(map[string]*list.Element, maxEntries),
	}
}

func (c *ByteLRU) Get(_ context.Context, key string) ([]byte, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return nil, false, nil
	}

	entry := elem.Value.(*cacheEntry)
	if c.isExpired(entry) {
		c.removeElement(elem)
		return nil, false, nil
	}

	c.ll.MoveToFront(elem)
	copied := make([]byte, len(entry.value))
	copy(copied, entry.value)
	return copied, true, nil
}

func (c *ByteLRU) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	if elem, ok := c.items[key]; ok {
		entry := elem.Value.(*cacheEntry)
		entry.value = append(entry.value[:0], value...)
		entry.expiresAt = c.expirationTime(ttl)
		c.ll.MoveToFront(elem)
		return nil
	}

	copied := make([]byte, len(value))
	copy(copied, value)

	entry := &cacheEntry{
		key:       key,
		value:     copied,
		expiresAt: c.expirationTime(ttl),
	}

	elem := c.ll.PushFront(entry)
	c.items[key] = elem

	for c.ll.Len() > c.maxEntries {
		last := c.ll.Back()
		if last == nil {
			break
		}
		c.removeElement(last)
	}
	return nil
}

func (c *ByteLRU) Close() error {
	return nil
}

func (c *ByteLRU) isExpired(entry *cacheEntry) bool {
	if entry.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(entry.expiresAt)
}

func (c *ByteLRU) expirationTime(ttl time.Duration) time.Time {
	if ttl <= 0 {
		return time.Time{}
	}
	return time.Now().Add(ttl)
}

func (c *ByteLRU) removeElement(elem *list.Element) {
	c.ll.Remove(elem)
	entry := elem.Value.(*cacheEntry)
	delete(c.items, entry.key)
}
