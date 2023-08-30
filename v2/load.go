package coral

import "time"

// Load ...
type Load[K comparable, V any] func(K) (V, *time.Time, error)
