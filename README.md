# timerpool
timer pool for high performance time out control (pooled time.AfterFunc/time.NewTimer)

use example with context
```go
ctx, cancel := context.WithCancel(parentCtx)
t := getTimerWithAfterFunc(time.Second*5, cancel)
putTimerWithAfterFunc(t)
cancel()
```
use example with timer.C
```go
timer:=getTimer(time.Second*5)
select {
case <-timer.C:
	// timeout
	putTimer(timer)
	return
	// other case ...
}
putTimer(timer)
```
benchmark result
```
goos: darwin
goarch: amd64
BenchmarkTimeoutControlViaContextWithTimeout-8           2251708               518 ns/op
BenchmarkTimeoutControlViaRawTimer-8                     4401963               267 ns/op
BenchmarkTimeoutControlViaTimerPool-8                    4477863               257 ns/op
```
