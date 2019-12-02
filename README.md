[![GoDoc](https://godoc.org/github.com/hanjm/timerpool?status.svg)](https://godoc.org/github.com/hanjm/timerpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/hanjm/timerpool)](https://goreportcard.com/report/github.com/hanjm/timerpool)
[![code-coverage](http://gocover.io/_badge/github.com/hanjm/timerpool)](http://gocover.io/github.com/hanjm/timerpool)

# timerpool
timer pool for high performance time out control (pooled time.AfterFunc/time.NewTimer)

use example with context
```go
ctx, cancel := context.WithCancel(parentCtx)
t := GetTimerWithAfterFunc(time.Second*5, cancel)
PutTimerWithAfterFunc(t)
cancel()
```
use example with timer.C
```go
timer:=GetTimer(time.Second*5)
select {
case <-timer.C:
	// timeout
	PutTimer(timer)
	return
	// other case ...
}
PutTimer(timer)
```
benchmark result
```
go test -run TimeoutControl -bench TimeoutControl
goos: darwin
goarch: amd64
BenchmarkTimeoutControlViaContextWithTimeout-8                   2272261               523 ns/op
BenchmarkTimeoutControlViaContextWithCancelAndRawTimer-8         4370292               267 ns/op
BenchmarkTimeoutControlViaContextWithCancelAndTimerPool-8        4662602               253 ns/op
BenchmarkTimeoutControlViaRawTimer-8                             4829988               245 ns/op
BenchmarkTimeoutControlViaTimerPool-8                            9144037               126 ns/op
```
