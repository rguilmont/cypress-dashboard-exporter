package metricsmap

import (
	"reflect"
	"testing"
	"time"
)

func TestMetricMapSumValues_Map(t *testing.T) {
	type fields struct {
		metrics   map[Key]Value
		keepUntil time.Duration
	}
	now := time.Now()
	tests := []struct {
		name   string
		fields fields
		want   map[Key]Value
	}{
		{
			name: "test old items",
			fields: fields{
				metrics: map[Key]Value{
					// Item to keep
					Key{nil, "future"}: Value{
						0.0,
						[]string{},
						now,
					},
					// Too old item that should be removed
					Key{nil, "old"}: Value{
						0.0,
						[]string{},
						now.Add(-time.Hour * 10),
					},
				},
				keepUntil: time.Duration(2 * time.Hour),
			},
			want: map[Key]Value{
				Key{nil, "future"}: Value{
					0.0,
					[]string{},
					now,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MetricMapSumValues{
				metrics:   tt.fields.metrics,
				KeepUntil: tt.fields.keepUntil,
			}
			if got := m.Map(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricMapSumValues.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
