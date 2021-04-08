package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/bluele/gcache"
)

func checkGcache() {

	var mux sync.RWMutex
	mux.Lock()

	load := func(k interface{}) (v interface{}, err error) {

		fmt.Println(`new load`, k)

		v = k.(int) * 100
		return
	}
	cache := gcache.New(20000).LRU().LoaderFunc(load).Build()

	for i := 0; i < 100; i++ {
		j := i
		go func() {
			fmt.Println(`start`, j)
			mux.RLock()
			v, _ := cache.Get(13)
			fmt.Println(`get`, v)
		}()
	}

	time.Sleep(time.Second)

	fmt.Println(`unlock`)
	mux.Unlock()
}
