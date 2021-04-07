package main

import (
	"fmt"

	"github.com/zhengkai/coral"
)

func main() {
	x := coral.BuildSimple(func(k interface{}) (v interface{}, err error) {
		return
	})

	x.Set(1, 2)
	i, err := x.Get(1)
	fmt.Println(i, err)
}
