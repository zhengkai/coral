package main

import (
	"time"

	"github.com/bluele/gcache"
	"github.com/zhengkai/coral"
)

var serialPool []chan int

func provider() {

	for i := 0; i < poolNum; i++ {

		ch := make(chan int)
		serialPool = append(serialPool, ch)
		go func() {
			j := 0
			for {
				j++
				ch <- j
			}
		}()
	}

	loadFn := func(k interface{}) (v interface{}, expire *time.Time, err error) {

		id := k.(int)

		ch := serialPool[id]

		serial := <-ch

		v = &val{
			id:     id,
			serial: serial,
		}
		return
	}

	c = coral.BuildLRU(loadFn, cacheCapacity, cacheEvictThreshold)

	loadFnG := func(k interface{}) (v interface{}, err error) {

		id := k.(int)

		ch := serialPool[id]

		serial := <-ch

		v = &val{
			id:     id,
			serial: serial,
		}
		return
	}

	c2 = gcache.New(cacheCapacity).LRU().LoaderFunc(loadFnG).Build()
}
