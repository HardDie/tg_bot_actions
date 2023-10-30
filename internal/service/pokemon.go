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
	hiddenImage = `<a href="%s">&#8205;</a>`
)

var (
	pokemonTypes = map[string]string{
		"bug":      "🪲",
		"dark":     "😈",
		"dragon":   "🐲",
		"electric": "⚡️",
		"fairy":    "🧚‍♂️",
		"fighting": "👊",
		"fire":     "🔥",
		"flying":   "🕊️",
		"ghost":    "👻",
		"grass":    "🌱",
		"ground":   "🪱",
		"ice":      "🧊",
		"normal":   "🗿",
		"poison":   "🧪",
		"psychic":  "🧠",
		"rock":     "⛰",
		"steel":    "⚙️",
		"water":    "🌊",

		"weed": "🤙",
	}
)

type PokemonService struct {
	pokemons     []models.Pokemon
	rarePokemons []models.Pokemon
}

func NewPokemonService() (*PokemonService, error) {
	s := PokemonService{}
	err := s.readPokemonsFromFile("pokemons.json")
	if err != nil {
		return nil, fmt.Errorf("error init pokemon service: %w", err)
	}
	s.initRarePokemons()
	return &s, nil
}

func (s PokemonService) GeneratePokemonIndex() (int, bool) {
	// Rare pokemon with %5 chance
	if utils.Random(100) < 5 {
		return s.generateRandomRareIndex(), true
	}
	return s.generateRandomIndex(), false
}
func (s PokemonService) generateRandomRareIndex() int {
	return utils.Random(len(s.rarePokemons))
}
func (s PokemonService) generateRandomIndex() int {
	return utils.Random(len(s.pokemons))
}

func (s PokemonService) GenerateDescription(index int, isRare bool) string {
	if isRare {
		return s.generateDescriptionForRare(index)
	}
	if len(s.pokemons) == 0 {
		logger.Error.Println("Pokemon records is empty")
		return ""
	}
	if index < 0 || index >= len(s.pokemons) {
		logger.Error.Printf("Invalid pokemon index: %d, have records: %d\n", index, len(s.pokemons))
		index = s.generateRandomIndex()
	}

	pokemon := s.pokemons[index]
	return fmt.Sprintf(hiddenImage+`Какой ты покемон:
<a href="%s"><b>#%04d</b></a> %s
Тип: %s
Рост: %s (%s)
Вес: %.01f lbs (%s)
Поколение: %s
Регион: <a href="%s"><b>%s</b></a>`,
		pokemon.ThumbnailImage,
		pokemon.DetailPageURL,
		pokemon.ID,
		pokemon.Name,
		s.typeOfPokemon(pokemon),
		s.inchToFootInch(pokemon.Height), s.inchToCm(pokemon.Height),
		pokemon.Weight, s.lbsToKg(pokemon.Weight),
		s.intToRoman(pokemon.Generation),
		pokemon.RegionLink, pokemon.Region,
	)
}

func (s PokemonService) generateDescriptionForRare(index int) string {
	if len(s.rarePokemons) == 0 {
		logger.Error.Println("Rare pokemon records is empty")
		return ""
	}
	if index < 0 || index >= len(s.rarePokemons) {
		logger.Error.Printf("Invalid rare pokemon index: %d, have records: %d\n", index, len(s.rarePokemons))
		index = s.generateRandomRareIndex()
	}

	pokemon := s.rarePokemons[index]
	return fmt.Sprintf(hiddenImage+`Какой ты покемон:
[RARE] <b>#????</b> %s
Тип: %s
Рост: %s (%s)
Вес: %.01f lbs (%s)`,
		pokemon.ThumbnailImage,
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
func (s *PokemonService) initRarePokemons() {
	s.rarePokemons = append(s.rarePokemons,
		models.Pokemon{
			Name:           "Травозавр",
			Type:           []string{"grass", "poison", "weed"},
			ThumbnailImage: "https://i.imgur.com/xySj0Vs.png",
			Weight:         155,
			Height:         68,
		},
		models.Pokemon{
			Name:           "#420",
			Type:           []string{"grass", "poison", "weed"},
			ThumbnailImage: "https://i.imgur.com/acP3eWP.jpeg",
			Weight:         220,
			Height:         74,
		},
		models.Pokemon{
			Name:           "Bubblehash",
			Type:           []string{"normal", "weed"},
			ThumbnailImage: "https://i.imgur.com/6X6SA2m.jpeg",
			Weight:         143,
			Height:         67,
		},
		models.Pokemon{
			Name:           "Honey Pot",
			Type:           []string{"fire", "weed"},
			ThumbnailImage: "https://i.imgur.com/Fiszpe7.jpg",
			Weight:         1468,
			Height:         69,
		},
		models.Pokemon{
			Name:           "Pizza",
			Type:           []string{"fairy", "weed"},
			ThumbnailImage: "https://i.imgur.com/zntAPFt.jpg",
			Weight:         151,
			Height:         68,
		},
		models.Pokemon{
			Name:           "Pikushu",
			Type:           []string{"electric", "weed"},
			ThumbnailImage: "https://i.imgur.com/tyWpIjo.jpg",
			Weight:         198,
			Height:         74,
		},
	)
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
func (s PokemonService) intToRoman(num int) string {
	var roman string = ""
	var numbers = []int{1, 4, 5, 9, 10, 40, 50, 90, 100, 400, 500, 900, 1000}
	var romans = []string{"I", "IV", "V", "IX", "X", "XL", "L", "XC", "C", "CD", "D", "CM", "M"}
	var index = len(romans) - 1
	for num > 0 {
		for numbers[index] <= num {
			roman += romans[index]
			num -= numbers[index]
		}
		index -= 1
	}
	return roman
}
