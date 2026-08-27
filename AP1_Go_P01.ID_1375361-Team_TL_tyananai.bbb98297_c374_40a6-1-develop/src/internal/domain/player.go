package domain

import "rogue/internal/domain/entity"

type Player struct {
	Character
	Backpack        *Backpack
	CurrentWeapon   *Weapon
	ActiveElixir    *Elixir
	ElixirMovesLeft int
	IsPlayerSleep   bool
}

func NewPlayer(initPos entity.Point) *Player {
	return &Player{
		Character: Character{
			Pos:       initPos,
			Symbol:    EnvPerson,
			Health:    100,
			MaxHealth: 100,
			Strength:  8,
			Dexterity: 8,
		},
		Backpack: NewBackpack(),
	}
}

func GeneratePlayerStartPosition(rooms []entity.Room) (startPoint entity.Point, idStartRoom int) {
	idStartRoom = randNum(0, 8)
	startRoom := rooms[idStartRoom]
	startPoint = GenerateRandomPointInRoom(startRoom)
	return startPoint, idStartRoom
}
