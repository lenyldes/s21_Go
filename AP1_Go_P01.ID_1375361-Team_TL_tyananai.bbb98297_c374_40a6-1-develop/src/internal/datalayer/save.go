package datalayer

import (
	"encoding/json"
	"os"
	"rogue/internal/domain"
)

const SavePath = "../../internal/datalayer/save.json"

func SaveGame(g *domain.Game, path string) error {
	data, err := json.Marshal(g)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadGame(path string) (*domain.Game, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var g domain.Game
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}
	return &g, nil
}
