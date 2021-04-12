package coral

import (
	"sync"
	"time"
)

type entry struct {
	value  interface{}
	err    error
	mux    sync.RWMutex
	expire *time.Time
}
