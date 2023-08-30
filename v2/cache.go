package coral

import (
	"container/list"
	"sync"
	"time"
)

// Cache ...
type Cache[K comparable, V any] interface {
	Set(p K, v V, expire *time.Time)
	Get(p K) (v V, err error)
	Reset()
	Delete(p K)

	SetStats(IStats)
	GetStats() IStats
}

type cache[K comparable, V any] struct {
	load  Load[K, V]
	pool  map[K]*entry[V]
	list  *list.List
	mux   sync.Mutex
	stats IStats
}

func (c *cache[K, V]) SetStats(s IStats) {
	c.stats = s
}

func (c *cache[K, V]) GetStats() IStats {
	return c.stats
}

func (c *cache[K, V]) statsHit() {
	if c.stats != nil {
		go c.stats.IncHit()
	}
}

func (c *cache[K, V]) statsMiss() {
	if c.stats != nil {
		go c.stats.IncMiss()
	}
}

func (c *cache[K, V]) statsWait() {
	if c.stats != nil {
		go c.stats.IncWait()
	}
}

func (c *cache[K, V]) statsEvict(i uint64) {
	if c.stats != nil {
		go c.stats.IncEvict(i)
	}
}
