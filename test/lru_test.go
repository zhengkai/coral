package test

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zhengkai/coral"
)

func TestLRUConcurrency(t *testing.T) {

	for j := 0; j < 100; j++ {

		concurrency := 20

		var mux sync.RWMutex
		mux.Lock()

		var count uint32

		var wg sync.WaitGroup
		wg.Add(concurrency)

		c := coral.BuildLRU(func(k interface{}) (v interface{}, expire *time.Time, err error) {
			atomic.AddUint32(&count, 1)
			v = k.(int) * 100
			return
		}, 10, 10)

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

func TestLRUConcurrencyDelay(t *testing.T) {
	// t.Fatal("not implemented")

	var mux sync.RWMutex
	mux.Lock()

	var count uint32

	var wg sync.WaitGroup
	wg.Add(100)

	c := coral.BuildLRU(func(k interface{}) (v interface{}, expire *time.Time, err error) {
		atomic.AddUint32(&count, 1)
		time.Sleep(time.Second / 100)
		v = k.(int) * 100
		return
	}, 200, 200)

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

func TestLRUMisc(t *testing.T) {

	var count int

	c := coral.BuildLRU(func(k interface{}) (v interface{}, expire *time.Time, err error) {
		count++
		time.Sleep(time.Second / 100)
		v = k.(int) * 100
		return
	}, 10, 10)

	c.StatsOff()

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

	c.Set(3, true, nil)
	c.Get(3)

	if count != 5 {
		t.Error(`simple set fail`)
	}
}

func TestLRUExpire(t *testing.T) {

	var prev uint64

	loadFn := func(k interface{}) (v interface{}, expire *time.Time, err error) {
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

	c := coral.BuildLRU(loadFn, 10, 10)

	v, _ := c.Get(1)

	i := v.(uint64)

	v, _ = c.Get(1)
	if i != v.(uint64) {
		t.Error(`cache not work`)
	}
	time.Sleep(time.Second / 100)

	v, _ = c.Get(1)

	if i == v.(uint64) {
		t.Error(`expire not work`)
	}
}

func TestLRUDeadline(t *testing.T) {

	loadFn := func(k interface{}) (v interface{}, expire *time.Time, err error) {
		time.Sleep(time.Second / 100)
		v = k.(int) * 100
		return
	}

	c := coral.BuildLRU(coral.LoadDeadline(loadFn, time.Second/150), 10, 10)
	_, err := c.Get(1)
	if err != coral.ErrLoadFuncTimeout {
		t.Error(`LoadDeadline fail`)
	}

	c = coral.BuildLRU(coral.LoadDeadline(loadFn, time.Second/50), 10, 10)
	_, err = c.Get(1)
	if err == coral.ErrLoadFuncTimeout {
		t.Error(`LoadDeadline fail`)
	}
}

func TestLRUBuild(t *testing.T) {
	var cnt uint32
	loadFn := func(k interface{}) (v interface{}, expire *time.Time, err error) {
		if k.(int) == 10 {
			atomic.AddUint32(&cnt, 1)
		}
		v = 1
		return
	}
	c := coral.BuildLRU(loadFn, 0, 0)
	for i := 0; i < 1051; i++ {
		c.Get(i)
	}
	for i := 0; i < 1051; i++ {
		c.Get(i)
	}

	if cnt != 2 {
		t.Error(`build default capacity not 1000`)
	}

	c.Clean()
	for i := 1000; i > 0; i-- {
		c.Get(i)
	}
	if c.Stats().Hit != 0 {
		t.Error(`clean fail`)
	}
	for i := 0; i < 2000; i++ {
		c.Get(i)
	}
	if c.Stats().Hit != 1000 {
		t.Error(`stats hit fail`)
	}
	c.Stats().Report()

	c.StatsOff()
	for i := 0; i < 2000; i++ {
		c.Get(i)
	}
	c.Stats().Report()
}

func TestLRUEviction(t *testing.T) {

	var cnt uint32

	loadFn := func(k interface{}) (v interface{}, expire *time.Time, err error) {
		if k.(int) == 5 {
			atomic.AddUint32(&cnt, 1)
		}
		v = 1
		return
	}

	c := coral.BuildLRU(loadFn, 10, 15)

	for i := 0; i < 12; i++ {
		c.Get(i)
	}
	for i := 0; i < 12; i++ {
		c.Get(i)
	}

	if cnt != 1 {
		t.Error(`eviction when not overflow`)
	}
	for i := 18; i >= 0; i-- {
		c.Get(i)
	}
	if cnt != 2 {
		t.Error(`eviction when not overflow`)
	}

	c.Set(5, 1, nil)
	for i := 0; i < 18; i++ {
		c.Get(i)
	}
	if cnt != 2 {
		t.Error(`eviction when not overflow`)
	}
}

func TestLRUSetSlim(t *testing.T) {

	var cnt uint32

	loadFn := func(k interface{}) (v interface{}, expire *time.Time, err error) {
		if k.(int) == 0 {
			atomic.AddUint32(&cnt, 1)
		}
		e := time.Now().Add(time.Second / 200)
		expire = &e
		v = 1
		return
	}

	c := coral.BuildLRU(loadFn, 10, 12)
	for i := 0; i < 12; i++ {
		c.Get(i)
	}

	time.Sleep(time.Second / 100)
	c.Set(20, 1, nil)
	c.Get(0)
	if cnt != 2 {
		t.Error(`slim error when set overflow`)
	}

	for i := 1; i < 15; i++ {
		c.Get(i)
		c.Get(0)
	}

	if cnt != 2 {
		t.Error(`LRU not work`)
	}

	c.Stats().Report()
}
