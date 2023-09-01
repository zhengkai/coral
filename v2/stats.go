package coral

import (
	"fmt"
	"sync/atomic"
)

// IStats 如果想要自定义统计，可以实现该接口并在 Cache 里 SetStats。Cache 默认无统计
type IStats interface {
	IncHit()
	IncMiss()
	IncWait()
	IncEvict(uint64)
}

// Stats IStats 的简单实现
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

// IncWait 如果有并发，其他 goroutine 等待次数
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
