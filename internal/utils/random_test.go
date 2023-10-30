package utils

import (
	"fmt"
	"math"
	"testing"
)

func TestRandom(t *testing.T) {
	min, max := math.MaxInt, 0
	for i := 0; i < 100_000; i++ {
		val := Random(100)
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
	fmt.Println("min:", min)
	fmt.Println("max:", max)
}
