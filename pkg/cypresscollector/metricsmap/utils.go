package metricsmap

import "github.com/prometheus/client_golang/prometheus"

func MultipleAdd(ms ...MetricMap) func(*prometheus.Desc, interface{}, ...string) {
	return func(k *prometheus.Desc, v interface{}, labels ...string) {
		for _, m := range ms {
			m.Add(k, v, labels...)
		}
	}
}
