package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HardDie/tg_bot_actions/internal/logger"
	"github.com/HardDie/tg_bot_actions/internal/models"
	"github.com/HardDie/tg_bot_actions/internal/utils"
)

const (
	criminalBook = "📕"
)

type CriminalService struct {
	criminals []models.Criminal
}

func NewCriminalService() (*CriminalService, error) {
	s := CriminalService{}
	err := s.readLawFromFile("criminals.json")
	if err != nil {
		return nil, fmt.Errorf("error init criminal service: %w", err)
	}
	return &s, nil
}

func (s CriminalService) GenerateCriminalIndex() int {
	return utils.Random(len(s.criminals))
}
func (s CriminalService) GenerateDescription(index int) string {
	if len(s.criminals) == 0 {
		logger.Error.Println("Criminal records is empty")
		return ""
	}
	if index < 0 || index >= len(s.criminals) {
		logger.Error.Printf("Invalid criminal index: %d, have records: %d\n", index, len(s.criminals))
		index = s.GenerateCriminalIndex()
	}

	return fmt.Sprintf(`%s <u>Твоя статья УК РФ</u>:
<a href="%s"><b>%s</b></a> - %s`,
		criminalBook,
		s.criminals[index].Link,
		s.criminals[index].Number,
		s.criminals[index].Description,
	)
}

func (s *CriminalService) readLawFromFile(filename string) error {
	file, err := os.Open("criminals.json")
	if err != nil {
		return fmt.Errorf("error open %s: %w", filename, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Warn.Printf("error closing file description %s: %v", filename, err.Error())
		}
	}()

	err = json.NewDecoder(file).Decode(&s.criminals)
	if err != nil {
		return fmt.Errorf("error parse %s file: %w", filename, err)
	}

	return nil
}
