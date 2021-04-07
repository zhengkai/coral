package coral

// Cache ...
type Cache interface {
	Set(k, v interface{}) (err error)
	Get(k interface{}) (v interface{}, err error)
	Clean()
}

// LoadFunc if Get miss, load data by this
type LoadFunc func(k interface{}) (v interface{}, err error)
