package coral

import (
	"errors"
	"time"
)

// ErrLoadFuncTimeout timeout error when running LoadFunc
var ErrLoadFuncTimeout = errors.New(`loadFn timeout`)

// LoadDeadline add timeout with LoadFunc
func LoadDeadline(loadFn LoadFunc, timeout time.Duration) LoadFunc {

	return func(k interface{}) (v interface{}, err error) {

		ch := make(chan *entry)

		go func() {
			e := &entry{}
			e.value, e.err = loadFn(k)
			ch <- e
		}()

		select {
		case e := <-ch:

			v = e.value
			err = e.err

		case <-time.After(timeout):

			err = ErrLoadFuncTimeout
		}

		return
	}
}
