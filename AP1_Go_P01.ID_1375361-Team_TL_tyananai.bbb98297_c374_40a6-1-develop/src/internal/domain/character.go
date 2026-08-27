package domain

import "rogue/internal/domain/entity"

type Character struct {
	Pos       entity.Point
	Symbol    rune
	Dexterity int
	Strength  int
	Health    int
	MaxHealth int
}
