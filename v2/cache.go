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

	SetStats(*Stats)
	GetStats() *Stats
}

type cache[K comparable, V any] struct {
	load  Load[K, V]
	pool  map[K]*entry[V]
	list  *list.List
	mux   sync.Mutex
	stats *Stats
}

func (c *cache[K, V]) SetStats(s *Stats) {
	c.stats = s
}

func (c *cache[K, V]) GetStats() *Stats {
	return c.stats
}
