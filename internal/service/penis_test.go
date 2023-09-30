package service

import "testing"

func TestPenisSadSmile(t *testing.T) {
	penis := NewPenisService()

	uniq := make(map[string]struct{})
	for i := 0; i < 100_000; i++ {
		uniq[penis.getSmile(1)] = struct{}{}
	}
	t.Log("sad smiles count:", len(uniq))
}

func TestPenisHappySmile(t *testing.T) {
	penis := NewPenisService()

	uniq := make(map[string]struct{})
	for i := 0; i < 100_000; i++ {
		uniq[penis.getSmile(20)] = struct{}{}
	}
	t.Log("happy smiles count:", len(uniq))
}

func TestPenisNouns(t *testing.T) {
	penis := NewPenisService()

	uniq := make(map[string]struct{})
	for i := 0; i < 100_000; i++ {
		uniq[penis.getNoun()] = struct{}{}
	}
	t.Log("nouns count:", len(uniq))
}

func TestPenisSize(t *testing.T) {
	penis := NewPenisService()
	minimal, maximum := 100, 0
	for i := 0; i < 100_000; i++ {
		size := penis.GenerateSize()
		if size < minimal {
			minimal = size
		}
		if size > maximum {
			maximum = size
		}
	}
	t.Log("min:", minimal, "max:", maximum)
}
