package coral

import (
	"container/list"
	"sync"
	"time"
)

type entry[V any] struct {
	v      V
	err    error
	expire *time.Time

	mux  sync.RWMutex
	done bool
	once sync.Once

	cur *list.Element
}

func newEntry[V any]() *entry[V] {
	ey := &entry[V]{}
	ey.mux.Lock()
	return ey
}

func (ey *entry[V]) wait() (V, error) {
	ey.mux.RLock()
	return ey.v, ey.err
}

func (ey *entry[V]) set(v V, expire *time.Time, err error) {
	ey.v = v
	ey.expire = expire
	ey.err = err
	ey.once.Do(func() {
		ey.done = true
		ey.mux.Unlock()
	})
}

func (ey *entry[V]) check() (*entry[V], bool) {
	if ey == nil || !ey.done {
		return ey, false
	}
	if ey.err != nil {
		return nil, false
	}
	if ey.expire != nil && ey.expire.Before(time.Now()) {
		return nil, false
	}
	return ey, true
}
