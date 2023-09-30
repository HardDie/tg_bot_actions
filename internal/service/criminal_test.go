package service

import "testing"

func TestCriminalValue(t *testing.T) {
	criminal, err := NewCriminalService()
	if err != nil {
		t.Fatal(err)
	}
	uniq := make(map[int]struct{})
	for i := 0; i < 100_000; i++ {
		uniq[criminal.GenerateCriminalIndex()] = struct{}{}
	}
	t.Log("criminal laws count:", len(uniq))
}
