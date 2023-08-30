package coral

import (
	"fmt"
	"sync/atomic"
)

// IStats interface of stats
type IStats interface {
	IncHit()
	IncMiss()
	IncWait()
	IncEvict(uint64)
}

// Stats ...
type Stats struct {
	Hit   uint64
	Miss  uint64
	Wait  uint64
	Evict uint64
}

// IncHit increment hit count
func (st *Stats) IncHit() {
	atomic.AddUint64(&st.Hit, 1)
}

// IncMiss increment miss count
func (st *Stats) IncMiss() {
	atomic.AddUint64(&st.Miss, 1)
}

// IncWait increment wait count
func (st *Stats) IncWait() {
	atomic.AddUint64(&st.Wait, 1)
}

// IncEvict increment evict count
func (st *Stats) IncEvict(i uint64) {
	atomic.AddUint64(&st.Evict, i)
}

func (st *Stats) String() string {

	if st == nil {
		return `stats turned off`
	}

	rate := float64(0)
	total := st.Hit + st.Miss + st.Wait
	if total > 0 {
		rate = float64(st.Hit) / float64(total)
	}

	return fmt.Sprintf(
		"cache hit: %d, miss: %d, wait: %d, evict: %d, hit rate: %.2f%%\n",
		st.Hit,
		st.Miss,
		st.Wait,
		st.Evict,
		rate*100,
	)
}
