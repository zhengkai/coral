package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/zhengkai/coral"
)

func simpleTest() {
	x := coral.BuildSimple(func(k interface{}) (v interface{}, err error) {
		return
	})

	x.Set(1, 2)
	i, err := x.Get(1)
	fmt.Println(i, err)
}

func simpleConcurrency() {

	var mux sync.RWMutex
	mux.Lock()

	c := coral.BuildSimple(func(k interface{}) (v interface{}, err error) {

		fmt.Println(`new load`, k)

		v = k.(int) * 100
		return
	})

	for i := 0; i < 100; i++ {
		j := i
		go func() {
			fmt.Println(`start`, j)
			mux.RLock()
			v, _ := c.Get(13)
			fmt.Println(`get`, v)
		}()
	}

	time.Sleep(time.Second)

	fmt.Println(`unlock`)
	mux.Unlock()
}

func simpleTimeout() {

	loadFn := func(k interface{}) (v interface{}, err error) {
		fmt.Println(`new load`, k)
		time.Sleep(time.Second / 4)
		v = k.(int) * 100
		return
	}

	c1 := coral.BuildSimple(coral.LoadDeadline(loadFn, time.Second/10))
	v, err := c1.Get(17)
	fmt.Println(`c1`, v, err)

	c2 := coral.BuildSimple(coral.LoadDeadline(loadFn, time.Second/3))
	v, err = c2.Get(17)
	fmt.Println(`c2`, v, err)
}
