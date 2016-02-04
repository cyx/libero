package libero

import (
	"testing"
	"time"

	"github.com/apg/ln"
	"github.com/rcrowley/go-metrics"
)

func init() {
	var _ ln.Filter = ln.FilterFunc(Librato)
}

func TestIntercept(t *testing.T) {
	oldFilters := ln.DefaultLogger.Filters
	ln.DefaultLogger.Filters = []ln.Filter{ln.FilterFunc(Librato)}
	defer func() {
		ln.DefaultLogger.Filters = oldFilters
	}()

	cases := []struct {
		input  string
		output string
		value  interface{}
	}{
		{"count#hello.count", "hello.count", 1},
		{"sample#hello.sample", "hello.sample", 2},
		{"measure#hello.measure", "hello.measure", 2},
		{"measure#hello.measure", "hello.measure", time.Second},
		{"gauge#hello.gauge", "hello.gauge", 1001},
	}

	for _, test := range cases {
		ln.Info(ln.F{test.input: test.value})
		if metrics.Get(test.output) == nil {
			t.Fatalf("Expected metric %s to be registered, got none", test.output)
		}
	}
}

func TestMetricName(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{"count#hello.count", "hello.count"},
		{"sample#hello.sample", "hello.sample"},
		{"measure#hello.measure", "hello.measure"},
		{"gauge#hello.gauge", "hello.gauge"},
	}

	for _, test := range cases {
		if m := metricName(test.input); m != test.output {
			t.Fatalf("Expected %s != %s", m, test.output)
		}
	}
}

func BenchmarkLibero(b *testing.B) {
	oldFilters := ln.DefaultLogger.Filters
	ln.DefaultLogger.Filters = []ln.Filter{ln.FilterFunc(Librato)}
	defer func() {
		ln.DefaultLogger.Filters = oldFilters
	}()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ln.Info(ln.F{"count#hello.count": 1})
	}
	b.StopTimer()
}

func BenchmarkLiberoBatched(b *testing.B) {
	oldFilters := ln.DefaultLogger.Filters
	ln.DefaultLogger.Filters = []ln.Filter{ln.FilterFunc(Librato)}
	defer func() {
		ln.DefaultLogger.Filters = oldFilters
	}()

	b.StartTimer()
	counter := 0
	for i := 0; i < b.N; i++ {
		counter++
		if counter%1000 == 0 {
			ln.Info(ln.F{"count#hello.count": counter})
		}
	}
	b.StopTimer()
}

func BenchmarkGoMetrics(b *testing.B) {
	counter := metrics.GetOrRegisterCounter("hello", metrics.DefaultRegistry)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		counter.Inc(1)
	}
	b.StopTimer()
}
