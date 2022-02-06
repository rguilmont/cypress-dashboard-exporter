package converter

import (
	"fmt"
	"reflect"
	"time"
)

func ConvertValueForPrometheus(value interface{}) (float64, error) {
	var prometheusValue float64

	switch vt := value.(type) {
	case int:
		prometheusValue = float64(vt)
	case float64:
		prometheusValue = vt
	case time.Time:
		prometheusValue = float64(vt.Unix())
	case bool:
		prometheusValue = func() float64 {
			if vt {
				return 1.0
			}
			return 0.0
		}()
	default:
		return 0.0, fmt.Errorf("convertion from type %v to float ( prometheus standard ) unknown", reflect.TypeOf(vt))

	}
	return prometheusValue, nil
}

func StateToValue(state string) float64 {
	if state == "PASSED" {
		return 1.0
	}
	return 0.0
}
