package cache

import (
	"sync"
)

// BranchToNSCache 使用 sync.Map 实现一个并发安全的内存缓存
type BranchToNSCache struct {
	cache sync.Map
}

// 缓存实例的全局变量
var (
	once            sync.Once
	branchToNSCache *BranchToNSCache
)

func GetInstance() *BranchToNSCache {
	once.Do(func() {
		branchToNSCache = &BranchToNSCache{}
	})
	return branchToNSCache
}

func (c *BranchToNSCache) Get(branchName string) (string, bool) {
	if val, ok := c.cache.Load(branchName); ok {
		return val.(string), true
	}
	return "", false
}

func (c *BranchToNSCache) Set(branchName, nsName string) {
	c.cache.Store(branchName, nsName)
}

// CacheItem 用于表示单个缓存键值对
type CacheItem struct {
	BranchName string `json:"branch_name"`
	NsName     string `json:"ns_name"`
}

func (c *BranchToNSCache) Items() []CacheItem {
	var items []CacheItem

	c.cache.Range(func(key, value any) bool {
		branchName, ok := key.(string)
		if !ok {
			// 如果键不是字符串，可以忽略或记录错误
			return true
		}
		nsName, ok := value.(string)
		if !ok {
			// 如果值不是字符串，可以忽略或记录错误
			return true
		}

		items = append(items, CacheItem{
			BranchName: branchName,
			NsName:     nsName,
		})
		return true // 返回 true 继续遍历
	})

	return items
}
