package domain

import (
	"encoding/json"

	"rogue/internal/domain/entity"
)

const (
	EnemyTypeZombie int = iota
	EnemyTypeVampire
	EnemyTypeGhost
	EnemyTypeOgre
	EnemyTypeSnake
)

type Enemy struct {
	Character
	EnemyType int
	Hostility HostilityLevel
	Alive     bool
	Behavior  EnemyBehavior
	IsVisible bool
	IsAwake   bool
}

func (e *Enemy) UnmarshalJSON(data []byte) error {
	var raw struct {
		Character
		EnemyType int             `json:"EnemyType"`
		Hostility  HostilityLevel  `json:"Hostility"`
		Alive      bool            `json:"Alive"`
		Behavior   json.RawMessage `json:"Behavior"`
		IsVisible  bool            `json:"IsVisible"`
		IsAwake    bool            `json:"IsAwake"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.Character = raw.Character
	e.EnemyType = raw.EnemyType
	e.Hostility = raw.Hostility
	e.Alive = raw.Alive
	e.IsVisible = raw.IsVisible
	e.IsAwake = raw.IsAwake

	switch raw.EnemyType {
	case EnemyTypeZombie:
		e.Behavior = &ZombieBehavior{}
	case EnemyTypeVampire:
		b := &VampireBehavior{}
		json.Unmarshal(raw.Behavior, b)
		e.Behavior = b
	case EnemyTypeGhost:
		b := &GhostBehavior{}
		json.Unmarshal(raw.Behavior, b)
		e.Behavior = b
	case EnemyTypeOgre:
		b := &OgreBehavior{}
		json.Unmarshal(raw.Behavior, b)
		e.Behavior = b
	case EnemyTypeSnake:
		b := &SnakeBehavior{}
		json.Unmarshal(raw.Behavior, b)
		e.Behavior = b
	}
	return nil
}

func (e *Enemy) Name() string {
	switch e.EnemyType {
	case EnemyTypeZombie:
		return "Zombie"
	case EnemyTypeVampire:
		return "Vampire"
	case EnemyTypeGhost:
		return "Ghost"
	case EnemyTypeOgre:
		return "Ogre"
	case EnemyTypeSnake:
		return "Snake"
	}
	return "Monster"
}

func (e *Enemy) AggroRadius() int {
	return e.Hostility.Radius()
}

func NewZombie(pos entity.Point) *Enemy {
	return &Enemy{
		Character: Character{
			Pos:       pos,
			Symbol:    'z',
			Dexterity: 2,
			Strength:  5,
			Health:    22,
			MaxHealth: 22,
		},
		EnemyType: EnemyTypeZombie,
		Hostility: HostilityAverage,
		Alive:     true,
		Behavior:  &ZombieBehavior{},
		IsVisible: true,
	}
}

func NewVampire(pos entity.Point) *Enemy {
	return &Enemy{
		Character: Character{
			Pos:       pos,
			Symbol:    'v',
			Dexterity: 12,
			Strength:  6,
			Health:    18,
			MaxHealth: 18,
		},
		EnemyType: EnemyTypeVampire,
		Hostility: HostilityHigh,
		Alive:     true,
		Behavior:  &VampireBehavior{FirstHitPending: true},
		IsVisible: true,
	}
}

func NewGhost(pos entity.Point) *Enemy {
	return &Enemy{
		Character: Character{
			Pos:       pos,
			Symbol:    'g',
			Dexterity: 11,
			Strength:  3,
			Health:    8,
			MaxHealth: 8,
		},
		EnemyType: EnemyTypeGhost,
		Hostility: HostilityLow,
		Alive:     true,
		Behavior:  &GhostBehavior{},
		IsVisible: true,
	}
}

func NewOgre(pos entity.Point) *Enemy {
	return &Enemy{
		Character: Character{
			Pos:       pos,
			Symbol:    'O',
			Dexterity: 3,
			Strength:  10,
			Health:    35,
			MaxHealth: 35,
		},
		EnemyType: EnemyTypeOgre,
		Hostility: HostilityAverage,
		Alive:     true,
		Behavior:  &OgreBehavior{},
		IsVisible: true,
	}
}

func NewSnake(pos entity.Point) *Enemy {
	return &Enemy{
		Character: Character{
			Pos:       pos,
			Symbol:    's',
			Dexterity: 14,
			Strength:  5,
			Health:    12,
			MaxHealth: 12,
		},
		EnemyType: EnemyTypeSnake,
		Hostility: HostilityHigh,
		Alive:     true,
		Behavior:  &SnakeBehavior{},
		IsVisible: true,
	}
}
