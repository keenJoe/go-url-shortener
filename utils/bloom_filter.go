package utils

import (
	"hash/fnv"
	"sync"
)

// 简单的布隆过滤器实现
type BloomFilter struct {
	bitset []bool
	size   uint
	mutex  sync.RWMutex
}

// NewBloomFilter 创建新的布隆过滤器
func NewBloomFilter(size uint) *BloomFilter {
	return &BloomFilter{
		bitset: make([]bool, size),
		size:   size,
	}
}

// Add 添加元素到布隆过滤器
func (bf *BloomFilter) Add(item string) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	h1 := bf.hash1(item)
	h2 := bf.hash2(item)
	h3 := bf.hash3(item)

	bf.bitset[h1%bf.size] = true
	bf.bitset[h2%bf.size] = true
	bf.bitset[h3%bf.size] = true
}

// Contains 检查元素是否可能存在
func (bf *BloomFilter) Contains(item string) bool {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	h1 := bf.hash1(item)
	h2 := bf.hash2(item)
	h3 := bf.hash3(item)

	return bf.bitset[h1%bf.size] && bf.bitset[h2%bf.size] && bf.bitset[h3%bf.size]
}

// 哈希函数
func (bf *BloomFilter) hash1(item string) uint {
	h := fnv.New32a()
	h.Write([]byte(item))
	return uint(h.Sum32())
}

func (bf *BloomFilter) hash2(item string) uint {
	h := fnv.New32a()
	h.Write([]byte(item + item))
	return uint(h.Sum32())
}

func (bf *BloomFilter) hash3(item string) uint {
	h := fnv.New64a()
	h.Write([]byte(item))
	return uint(h.Sum64())
}

// 全局布隆过滤器实例
var (
	ShortCodeFilter   *BloomFilter
	OriginalURLFilter *BloomFilter
)

// InitBloomFilters 初始化布隆过滤器
func InitBloomFilters() {
	// 预估1亿个元素，假阳性率0.01
	ShortCodeFilter = NewBloomFilter(1000000000)
	OriginalURLFilter = NewBloomFilter(1000000000)
}
