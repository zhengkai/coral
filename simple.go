package coral

import "sync"

// no expire, no eviction

type simple struct {
	storeMux sync.RWMutex
	store    map[interface{}]*entry

	loadFn  LoadFunc
	loadMux sync.Mutex
	load    map[interface{}]*entry
}

// Get
func (s *simple) Get(k interface{}) (v interface{}, err error) {

	ok, v, err := s.storeGet(k)
	if ok {
		return
	}

	ok, e, v, err := s.loadGet(k)
	if ok {
		return
	}

	v, err = s.loadExec(k, e)
	return
}

func (s *simple) storeGet(k interface{}) (ok bool, v interface{}, err error) {
	s.storeMux.RLock()
	e, ok := s.store[k]
	s.storeMux.RUnlock()
	if ok {
		v = e.value
		err = e.err
	}
	return
}

func (s *simple) loadExec(k interface{}, e *entry) (v interface{}, err error) {

	e.value, e.err = s.loadFn(k)
	v = e.value
	err = e.err
	e.mux.Unlock()

	s.storeMux.Lock()
	s.store[k] = e
	s.storeMux.Unlock()

	s.loadMux.Lock()
	delete(s.load, k)
	s.loadMux.Unlock()

	return
}

func (s *simple) loadGet(k interface{}) (ok bool, e *entry, v interface{}, err error) {
	s.loadMux.Lock()
	e, ok = s.load[k]
	if !ok {
		ok, v, err = s.storeGet(k)
		if ok {
			s.loadMux.Unlock()
			return
		}
		e = &entry{}
		e.mux.Lock()
		s.load[k] = e
	}
	s.loadMux.Unlock()

	if ok {
		e.mux.RLock()
		v = e.value
		err = e.err
		e.mux.RUnlock()
	}

	return
}

func (s *simple) Set(k interface{}, v interface{}) (err error) {

	s.storeMux.Lock()
	s.store[k] = &entry{
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

func (s *simple) Delete(k interface{}) {

	s.storeMux.Lock()
	delete(s.store, k)
	s.storeMux.Unlock()
}
