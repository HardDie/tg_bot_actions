package service

import (
	"log"
	"testing"
)

func TestRarePokemon(t *testing.T) {
	totalCall := 0
	rareCall := 0

	pok, err := NewPokemonService()
	if err != nil {
		t.Fatal(err.Error())
	}
	for i := 0; i < 100_000; i++ {
		_, isRare := pok.GeneratePokemonIndex()
		totalCall++
		if isRare {
			rareCall++
		}
	}

	log.Printf("Calls: %d/%d\n", totalCall, rareCall)
	log.Printf("Percentage: %d%%\n", rareCall*100/totalCall)
}
