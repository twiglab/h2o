package eyes

import (
	"math"
)

func mean(v ...float64) float64 {
	var res float64
	for _, x := range v {
		res += x
	}
	return res / float64(len(v))
}

func variance(v ...float64) float64 {
	var res float64
	var m = mean(v...)
	for _, x := range v {
		res += (x - m) * (x - m)
	}
	return res / float64(len(v))
}

func std(v ...float64) float64 {
	return math.Sqrt(variance(v...))
}
