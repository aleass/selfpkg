package main

import (
	ants "github.com/panjf2000/ants/v2"
	"runtime"
	"sync"
	"sync/atomic"
)

func pt(i int) {

}

type originSpinLock uint32

func (sl *originSpinLock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		runtime.Gosched()
	}
}
func (sl *originSpinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

func selff() {
	var s sync.WaitGroup
	var j int64
	s.Add(100000)
	var o originSpinLock
	var sc = sync.NewCond(&o)
	for i := 0; i < 100000; i++ {
		go func() {
			for true {
				if atomic.LoadInt64(&j) < 4 {
					atomic.AddInt64(&j, 1)
					_ = new(int)
					atomic.AddInt64(&j, -1)
					sc.Signal()
					s.Done()
					return
				} else {
					sc.L.Lock()
					sc.Wait()
					sc.L.Unlock()
				}
			}
		}()
	}

	s.Wait()
}
func antsf() { //6.79
	var s sync.WaitGroup
	s.Add(100000)
	f, _ := ants.NewPoolWithFunc(4, func(i interface{}) {
		_ = new(int)
		s.Done()
	})
	for i := 0; i < 100000; i++ {
		f.Invoke(i)
	}
	s.Wait()
}

func forf() {
	for i := 0; i < 100000; i++ {
		_ = new(int)
	}
}

func chanf() {
	var ic = make(chan int, 1000)
	for i := 0; i < 3; i++ {
		go func() {
			for range ic {
				_ = new(int)
			}
		}()
	}
	for i := 0; i < 100000; i++ {
		ic <- i
	}
	close(ic)
	for range ic {
		_ = new(int)
	}
}
