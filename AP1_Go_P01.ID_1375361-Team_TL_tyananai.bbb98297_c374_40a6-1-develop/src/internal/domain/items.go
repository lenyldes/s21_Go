package domain

import (
	"math"
	"slices"

	"rogue/internal/domain/entity"
)

type Food struct {
	HP       int
	PosOnMap entity.Point
	Ch       rune
}

type Elixir struct {
	NumMoves int
	MaxHP    int
	Agility  int
	Power    int
	PosOnMap entity.Point
	Ch       rune
}

type Roll struct {
	MaxHP    int
	Agility  int
	Power    int
	PosOnMap entity.Point
	Ch       rune
}

type Weapon struct {
	Power    int
	PosOnMap entity.Point
	Ch       rune
}

func GetAllFoodPositionOnMap(level LevelMap) []entity.Point {
	positions := []entity.Point{}
	for _, f := range level.Food {
		positions = append(positions, f.PosOnMap)
	}
	return positions
}

func GenerateFood(level *LevelMap) Food {
	food := Food{}
	for {
		idRoom := randNum(0, 8)
		freeCells := GetFreeCellsInRoom(idRoom, level)
		if len(freeCells) == 0 {
			continue
		}
		randP := randNum(0, len(freeCells)-1)
		food.HP = randNum(15, 30)
		food.Ch = EnvFood
		food.PosOnMap = freeCells[randP]
		break
	}
	return food
}

func GenerateFoodOnMap(level *LevelMap) {
	count := int(math.Floor(max(1.0, 6.0-float64(level.Level)/4)))
	for i := 0; i < count; i++ {
		food := GenerateFood(level)
		level.Food = append(level.Food, food)
	}
}

func GetAllElixirPositionOnMap(level LevelMap) []entity.Point {
	positions := []entity.Point{}
	for _, e := range level.Elixir {
		positions = append(positions, e.PosOnMap)
	}
	return positions
}

func GenerateElixir(level *LevelMap) Elixir {
	elixir := Elixir{}
	for {
		idRoom := randNum(0, 8)
		freeCells := GetFreeCellsInRoom(idRoom, level)
		if len(freeCells) == 0 {
			continue
		}
		randP := randNum(0, len(freeCells)-1)
		elixir.NumMoves = randNum(10, 30)
		elixir.Ch = EnvElixir
		elixir.PosOnMap = freeCells[randP]
		switch randNum(0, 2) {
		case 0:
			elixir.MaxHP = randNum(10, 20)
		case 1:
			elixir.Agility = randNum(1, 4)
		case 2:
			elixir.Power = randNum(1, 4)
		}
		break
	}
	return elixir
}

func GenerateElixirOnMap(level *LevelMap) {
	count := int(math.Floor(max(1.0, 4.0-float64(level.Level)/6)))
	for i := 0; i < count; i++ {
		elixir := GenerateElixir(level)
		level.Elixir = append(level.Elixir, elixir)
	}
}

func GetAllRollPositionOnMap(level LevelMap) []entity.Point {
	positions := []entity.Point{}
	for _, r := range level.Roll {
		positions = append(positions, r.PosOnMap)
	}
	return positions
}

func GenerateRoll(level *LevelMap) Roll {
	roll := Roll{}
	for {
		idRoom := randNum(0, 8)
		freeCells := GetFreeCellsInRoom(idRoom, level)
		if len(freeCells) == 0 {
			continue
		}
		randP := randNum(0, len(freeCells)-1)
		roll.Ch = EnvRoll
		roll.PosOnMap = freeCells[randP]
		switch randNum(0, 2) {
		case 0:
			roll.MaxHP = randNum(10, 15)
		case 1:
			roll.Agility = randNum(1, 2)
		case 2:
			roll.Power = randNum(1, 2)
		}
		break
	}
	return roll
}

func GenerateRollOnMap(level *LevelMap) {
	count := int(math.Floor(max(0.0, 3.0-float64(level.Level)/5)))
	for i := 0; i < count; i++ {
		roll := GenerateRoll(level)
		level.Roll = append(level.Roll, roll)
	}
}

func GetAllWeaponPositionOnMap(level LevelMap) []entity.Point {
	positions := []entity.Point{}
	for _, w := range level.Weapon {
		positions = append(positions, w.PosOnMap)
	}
	return positions
}

func GenerateWeapon(level *LevelMap) Weapon {
	weapon := Weapon{}
	for {
		idRoom := randNum(0, 8)
		freeCells := GetFreeCellsInRoom(idRoom, level)
		if len(freeCells) == 0 {
			continue
		}
		randP := randNum(0, len(freeCells)-1)
		weapon.Ch = EnvWeapon
		weapon.PosOnMap = freeCells[randP]
		switch {
		case level.Level >= 17:
			weapon.Power = 20
		case level.Level >= 12:
			weapon.Power = 15
		case level.Level >= 7:
			weapon.Power = 10
		case level.Level >= 3:
			weapon.Power = 6
		default:
			weapon.Power = 3
		}
		break
	}
	return weapon
}

func GenerateWeaponOnMap(level *LevelMap) {
	count := randNum(0, 2)
	for i := 0; i < count; i++ {
		weapon := GenerateWeapon(level)
		level.Weapon = append(level.Weapon, weapon)
	}
}

func TryPickupItems(lm *LevelMap, player *Player) {
	pos := player.Pos

	for i, f := range lm.Food {
		if f.PosOnMap == pos {
			if len(player.Backpack.Food) < BackpackMaxItemsSize {
				player.Backpack.Food = append(player.Backpack.Food, f)
				lm.Food = slices.Delete(lm.Food, i, i+1)
				lm.Field[pos.Y][pos.X].EnvType = EnvRoom
			}
			return
		}
	}

	for i, e := range lm.Elixir {
		if e.PosOnMap == pos {
			if len(player.Backpack.Elixir) < BackpackMaxItemsSize {
				player.Backpack.Elixir = append(player.Backpack.Elixir, e)
				lm.Elixir = slices.Delete(lm.Elixir, i, i+1)
				lm.Field[pos.Y][pos.X].EnvType = EnvRoom
			}
			return
		}
	}

	for i, r := range lm.Roll {
		if r.PosOnMap == pos {
			if len(player.Backpack.Roll) < BackpackMaxItemsSize {
				player.Backpack.Roll = append(player.Backpack.Roll, r)
				lm.Roll = slices.Delete(lm.Roll, i, i+1)
				lm.Field[pos.Y][pos.X].EnvType = EnvRoom
			}
			return
		}
	}

	for i, w := range lm.Weapon {
		if w.PosOnMap == pos {
			if len(player.Backpack.Weapon) < BackpackMaxItemsSize {
				player.Backpack.Weapon = append(player.Backpack.Weapon, w)
				lm.Weapon = slices.Delete(lm.Weapon, i, i+1)
				lm.Field[pos.Y][pos.X].EnvType = EnvRoom
			}
			return
		}
	}
}
