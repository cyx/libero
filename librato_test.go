package ln

import "testing"

func init() {
	var _ Filter = FilterFunc(Librato)
}

func TestIntercept(t *testing.T) {
	oldFilters := DefaultLogger.Filters
	DefaultLogger.Filters = []Filter{FilterFunc(Librato)}
	defer func() {
		DefaultLogger.Filters = oldFilters
	}()

	Info(F{"count#hello.counter": 1})
	Info(F{"sample#hello.sample": 2})
	Info(F{"measure#hello.measure": 2})
}

func TestMetricName(t *testing.T) {
	if m := metricName("count#hello"); m != "hello" {
		t.Fatalf("Expected %s != %s", m, "hello")
	}
}
