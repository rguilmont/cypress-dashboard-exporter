package optional

import (
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

// I NEED GENERICS PLEEEAAAASE :)
type Nothing struct{}

// Optional time.Time

// Option of Time
type OptionalTime interface {
	rmAj2vb2A98wYaRuT3C2wX5TV8y67tEc()
}

type SomeTime time.Time

func (SomeTime) rmAj2vb2A98wYaRuT3C2wX5TV8y67tEc() {}

func (Nothing) rmAj2vb2A98wYaRuT3C2wX5TV8y67tEc() {}

func NewOptionalTime(t *time.Time) OptionalTime {
	if t == nil {
		return Nothing{}
	} else {
		return SomeTime(*t)
	}
}

func OrElseTime(o OptionalTime, def time.Time) time.Time {
	switch value := o.(type) {
	case Nothing:
		return def
	case SomeTime:
		return time.Time(value)
	default:
		logrus.Panicln("Impossible switch condition : unknown type %v - value %v", reflect.TypeOf(value), value)
	}
	return time.Time{} // Should never enter here. Just make it compile.
}

type OptionalInt interface {
	amAj2vb2A98wYaRuT3C2wX5TV8y67tEc()
}

type SomeInt int

func (SomeInt) amAj2vb2A98wYaRuT3C2wX5TV8y67tEc() {}

func (Nothing) amAj2vb2A98wYaRuT3C2wX5TV8y67tEc() {}

func NewOptionalInt(t *int) OptionalInt {
	if t == nil {
		return Nothing{}
	} else {
		return SomeInt(*t)
	}
}

func OrElseInt(o OptionalInt, def int) int {
	switch value := o.(type) {
	case Nothing:
		return def
	case SomeInt:
		return int(value)
	default:
		logrus.Panicln("Impossible switch condition : unknown type %v - value %v", reflect.TypeOf(value), value)
	}
	return 0 // Should never enter here. Just make it compile.
}
