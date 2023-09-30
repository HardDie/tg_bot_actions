package service

import "testing"

func TestGayMeterValue(t *testing.T) {
	gayMeter := NewGayMeterService()
	minimal, maximum := 100, 0
	for i := 0; i < 100_000; i++ {
		size := gayMeter.GenerateValue()
		if size < minimal {
			minimal = size
		}
		if size > maximum {
			maximum = size
		}
	}
	t.Log("min:", minimal, "max:", maximum)
}
