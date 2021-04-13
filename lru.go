package coral

import (
	"container/list"
	"sync"
	"time"
)

// no eviction

type lru struct {
	storeMux sync.Mutex
	store    map[interface{}]*entryRU

	loadFn  LoadFunc
	loadMux sync.Mutex
	load    map[interface{}]*entryRU

	list *list.List

	evictThreshold int
	capacity       int
}

// Get
func (s *lru) Get(k interface{}) (v interface{}, err error) {

	e := s.storeGet(k)
	if e != nil {
		v = e.value
		err = e.err
		return
	}

	ok, e := s.loadGet(k)
	if ok {
		v = e.value
		err = e.err
		return
	}

	v, err = s.loadExec(k, e)
	return
}

func (s *lru) storeGet(k interface{}) (en *entryRU) {

	s.storeMux.Lock()
	en, ok := s.store[k]
	defer s.storeMux.Unlock()
	if ok {
		if en.expire != nil && time.Now().After(*en.expire) {
			en = nil
			return
		}
		s.list.MoveToFront(en.cur)
	}
	return
}

func (s *lru) loadExec(k interface{}, e *entryRU) (v interface{}, err error) {

	e.value, e.expire, e.err = s.loadFn(k)
	v = e.value
	err = e.err
	e.mux.Unlock()

	if err == nil || e.expire != nil {
		s.storeMux.Lock()
		prev, found := s.store[k]
		if found {
			s.list.Remove(prev.cur)
		}
		cur := s.list.PushFront(k)
		e.cur = cur
		s.store[k] = e

		if s.list.Len() > s.evictThreshold {
			s.slim()
		}
		s.storeMux.Unlock()
	}

	s.loadMux.Lock()
	delete(s.load, k)
	s.loadMux.Unlock()

	return
}

func (s *lru) loadGet(k interface{}) (ok bool, e *entryRU) {
	s.loadMux.Lock()
	e, ok = s.load[k]
	if !ok {
		e = s.storeGet(k)
		if e != nil {
			s.loadMux.Unlock()
			ok = true
			return
		}
		e = &entryRU{}
		e.mux.Lock()
		s.load[k] = e
	}
	s.loadMux.Unlock()

	if ok {
		e.mux.RLock()
	}

	return
}

func (s *lru) Set(k interface{}, v interface{}, expire *time.Time) (err error) {

	s.storeMux.Lock()
	e, ok := s.store[k]
	if ok {
		e.value = v
		e.expire = expire
		s.list.MoveToFront(e.cur)
	} else {
		e = &entryRU{
			value:  v,
			expire: expire,
		}
		cur := s.list.PushFront(k)
		e.cur = cur
		s.store[k] = e
		if s.list.Len() > s.evictThreshold {
			s.slim()
		}
	}
	s.storeMux.Unlock()

	return
}

func (s *lru) Clean() {

	m := make(map[interface{}]*entryRU)
	s.storeMux.Lock()
	s.store = m
	s.list = list.New()
	s.storeMux.Unlock()
}

func (s *lru) slim() {

	now := time.Now()
	for k, v := range s.store {
		if v.expire != nil && v.expire.Before(now) {
			delete(s.store, k)
			s.list.Remove(v.cur)
		}
	}
	for i := s.list.Len(); i >= s.capacity; i-- {
		k := s.list.Remove(s.list.Back())
		delete(s.store, k)
	}
}

func (s *lru) Delete(k interface{}) {

	s.storeMux.Lock()
	en, ok := s.store[k]
	if ok {
		s.list.Remove(en.cur)
		delete(s.store, k)
	}
	s.storeMux.Unlock()
}
