package main

import (
	"fmt"
	"time"

	"github.com/zhengkai/coral"
)

func simpleTest() {
	x := coral.BuildSimple(func(k interface{}) (v interface{}, expire *time.Time, err error) {
		return
	})

	x.Set(1, 2)
	i, err := x.Get(1)
	fmt.Println(i, err)
}

func simpleTimeout() {

	loadFn := func(k interface{}) (v interface{}, expire *time.Time, err error) {
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
