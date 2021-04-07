package coral

// Algorithms list https://en.wikipedia.org/wiki/Cache_replacement_policies
const (
	LRU Algorithms = iota
	LFU
	MRU
	FIFO
	RR
)

// Algorithms cache algorithms
type Algorithms uint8
