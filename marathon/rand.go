package main

import (
	"math/rand"
)

var randTable []int
var randMax int

type rnd struct {
	core *rand.Rand
}

func buildRandTable() {
	for i := 0; i < poolNum; i++ {
		num := i + poolBase
		for j := 0; j < num; j++ {
			randTable = append(randTable, i)
		}
	}
	randMax = len(randTable)
}

func newRnd(seed int64) *rnd {
	s := rand.NewSource(seed)
	r := rand.New(s)
	return &rnd{
		core: r,
	}
}

func (r *rnd) Get() int {
	x := r.core.Intn(randMax)
	return randTable[x]
}
