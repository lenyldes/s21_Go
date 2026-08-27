package datalayer

import (
	"encoding/json"
	"os"
	"sort"
)

const LeadersPath = "../../internal/datalayer/leaders.json"

type Leaders struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Gold  int    `json:"gold"`
}

func LoadLeaders(path string) []Leaders {
	data, err := os.ReadFile(path)
	if err != nil {
		return []Leaders{}
	}
	var entries []Leaders
	if err := json.Unmarshal(data, &entries); err != nil {
		return []Leaders{}
	}
	return entries
}

func SaveLeader(name string, level, gold int) error {
	return addLeader(Leaders{Name: name, Level: level, Gold: gold}, LeadersPath)
}

func addLeader(entry Leaders, path string) error {
	entries := LoadLeaders(path)
	entries = append(entries, entry)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Level != entries[j].Level {
			return entries[i].Level > entries[j].Level
		}
		return entries[i].Gold > entries[j].Gold
	})
	if len(entries) > 10 {
		entries = entries[:10]
	}
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
