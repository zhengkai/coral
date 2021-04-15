package main

import (
	"time"

	"github.com/bluele/gcache"
	"github.com/zhengkai/coral"
)

var c coral.Cache
var c2 gcache.Cache

type val struct {
	id     int
	serial int
}

func main() {

	buildRandTable()

	provider()

	serialTable = make([]int, poolNum)
	getTable = make([]int, poolNum)

	ts := time.Now().UnixNano()
	for i := 0; i < customerNum; i++ {
		go customer(i, ts)
	}

	reporter()
}
