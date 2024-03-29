package coral

import (
	"time"
)

type simple[K comparable, V any] struct {
	cache[K, V]
}

// NewSimple 无淘汰全部缓存，不太可能用上，只是功能最少，作为其他缓存的参考
func NewSimple[K comparable, V any](load Load[K, V]) Cache[K, V] {
	o := simple[K, V]{
		cache[K, V]{
			load: load,
		},
	}
	o.Reset()
	return &o
}

func (c *simple[K, V]) Set(k K, v V, expire *time.Time) {
	en := &entry[V]{
		v:      v,
		expire: expire,
		done:   true,
	}
	c.mux.Lock()
	c.pool[k] = en
	c.mux.Unlock()
}

func (c *simple[K, V]) Get(k K) (v V, err error) {
	c.mux.Lock()
	ey, ok := c.pool[k].check()
	if ok {
		c.mux.Unlock()
		c.statsHit()
		return ey.v, nil
	}
	if ey == nil {
		ey = newEntry[V]()
		c.pool[k] = ey
		c.mux.Unlock()
		go func() {
			ey.set(c.load(k))
		}()
		c.statsMiss()
	} else {
		c.mux.Unlock()
		c.statsWait()
	}
	ey.wait()
	return ey.v, ey.err
}

func (c *simple[K, V]) Reset() {
	c.mux.Lock()
	c.pool = make(map[K]*entry[V])
	c.mux.Unlock()
}

func (c *simple[K, V]) Delete(k K) {
	c.mux.Lock()
	delete(c.pool, k)
	c.mux.Unlock()
}
