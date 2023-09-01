package coral

import "time"

// TimeExpired 如果希望结果只使用一次而不缓存，在 Load 函数里使用该值而别用 nil。nil 表示永不过期，用过期值表示仅仅是为了阻止并发问题。如果 Load 没结束时有其他协程调 Get，他们会返回相同的结果，且结果不会被缓存。
var TimeExpired = time.Unix(0, 0)

// Load 如果缓存里没有，使用该函数读取上游数据
type Load[K comparable, V any] func(K) (V, *time.Time, error)
