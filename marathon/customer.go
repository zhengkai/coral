package main

import (
	"fmt"
	"time"
)

var serialTable []int
var getTable []int

var ready = false

func customer(id int, ts int64) {
	r := newRnd(int64(id) + ts)
	for {

		if !ready {
			time.Sleep(time.Microsecond)
			continue
		}

		k := r.Get()
		x, _ := c.Get(k)
		v := x.(*val)
		if v.id != k {
			fmt.Println(id, v)
			panic(`value wrong`)
		}
		serialTable[k] = v.serial
		getTable[k]++
	}
}
