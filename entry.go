package coral

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	value  interface{}
	err    error
	mux    sync.RWMutex
	expire *time.Time
}

type entryRU struct {
	value  interface{}
	err    error
	mux    sync.RWMutex
	expire *time.Time

	cur *list.Element
}
