package main

import (
	"fmt"
	"time"

	"github.com/zhengkai/coral/v2"
)

func simpleTest() {
	c := coral.NewSimple(func(k uint32) (int32, *time.Time, error) {
		return -int32(k), nil, nil
	})

	c.Set(1, 2, nil)
	fmt.Println(c.Get(1))
}

func simpleTimeout() {

	/*
		loadFn := func(k interface{}) (v interface{}, expire *time.Time, err error) {
			fmt.Println(`new load`, k)
			time.Sleep(time.Second / 4)
			v = k.(int) * 100
			return
		}

		c1 := coral.NewSimple(coral.LoadDeadline(loadFn, time.Second/10))
		v, err := c1.Get(17)
		fmt.Println(`c1`, v, err)

		c2 := coral.NewSimple(coral.LoadDeadline(loadFn, time.Second/3))
		v, err = c2.Get(17)
		fmt.Println(`c2`, v, err)
	*/
}
