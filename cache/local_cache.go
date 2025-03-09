package cache

import (
	"sync"
	"time"
)

// 简单的LRU缓存实现
type CacheItem struct {
	Value      string
	Expiration int64
}

type LocalCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

var localCache *LocalCache

// InitLocalCache 初始化本地缓存
func InitLocalCache() {
	localCache = &LocalCache{
		items: make(map[string]CacheItem),
	}

	// 启动清理过期项的goroutine
	go localCache.cleanupLoop()
}

// Set 设置缓存
func (c *LocalCache) Set(key string, value string, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(duration).UnixNano()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Get 获取缓存
func (c *LocalCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return "", false
	}

	// 检查是否过期
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return "", false
	}

	return item.Value, true
}

// Delete 删除缓存
func (c *LocalCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// 定期清理过期项
func (c *LocalCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.cleanup()
	}
}

func (c *LocalCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
		}
	}
}

// SetURLLocal 设置本地缓存
func SetURLLocal(shortCode, originalURL string, duration time.Duration) {
	localCache.Set("url:"+shortCode, originalURL, duration)
}

// GetURLLocal 获取本地缓存
func GetURLLocal(shortCode string) (string, bool) {
	return localCache.Get("url:" + shortCode)
}
