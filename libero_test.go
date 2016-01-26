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
