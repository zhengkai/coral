package test

import (
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zhengkai/coral/v2"
)

func TestSimpleConcurrency(t *testing.T) {

	for j := 0; j < 200; j++ {

		concurrency := 20

		var mux sync.RWMutex
		mux.Lock()

		var count uint32

		var wg sync.WaitGroup
		wg.Add(concurrency)

		c := coral.NewSimple(func(k int) (int, *time.Time, error) {
			atomic.AddUint32(&count, 1)
			v := k * 100
			return v, nil, nil
		})

		if j%2 == 0 {
			c.SetStats(&coral.Stats{})
		}

		for i := 0; i < concurrency; i++ {
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
}

func TestSimpleConcurrencyDelay(t *testing.T) {
	// t.Fatal("not implemented")

	var mux sync.RWMutex
	mux.Lock()

	var wg sync.WaitGroup
	wg.Add(100)

	count := &atomic.Uint32{}

	c := coral.NewSimple(func(k int) (v int, expire *time.Time, err error) {
		count.Add(1)
		time.Sleep(time.Second / 100)
		v = k * 100
		return
	})
	st := &coral.Stats{}
	c.SetStats(st)

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

	if count.Load() != 1 {
		t.Error(`simple concurrency fail`)
	}
	if st.Wait != 99 {
		t.Error(`simple concurrency wait count fail`)
	}
	if !strings.Contains(st.String(), `wait: 99,`) {
		t.Error(`simple concurrency wait count fail`)
	}
}

func TestSimpleMisc(t *testing.T) {

	var count int

	c := coral.NewSimple(func(k int) (v int, expire *time.Time, err error) {
		count++
		time.Sleep(time.Second / 100)
		v = k * 100
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

	c.Reset()
	c.Get(1)
	c.Get(2)

	if count != 5 {
		t.Error(`simple clean fail`)
	}

	c.Set(3, 123123, nil)
	c.Get(3)

	if count != 5 {
		t.Error(`simple set fail`)
	}
}

func TestSimpleExpire(t *testing.T) {

	var prev uint64

	loadFn := func(k int) (v uint64, expire *time.Time, err error) {
		e := time.Now().Add(time.Second / 200)
		expire = &e
		for {
			r := rand.Uint64()
			if r != prev {
				v = r
				prev = r
				break
			}
		}
		return
	}

	c := coral.NewSimple(loadFn)

	v, _ := c.Get(1)

	i := v

	v, _ = c.Get(1)
	if i != v {
		t.Error(`cache not work`)
	}
	time.Sleep(time.Second / 100)

	v, _ = c.Get(1)

	if i == v {
		t.Error(`expire not work`)
	}
}

func TestSimpleDeadline(t *testing.T) {

	/*
		loadFn := func(k interface{}) (v interface{}, expire *time.Time, err error) {
			time.Sleep(time.Second / 100)
			v = k.(int) * 100
			return
		}

		c := coral.NewSimple(coral.LoadDeadline(loadFn, time.Second/150))
		_, err := c.Get(1)
		if err != coral.ErrLoadFuncTimeout {
			t.Error(`LoadDeadline fail`)
		}

		c = coral.BuildSimple(coral.LoadDeadline(loadFn, time.Second/50))
		_, err = c.Get(1)
		if err == coral.ErrLoadFuncTimeout {
			t.Error(`LoadDeadline fail`)
		}
	*/
}
