package service

import "testing"

func TestPressureValue(t *testing.T) {
	pressure := NewPressureService()
	minSys, maxSys := 10000, 0
	minDia, maxDia := 10000, 0
	for i := 0; i < 100_000; i++ {
		sys, dia := pressure.generatePressureValues()

		if sys < minSys {
			minSys = sys
		}
		if sys > maxSys {
			maxSys = sys
		}

		if dia < minDia {
			minDia = dia
		}
		if dia > maxDia {
			maxDia = dia
		}

		if sys == dia {
			t.Fatal("sys == dia", sys, i)
		}
		if sys < dia {
			t.Fatal("sys < dia", sys, dia, i)
		}
	}
	t.Log("minSys:", minSys, "maxSys:", maxSys)
	t.Log("minDia:", minDia, "maxDia:", maxDia)
}

func TestPressurePrefixes(t *testing.T) {
	pressure := NewPressureService()

	uniq := make(map[string]struct{})
	for i := 0; i < 100_000; i++ {
		uniq[pressure.getPrefix()] = struct{}{}
	}
	t.Log("prefixes count:", len(uniq))
}
