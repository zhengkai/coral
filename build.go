package coral

import "container/list"

// BuildSimple ...
func BuildSimple(loadFn LoadFunc) (c Cache) {

	c = &simple{
		store:  make(map[interface{}]*entry),
		load:   make(map[interface{}]*entry),
		loadFn: loadFn,
		stats:  &Stats{},
	}

	return
}

// BuildLRU ...
func BuildLRU(loadFn LoadFunc, capacity int, evictThreshold int) (c Cache) {

	if capacity < 1 {
		capacity = 1000
	}
	if evictThreshold < capacity {
		evictThreshold = capacity + capacity/20
	}

	c = &lru{
		store:  make(map[interface{}]*entryRU),
		load:   make(map[interface{}]*entryRU),
		loadFn: loadFn,
		list:   list.New(),
		stats:  &Stats{},

		capacity:       capacity,
		evictThreshold: evictThreshold,
	}

	return
}
