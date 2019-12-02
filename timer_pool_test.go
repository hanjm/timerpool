package timerpool

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
go test -run TimeoutControl -bench TimeoutControl
goos: darwin
goarch: amd64
BenchmarkTimeoutControlViaContextWithTimeout-8                   2272261               523 ns/op
BenchmarkTimeoutControlViaContextWithCancelAndRawTimer-8         4370292               267 ns/op
BenchmarkTimeoutControlViaContextWithCancelAndTimerPool-8        4662602               253 ns/op
BenchmarkTimeoutControlViaRawTimer-8                             4829988               245 ns/op
BenchmarkTimeoutControlViaTimerPool-8                            9144037               126 ns/op
*/

func BenchmarkTimeoutControlViaContextWithTimeout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, cancel := context.WithTimeout(context.Background(), time.Second*5)
		cancel()
	}
}

func BenchmarkTimeoutControlViaContextWithCancelAndRawTimer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, cancel := context.WithCancel(context.Background())
		t := time.AfterFunc(time.Second*5, cancel)
		t.Stop()
		cancel()
	}
}

func BenchmarkTimeoutControlViaContextWithCancelAndTimerPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, cancel := context.WithCancel(context.Background())
		t := GetTimerWithAfterFunc(time.Second*5, cancel)
		PutTimerWithAfterFunc(t)
		cancel()
	}
}

func BenchmarkTimeoutControlViaRawTimer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer := time.NewTimer(time.Second * 5)
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}
}

func BenchmarkTimeoutControlViaTimerPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer := GetTimer(time.Second * 5)
		PutTimer(timer)
	}
}

func TestGetPutTimerWithAfterFunc(t *testing.T) {
	a := assert.New(t)
	var v int
	timer := GetTimerWithAfterFunc(time.Second, func() {
		v = 1
	})
	time.Sleep(time.Second)
	PutTimerWithAfterFunc(timer)
	timer = GetTimerWithAfterFunc(time.Second*2, func() {
		v = 2
	})
	time.Sleep(time.Second * 2)
	PutTimerWithAfterFunc(timer)
	a.Equal(2, v)
}

func BenchmarkGetPutTimerWithAfterFunc(b *testing.B) {
	runtime.GOMAXPROCS(10)
	b.SetParallelism(1)
	b.N = 1000
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			timer := GetTimerWithAfterFunc(time.Second, func() {
			})
			PutTimerWithAfterFunc(timer)
		}
	})
}

func TestGetPutTimer(t *testing.T) {
	a := assert.New(t)
	d := time.Second
	timer := GetTimer(d)
	select {
	case <-timer.C:
	case <-time.After(d + time.Millisecond):
		a.Fail("unexpected")
	}
	PutTimer(timer)
	//
	d = time.Second * 2
	timer = GetTimer(d)
	select {
	case <-timer.C:
	case <-time.After(d + time.Millisecond):
		a.Fail("unexpected")
	}
}
