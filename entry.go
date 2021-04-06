package coral

import "sync"

type entry struct {
	value interface{}
	err   error
	mux   sync.RWMutex
}
