package main

import (
	"fmt"
	"time"
)

func reporter() {

	var prevSum int

	var sec = 2

	for {
		ready = true
		time.Sleep(time.Second * time.Duration(sec))
		ready = false
		time.Sleep(time.Second / 10)
		st := c.Stats()
		st.Report()

		total := st.Hit + st.Miss + st.Wait

		missSum := 0
		for _, v := range serialTable {
			missSum += v
		}
		sum := 0
		for _, v := range getTable {
			sum += v
		}

		fmt.Println(total, sum, missSum, uint64(sum) == st.Miss, (sum-prevSum)/sec)
		prevSum = sum
	}
}
