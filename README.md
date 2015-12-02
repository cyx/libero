# libero

A quick filter for use with ln that complies to l2met.

## example

```go
// Will send counter / histogram / timer metrics to go-metrics
ln.Info(ln.F({
	"count#hello.counter": 1,
	"sample#hello.sample": 1,
	"measure#hello.timer": 1,
}))
```
