package coral

import (
	"fmt"
	"sync/atomic"
)

// Stats ...
type Stats struct {
	Hit   uint64
	Miss  uint64
	Wait  uint64
	Evict uint64
}

// IncHit increment hit count
func (st *Stats) IncHit() {
	if st == nil {
		return
	}
	atomic.AddUint64(&st.Hit, 1)
}

// IncMiss increment miss count
func (st *Stats) IncMiss() {
	if st == nil {
		return
	}
	atomic.AddUint64(&st.Miss, 1)
}

// IncWait increment wait count
func (st *Stats) IncWait() {
	if st == nil {
		return
	}
	atomic.AddUint64(&st.Wait, 1)
}

// IncEvict increment evict count
func (st *Stats) IncEvict(i uint64) {
	if st == nil {
		return
	}
	atomic.AddUint64(&st.Evict, i)
}

// Report ...
func (st *Stats) Report() {

	if st == nil {
		fmt.Println(`stats turned off`)
		return
	}

	rate := float64(0)
	total := st.Hit + st.Miss + st.Wait
	if total > 0 {
		rate = float64(st.Hit) / float64(total)
	}

	fmt.Printf(
		"hit: %d, miss: %d, wait: %d, evict: %d, hit rate: %f\n",
		st.Hit,
		st.Miss,
		st.Wait,
		st.Evict,
		rate,
	)
}
