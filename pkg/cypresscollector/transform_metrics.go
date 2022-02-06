package cypresscollector

// metricTransformer allows to transform the format of a metric ( such as converting ms to seconds )
type metricTransformer func(float64) float64

func noopTransformer(v float64) float64 {
	return v
}

func msToSec(v float64) float64 {
	return v / 1000
}

func promValueFromState(value string, expected string) float64 {
	if value == expected {
		return 1.0
	} else {
		return 0.0
	}
}
