package coral

import "time"

// Cache ...
type Cache interface {
	Set(k, v interface{}, expire *time.Time) (err error)
	Get(k interface{}) (v interface{}, err error)
	Clean()
	Delete(k interface{})

	StatsOff()
	Stats() (st *Stats)
}

// LoadFunc if Get miss, load data by this
type LoadFunc func(k interface{}) (v interface{}, expire *time.Time, err error)
