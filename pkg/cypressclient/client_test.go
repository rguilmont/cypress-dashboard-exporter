package cypressclient

import (
	"reflect"
	"testing"
)

func TestRunResults_Reverse(t *testing.T) {
	tests := []struct {
		name string
		res  RunResults
		want []RunResult
	}{
		{
			"Should reverse a simple slice",
			[]RunResult{
				{
					ID: "1",
				},
				{
					ID: "2",
				},
				{
					ID: "3",
				},
				{
					ID: "4",
				},
				{
					ID: "5",
				},
			},
			[]RunResult{
				{
					ID: "5",
				},
				{
					ID: "4",
				},
				{
					ID: "3",
				},
				{
					ID: "2",
				},
				{
					ID: "1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.res.Reverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunResults.Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
