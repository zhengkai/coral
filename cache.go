package coral

// Cache ...
type Cache interface {
	Set(k, v interface{}) (err error)
	Get(k interface{}) (v interface{}, err error)
}

type loadFunc func(k interface{}) (v interface{}, err error)
