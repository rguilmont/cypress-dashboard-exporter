package metricsmap

import "fmt"

type labelsHash string

func StringSliceHash(s []string) labelsHash {

	toHash := ""
	separator := "___"

	for _, i := range s {
		toHash = fmt.Sprintf("%v%v%v", i, separator, toHash)
	}
	return labelsHash(toHash)

}
