package coral

import "sync"

type simple struct {
	storeMux sync.RWMutex
	store    map[interface{}]*entry

	loadFn  LoadFunc
	loadMux sync.Mutex
	load    map[interface{}]*entry
}

// Get
func (s *simple) Get(key interface{}) (v interface{}, err error) {

	ok, v, err := s.get(key)
	if ok {
		return
	}

	s.loadMux.Lock()
	e, ok := s.load[key]
	if !ok {
		ok, v, err = s.get(key)
		if ok {
			return
		}
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

	s.loadMux.Lock()
	delete(s.load, key)
	s.loadMux.Unlock()

	return
}

func (s *simple) get(key interface{}) (ok bool, v interface{}, err error) {
	s.storeMux.RLock()
	e, ok := s.store[key]
	s.storeMux.RUnlock()
	if ok {
		v = e.value
		err = e.err
	}
	return
}

func (s *simple) Set(key interface{}, v interface{}) (err error) {

	s.storeMux.Lock()
	s.store[key] = &entry{
		value: v,
	}
	s.storeMux.Unlock()

	return
}

func (s *simple) Clean() {

	s.storeMux.Lock()
	s.store = make(map[interface{}]*entry)
	s.storeMux.Unlock()
}
