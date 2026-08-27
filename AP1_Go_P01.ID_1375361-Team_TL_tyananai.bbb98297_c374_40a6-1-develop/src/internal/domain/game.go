package domain

import (
	"rogue/internal/domain/entity"
)

type Game struct {
	CurrentLevelMap *LevelMap `json:"currentLevelMap"`
	Player          *Player   `json:"player"`
	LastLog         string    `json:"lastLog"`
	PlayerDead      bool      `json:"playerDead"`
}

func NewGame() *Game {
	lm := NewLevelMap(1)
	g := &Game{
		CurrentLevelMap: lm,
		Player:          NewPlayer(lm.Spawn),
	}
	g.spawnEnemies()
	return g
}

func (g *Game) Log(msg string) {
	if g.LastLog != "" {
		g.LastLog += " " + msg
	} else {
		g.LastLog = msg
	}
}

func (g *Game) IsCellMovable(y int, x int) bool {
	if y < 0 || y >= HEIGHT || x < 0 || x >= WIDTH {
		return false
	}
	field := g.CurrentLevelMap.Field
	t := field[y][x].EnvType
	return t != EnvWall && t != EnvEmpty
}

func (g *Game) spawnEnemies() {
	lm := g.CurrentLevelMap
	sp := g.pickSpreadSpawnPoints(5)
	lm.Enemies = []*Enemy{
		NewZombie(sp[0]),
		NewVampire(sp[1]),
		NewGhost(sp[2]),
		NewOgre(sp[3]),
		NewSnake(sp[4]),
	}
}

func (g *Game) Update(event string) {
	if g.PlayerDead {
		return
	}
	if event != "" {
		g.LastLog = ""
	}

	lm := g.CurrentLevelMap
	player := g.Player

	if player.IsPlayerSleep {
		player.IsPlayerSleep = false
		g.Log("You slept through your turn!")
		g.StepAllEnemies()
		UpdateFOW(lm, player.Pos)
		return
	}

	px, py := player.Pos.X, player.Pos.Y
	turnTaken := false

	tryMove := func(ny, nx int) {
		if !g.IsCellMovable(ny, nx) {
			return
		}
		target := entity.Point{Y: ny, X: nx}
		if e := g.enemyAtPoint(target); e != nil {
			g.MeleeExchange(e)
			turnTaken = true
			return
		}
		player.Pos = target
		turnTaken = true
	}

	switch event {
	case "w", "W":
		tryMove(py-1, px)
	case "s", "S":
		tryMove(py+1, px)
	case "a", "A":
		tryMove(py, px-1)
	case "d", "D":
		tryMove(py, px+1)
	}

	if turnTaken {
		tickElixir(player)
		TryPickupItems(lm, player)

		if player.Health <= 0 {
			player.Health = 0
			g.PlayerDead = true
			g.Log("YOU DIED!")
			return
		}

		if CheckFindKey(*lm, player) {
			*lm = *NewLevelMap(lm.Level + 1)
			player.Pos = lm.Spawn
			g.spawnEnemies()
		} else {
			g.StepAllEnemies()
		}

		UpdateFOW(lm, player.Pos)
	}
}
