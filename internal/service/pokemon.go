package service

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/HardDie/tg_bot_actions/internal/logger"
	"github.com/HardDie/tg_bot_actions/internal/models"
	"github.com/HardDie/tg_bot_actions/internal/utils"
)

const (
	pokemonSmile = "⚡️"
	hiddenImage  = `<a href="%s">&#8205;</a>`
)

var (
	pokemonTypes = map[string]string{
		"bug":      "🪲",
		"dark":     "dark",
		"dragon":   "🐉",
		"electric": "⚡️",
		"fairy":    "🧚‍♂️",
		"fighting": "👊",
		"fire":     "🔥",
		"flying":   "✈️",
		"ghost":    "👻",
		"grass":    "🍀",
		"ground":   "ground",
		"ice":      "🧊",
		"normal":   "normal",
		"poison":   "🟣",
		"physic":   "physic",
		"rock":     "🪨",
		"steel":    "steel",
		"water":    "🌊",
	}
)

type PokemonService struct {
	pokemons []models.Pokemon
}

func NewPokemonService() (*PokemonService, error) {
	s := PokemonService{}
	err := s.readPokemonsFromFile("pokemons.json")
	if err != nil {
		return nil, fmt.Errorf("error init pokemon service: %w", err)
	}
	return &s, nil
}

func (s PokemonService) GeneratePokemonIndex() int {
	return utils.Random(len(s.pokemons))
}

func (s PokemonService) GenerateDescription(index int) string {
	if len(s.pokemons) == 0 {
		logger.Error.Println("Pokemon records is empty")
		return ""
	}
	if index < 0 || index >= len(s.pokemons) {
		logger.Error.Printf("Invalid pokemon index: %d, have records: %d\n", index, len(s.pokemons))
		index = s.GeneratePokemonIndex()
	}

	pokemon := s.pokemons[index]
	return fmt.Sprintf(hiddenImage+`Какой ты покемон:
<a href="%s"><b>#%04d</b></a> %s
Тип: %s
Рост: %s (%s)
Вес: %.01f lbs (%s)`,
		pokemon.ThumbnailImage,
		pokemon.DetailPageURL,
		pokemon.ID,
		pokemon.Name,
		s.typeOfPokemon(pokemon),
		s.inchToFootInch(pokemon.Height), s.inchToCm(pokemon.Height),
		pokemon.Weight, s.lbsToKg(pokemon.Weight),
	)
}

func (s *PokemonService) readPokemonsFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error open %s: %w", filename, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Warn.Printf("error closing file description %s: %v", filename, err.Error())
		}
	}()

	err = json.NewDecoder(file).Decode(&s.pokemons)
	if err != nil {
		return fmt.Errorf("error parse %s file: %w", filename, err)
	}

	return nil
}
func (s PokemonService) typeOfPokemon(p models.Pokemon) string {
	var res []string
	for _, t := range p.Type {
		res = append(res, pokemonTypes[t])
	}
	return strings.Join(res, " ")
}
func (s PokemonService) lbsToKg(val float32) string {
	return fmt.Sprintf("%.01f кг", float32(math.Round(float64(val/2.205))))
}
func (s PokemonService) inchToFootInch(val float32) string {
	foot := int(val) / 12
	inch := int(val - float32(foot*12))
	if foot > 0 {
		return fmt.Sprintf("%d' %02d''", foot, inch)
	}
	return fmt.Sprintf("%02d''", inch)
}
func (s PokemonService) inchToCm(val float32) string {
	return fmt.Sprintf("%.01f см", math.Round(float64(val*2.54)))
}
