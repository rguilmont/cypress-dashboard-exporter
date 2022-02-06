package metricsmap

import (
	"fmt"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rguilmont/cypress-dashboard-exporter/pkg/cypresscollector/converter"
	"github.com/sirupsen/logrus"
)

type Key struct {
	Prom *prometheus.Desc
	Hash labelsHash
}

type Value struct {
	Value      float64
	Labels     []string
	updated_at time.Time
}

func keyNotFoundError(k Key) error {
	return fmt.Errorf("Cannot find the key %v", k)
}

// MetricMap is an interface to allow to store metrics, and then
//  get them passed to prometheus.
//  The idea behind is that, for some of these, we want to keep only the first value,
//  whereas for other the sum, or ever something else in the future.
type MetricMap interface {
	// Add to the map
	Add(*prometheus.Desc, interface{}, ...string)

	// Get value from the map
	Get(Key) (*Value, error)
	// Map returns the map
	Map() map[Key]Value
	// Free old items, to not keep metrics from branches or web browser that hasn't been found for some time
	FreeOldItems()
}

// For Gauge value
type MetricMapKeepFirst struct {
	metrics   map[Key]Value
	KeepUntil time.Duration
}

func (m *MetricMapKeepFirst) Add(k *prometheus.Desc, value interface{}, labels ...string) {
	if m.metrics == nil {
		m.metrics = map[Key]Value{}
	}

	labelsHash := StringSliceHash(labels)
	v, err := converter.ConvertValueForPrometheus(value)
	if err != nil {
		logrus.Errorf("Can't convert metric %v type %v to float64 - original error: %v", k.String(), reflect.TypeOf(value), err)
	}
	key := Key{
		k,
		labelsHash,
	}
	m.metrics[key] = Value{v, labels, time.Now()}
}

func (m MetricMapKeepFirst) Get(k Key) (*Value, error) {
	if v, ok := m.metrics[k]; ok {
		return &v, nil
	}
	return nil, keyNotFoundError(k)
}

func (m MetricMapKeepFirst) Map() map[Key]Value {
	freeOldItems(m.metrics, m.KeepUntil)
	return m.metrics
}

// For Summary value
type MetricMapSumValues struct {
	metrics   map[Key]Value
	KeepUntil time.Duration
}

func (m *MetricMapSumValues) Add(k *prometheus.Desc, value interface{}, labels ...string) {
	if m.metrics == nil {
		m.metrics = map[Key]Value{}
	}

	labelsHash := StringSliceHash(labels)
	v, err := converter.ConvertValueForPrometheus(value)
	if err != nil {
		logrus.Errorf("Can't convert metric %v type %v to float64 - original error: %v", k.String(), reflect.TypeOf(value), err)
	}
	key := Key{
		k,
		labelsHash,
	}
	if currentValue, ok := m.metrics[key]; ok {
		m.metrics[key] = Value{currentValue.Value + v, labels, time.Now()}
		return
	}
	m.metrics[key] = Value{v, labels, time.Now()}
}

func (m MetricMapSumValues) Get(k Key) (*Value, error) {
	if v, ok := m.metrics[k]; ok {
		return &v, nil
	}
	return nil, keyNotFoundError(k)
}

func (m MetricMapSumValues) Map() map[Key]Value {
	freeOldItems(m.metrics, m.KeepUntil)
	return m.metrics
}

func freeOldItems(m map[Key]Value, keepUntil time.Duration) {
	for k, v := range m {
		if v.updated_at.Add(keepUntil).Before(time.Now()) {
			logrus.Debugf("Removing entry from map %v : %v\n", k, v)
			delete(m, k)
		}
	}
}
