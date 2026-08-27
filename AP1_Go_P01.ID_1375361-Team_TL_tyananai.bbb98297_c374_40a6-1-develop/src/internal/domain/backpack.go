package domain

import "slices"

const BackpackMaxItemsSize = 9

type Backpack struct {
	Gold   int
	Food   []Food
	Elixir []Elixir
	Roll   []Roll
	Weapon []Weapon
}

func NewBackpack() *Backpack {
	return &Backpack{}
}

func applyElixirStats(player *Player, e *Elixir) {
	player.MaxHealth += e.MaxHP
	player.Health += e.MaxHP
	player.Dexterity += e.Agility
	player.Strength += e.Power
}

func revertElixir(player *Player) {
	if player.ActiveElixir == nil {
		return
	}

	e := player.ActiveElixir

	player.MaxHealth -= e.MaxHP
	player.Health -= e.MaxHP
	if player.Health <= 0 {
		player.Health = 1 // подумать
	}

	player.Dexterity -= e.Agility
	player.Strength -= e.Power
	player.ActiveElixir = nil
	player.ElixirMovesLeft = 0
}

func tickElixir(player *Player) {
	if player.ActiveElixir == nil {
		return
	}

	player.ElixirMovesLeft--
	if player.ElixirMovesLeft <= 0 {
		revertElixir(player)
	}
}

func dropWeapon(player *Player, lm *LevelMap) {
	if player.CurrentWeapon == nil {
		return
	}

	w := *player.CurrentWeapon
	py, px := player.Pos.Y, player.Pos.X

	for x := px - 1; x <= px+1; x++ {
		for y := py - 1; y <= py+1; y++ {
			if x == px && y == py {
				continue
			}

			t := lm.Field[y][x].EnvType
			if t == EnvRoom || t == EnvCorridor {
				w.PosOnMap.Y = y
				w.PosOnMap.X = x
				lm.Weapon = append(lm.Weapon, w)
				lm.Field[y][x].EnvType = EnvWeapon
				player.CurrentWeapon = nil
				return
			}
		}
	}

	player.CurrentWeapon = nil // если свободных клеток не оказалось??
}

func UseFood(player *Player, idx int) bool {
	if idx < 1 || idx > len(player.Backpack.Food) {
		return false
	}

	food := player.Backpack.Food[idx-1]

	player.Health += food.HP
	if player.Health > player.MaxHealth {
		player.Health = player.MaxHealth
	}
	player.Backpack.Food = slices.Delete(player.Backpack.Food, idx-1, idx)
	return true
}

func UseElixir(player *Player, idx int) bool {
	if idx < 1 || idx > len(player.Backpack.Elixir) {
		return false
	}

	if player.ActiveElixir != nil {
		revertElixir(player)
	}

	elixir := player.Backpack.Elixir[idx-1]

	player.Backpack.Elixir = slices.Delete(player.Backpack.Elixir, idx-1, idx)
	elixirCopy := elixir
	player.ActiveElixir = &elixirCopy
	player.ElixirMovesLeft = elixir.NumMoves
	applyElixirStats(player, player.ActiveElixir)
	return true
}

func UseRoll(player *Player, idx int) bool {
	if idx < 1 || idx > len(player.Backpack.Roll) {
		return false
	}

	roll := player.Backpack.Roll[idx-1]

	player.MaxHealth += roll.MaxHP
	player.Health += roll.MaxHP
	player.Dexterity += roll.Agility
	player.Strength += roll.Power
	player.Backpack.Roll = slices.Delete(player.Backpack.Roll, idx-1, idx)
	return true
}

func EquipWeapon(player *Player, lm *LevelMap, idx int) bool {
	if idx == 0 {
		dropWeapon(player, lm)
		return true
	}
	if idx < 1 || idx > len(player.Backpack.Weapon) {
		return false
	}
	if player.CurrentWeapon != nil {
		dropWeapon(player, lm)
	}

	w := player.Backpack.Weapon[idx-1]

	player.Backpack.Weapon = slices.Delete(player.Backpack.Weapon, idx-1, idx)
	wCopy := w
	player.CurrentWeapon = &wCopy
	return true
}
