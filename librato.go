package libero

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"
)

var DefaultSample = metrics.NewUniformSample(100)

func Librato(e Event) bool {
	for k, v := range e.Data {
		if strings.HasPrefix(k, "count#") {
			return update("count", metricName(k), v)
		}
		if strings.HasPrefix(k, "sample#") {
			return update("sample", metricName(k), v)
		}
		if strings.HasPrefix(k, "measure#") {
			return update("sample", metricName(k), v)
		}
	}
	return false
}

func update(kind, metric string, v interface{}) bool {
	n, err := cast(v)
	if err != nil {
		log.Printf("librato.update err: %s", err)
		return false
	}

	switch kind {
	case "count":
		metrics.GetOrRegisterCounter(metric, metrics.DefaultRegistry).Inc(n)
	case "sample":
		metrics.GetOrRegisterHistogram(metric, metrics.DefaultRegistry, DefaultSample).Update(n)
	case "measure":
		metrics.GetOrRegisterTimer(metric, metrics.DefaultRegistry).Update(time.Duration(n))
	default:
		log.Printf("librato.update err: Unknown kind %s", kind)
		return false
	}

	return true
}

func cast(v interface{}) (int64, error) {
	if v, ok := v.(int64); ok {
		return v, nil
	}
	return 0, fmt.Errorf("Unable to cast %v %T to int64", v, v)
}

func metricName(key string) string {
	return key[strings.Index(key, "#")+1:]
}
