package test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zhengkai/coral"
)

func TestSimpleConcurrency(t *testing.T) {
	// t.Fatal("not implemented")

	var mux sync.RWMutex
	mux.Lock()

	var count uint32

	var wg sync.WaitGroup
	wg.Add(100)

	c := coral.BuildSimple(func(k interface{}) (v interface{}, err error) {
		atomic.AddUint32(&count, 1)
		v = k.(int) * 100
		return
	})

	for i := 0; i < 100; i++ {
		go func() {
			mux.RLock()
			v, err := c.Get(13)
			wg.Done()
			if v != 1300 {
				t.Error(`simple load fail`)
			}
			if err != nil {
				t.Error(`simple load err`)
			}
		}()
	}

	time.Sleep(time.Second / 100)

	mux.Unlock()
	wg.Wait()

	if count != 1 {
		t.Error(`simple concurrency failed`)
	}
}

func TestSimpleConcurrencyDelay(t *testing.T) {
	// t.Fatal("not implemented")

	var mux sync.RWMutex
	mux.Lock()

	var count uint32

	var wg sync.WaitGroup
	wg.Add(100)

	c := coral.BuildSimple(func(k interface{}) (v interface{}, err error) {
		atomic.AddUint32(&count, 1)
		time.Sleep(time.Second / 100)
		v = k.(int) * 100
		return
	})

	for i := 0; i < 100; i++ {
		go func() {
			mux.RLock()
			v, err := c.Get(13)
			wg.Done()
			if v != 1300 {
				t.Error(`simple load fail`)
			}
			if err != nil {
				t.Error(`simple load err`)
			}
		}()
	}

	time.Sleep(time.Second / 50)

	mux.Unlock()
	wg.Wait()

	if count != 1 {
		t.Error(`simple concurrency failed`)
	}
}

func TestSimpleMisc(t *testing.T) {

	var count int

	c := coral.BuildSimple(func(k interface{}) (v interface{}, err error) {
		count++
		time.Sleep(time.Second / 100)
		v = k.(int) * 100
		return
	})

	c.Get(1)
	c.Get(2)
	c.Delete(1)
	c.Get(1)
	c.Get(2)

	if count != 3 {
		t.Error(`simple delete fail`)
	}

	c.Clean()
	c.Get(1)
	c.Get(2)

	if count != 5 {
		t.Error(`simple clean fail`)
	}

	c.Set(3, true)
	c.Get(3)

	if count != 5 {
		t.Error(`simple set fail`)
	}
}

func TestSimpleDeadline(t *testing.T) {

	loadFn := func(k interface{}) (v interface{}, err error) {
		time.Sleep(time.Second / 100)
		v = k.(int) * 100
		return
	}

	c := coral.BuildSimple(coral.LoadDeadline(loadFn, time.Second/150))
	_, err := c.Get(1)
	if err != coral.ErrLoadFuncTimeout {
		t.Error(`LoadDeadline fail`)
	}

	c = coral.BuildSimple(coral.LoadDeadline(loadFn, time.Second/50))
	_, err = c.Get(1)
	if err == coral.ErrLoadFuncTimeout {
		t.Error(`LoadDeadline fail`)
	}
}
