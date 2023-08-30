package coral

import (
	"container/list"
	"time"
)

type lru[K comparable, V any] struct {
	cache[K, V]
	list *list.List

	capacity       int
	evictThreshold int
}

// NewLRU ...
func NewLRU[K comparable, V any](load Load[K, V], capacity, evictBuffer uint32) Cache[K, V] {
	o := lru[K, V]{
		cache[K, V]{
			load: load,
		},
		list.New(),
		int(capacity),
		int(capacity + evictBuffer),
	}
	o.Reset()
	return &o
}

func (c *lru[K, V]) Set(k K, v V, expire *time.Time) {

	c.mux.Lock()
	ey, ok := c.pool[k].check()
	if ok {
		ey.set(v, expire, nil)
		c.list.MoveToFront(ey.cur)
		c.mux.Unlock()
		return
	}

	ne := &entry[V]{
		v:      v,
		expire: expire,
		done:   true,
		cur:    c.list.PushFront(k),
	}
	c.pool[k] = ne
	c.slim()
	c.mux.Unlock()
}

func (c *lru[K, V]) slim() {

	if c.list.Len() <= c.evictThreshold {
		return
	}

	now := time.Now()
	var cnt uint64
	for k, v := range c.pool {
		if v.err != nil || (v.expire != nil && v.expire.Before(now)) {
			c.list.Remove(v.cur)
			delete(c.pool, k)
			cnt++
		}
	}
	for i := c.list.Len(); i > c.capacity; i-- {
		k := c.list.Remove(c.list.Back()).(K)
		delete(c.pool, k)
		cnt++
	}
	if cnt > 0 {
		c.stats.IncEvict(cnt)
	}
}

func (c *lru[K, V]) Get(k K) (v V, err error) {

	c.mux.Lock()
	ey, ok := c.pool[k].check()
	if ok {
		c.list.MoveToFront(ey.cur)
		c.mux.Unlock()
		c.stats.IncHit()
		return ey.v, nil
	}
	if ey == nil {
		ey = newEntry[V]()
		ey.cur = c.list.PushFront(k)
		c.pool[k] = ey
		c.slim()
		c.mux.Unlock()
		go func() {
			ey.set(c.load(k))
		}()
		c.stats.IncMiss()
	} else {
		c.mux.Unlock()
		c.stats.IncWait()
	}
	ey.wait()
	return ey.v, ey.err
}

func (c *lru[K, V]) Reset() {
	c.mux.Lock()
	c.pool = make(map[K]*entry[V])
	c.list = list.New()
	c.mux.Unlock()
}

func (c *lru[K, V]) Delete(k K) {
	c.mux.Lock()
	ey, ok := c.pool[k]
	if ok {
		delete(c.pool, k)
		c.list.Remove(ey.cur)
	}
	c.mux.Unlock()
}
