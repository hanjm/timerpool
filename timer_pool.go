package timerpool

import (
	"sync"
	"time"
)

var _timerWithAfterFuncPool = sync.Pool{
}

type timerWithAfterFunc struct {
	t *time.Timer
	f func()
}

// GetTimer get a timer from pool or create from time.AfterFunc
func GetTimerWithAfterFunc(d time.Duration, f func()) *timerWithAfterFunc {
	//defer func() func() {
	//	start := time.Now()
	//	return func() {
	//		log.Printf("getTimerWithAfterFunc:%s\n", time.Since(start))
	//	}
	//}()()
	if v := _timerWithAfterFuncPool.Get(); v != nil {
		t := v.(*timerWithAfterFunc)
		t.f = f
		t.t.Reset(d)
		return t
	}
	tf := &timerWithAfterFunc{f: f}
	tf.t = time.AfterFunc(d, func() {
		tf.f()
	})
	return tf
}

// PutTimer stop a timer and return into pool
func PutTimerWithAfterFunc(t *timerWithAfterFunc) {
	//defer func() func() {
	//	start := time.Now()
	//	return func() {
	//		log.Printf("getTimerWithAfterFunc:%s\n", time.Since(start))
	//	}
	//}()()
	t.t.Stop()
	// time.AfterFunc使用的timer.C是nil, 不需要clean
	_timerWithAfterFuncPool.Put(t)
}

var _timerPool = sync.Pool{
}

// GetTimer get a timer from pool or create from time.NewTimer
func GetTimer(d time.Duration) *time.Timer {
	if v := _timerPool.Get(); v != nil {
		t := v.(*time.Timer)
		t.Reset(d)
		return t
	}
	return time.NewTimer(d)
}

// PutTimer stop a timer and return into pool
func PutTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	_timerPool.Put(t)
}
