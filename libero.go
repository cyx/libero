package libero

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/apg/ln"
	"github.com/rcrowley/go-metrics"
)

var (
	defaultSampleReservoirSize = 100
	samplesRegistry            = make(map[string]metrics.Sample)
)

func Librato(e ln.Event) bool {
	found := false

	for k, v := range e.Data {
		if strings.HasPrefix(k, "count#") {
			found = update("count", metricName(k), v)
		}
		if strings.HasPrefix(k, "sample#") {
			found = update("sample", metricName(k), v)
		}
		if strings.HasPrefix(k, "measure#") {
			found = update("sample", metricName(k), v)
		}
		if strings.HasPrefix(k, "gauge#") {
			found = update("gauge", metricName(k), v)
		}
	}
	return !found
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
		metrics.GetOrRegisterHistogram(metric, metrics.DefaultRegistry, getOrRegisterSample(metric)).Update(n)
	case "measure":
		metrics.GetOrRegisterTimer(metric, metrics.DefaultRegistry).Update(time.Duration(n))
	case "gauge":
		metrics.GetOrRegisterGauge(metric, metrics.DefaultRegistry).Update(n)
	default:
		log.Printf("librato.update err: Unknown kind %s", kind)
		return false
	}

	return true
}

func cast(v interface{}) (n int64, err error) {
	switch v.(type) {
	case int32, int:
		n = int64(v.(int))
	case int64:
		n = v.(int64)
	case time.Duration:
		n = int64(v.(time.Duration))
	default:
		err = fmt.Errorf("Unable to cast %v %T to int64/int/time.Duration", v, v)
	}
	return n, err
}

func metricName(key string) string {
	return key[strings.Index(key, "#")+1:]
}

func getOrRegisterSample(metric string) metrics.Sample {
	if samplesRegistry[metric] == nil {
		samplesRegistry[metric] = metrics.NewUniformSample(defaultSampleReservoirSize)
	}
	return samplesRegistry[metric]
}
