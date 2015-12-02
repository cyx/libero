package libero

import (
	"testing"
	"time"

	"github.com/apg/ln"
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

	ln.Info(ln.F{"count#hello.counter": 1})
	ln.Info(ln.F{"sample#hello.sample": 2})
	ln.Info(ln.F{"measure#hello.measure": 2})
	ln.Info(ln.F{"measure#hello.measure": time.Second})
}

func TestMetricName(t *testing.T) {
	if m := metricName("count#hello"); m != "hello" {
		t.Fatalf("Expected %s != %s", m, "hello")
	}
}
