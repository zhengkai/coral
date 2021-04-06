package coral

import "sync"

type simple struct {
	storeMux sync.RWMutex
	store    map[interface{}]*entry

	loadFn  loadFunc
	loadMux sync.Mutex
	load    map[interface{}]*entry
}

// Get
func (s *simple) Get(key interface{}) (v interface{}, err error) {

	s.storeMux.RLock()
	e, ok := s.store[key]
	s.storeMux.RUnlock()
	if ok {
		return e.value, e.err
	}

	s.loadMux.Lock()
	e, ok = s.load[key]
	if !ok {
		e = &entry{}
		e.mux.Lock()
		s.load[key] = e
	}
	s.loadMux.Unlock()

	if ok {
		e.mux.RLock()
		v = e.value
		err = e.err
		e.mux.RUnlock()
		return
	}

	e.value, e.err = s.loadFn(key)
	v = e.value
	err = e.err
	e.mux.Unlock()

	s.storeMux.Lock()
	s.store[key] = e
	s.storeMux.Unlock()

	return
}
