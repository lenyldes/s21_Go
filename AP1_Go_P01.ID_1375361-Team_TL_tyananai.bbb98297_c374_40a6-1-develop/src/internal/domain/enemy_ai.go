package domain

import (
	"math/rand/v2"

	"rogue/internal/domain/entity"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattan(a, b entity.Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func (g *Game) HasLineOfSight(p1, p2 entity.Point) bool {
	x0, y0 := p1.X, p1.Y
	x1, y1 := p2.X, p2.Y
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy

	for {
		if x0 == x1 && y0 == y1 {
			return true
		}
		if (x0 != p1.X || y0 != p1.Y) && (x0 != p2.X || y0 != p2.Y) {
			if !g.IsCellMovable(y0, x0) {
				return false
			}
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

var cardinalDirs = []entity.Point{
	{X: -1, Y: 0},
	{X: 1, Y: 0},
	{X: 0, Y: -1},
	{X: 0, Y: 1},
}

func (g *Game) walkableForEnemy(p entity.Point, self *Enemy) bool {
	if !g.IsCellMovable(p.Y, p.X) {
		return false
	}
	for i := range g.CurrentLevelMap.Enemies {
		e := g.CurrentLevelMap.Enemies[i]
		if e.Alive && e != self && e.Pos == p {
			return false
		}
	}
	return p != g.Player.Pos
}

// getRoomAt finds the room that contains point p, using StartPos/Width/Height.
func (g *Game) getRoomAt(p entity.Point) (entity.Room, bool) {
	for _, r := range g.CurrentLevelMap.Rooms {
		brY := r.StartPos.Y + r.Height - 1
		brX := r.StartPos.X + r.Width - 1
		if p.Y >= r.StartPos.Y && p.Y <= brY && p.X >= r.StartPos.X && p.X <= brX {
			return r, true
		}
	}
	return entity.Room{}, false
}

func (g *Game) chasePlayer(e *Enemy, steps int, restrictToRoom bool) bool {
	playerPos := g.Player.Pos

	if restrictToRoom {
		eRoom, inERoom := g.getRoomAt(e.Pos)
		pRoom, inPRoom := g.getRoomAt(playerPos)
		if !inERoom || !inPRoom || eRoom.ID != pRoom.ID {
			return false
		}
	}

	moved := false
	room, inRoom := g.getRoomAt(e.Pos)

	for s := 0; s < steps; s++ {
		next, ok := g.bfsNextToward(e, playerPos)
		if !ok {
			break
		}

		if restrictToRoom && inRoom {
			brY := room.StartPos.Y + room.Height - 1
			brX := room.StartPos.X + room.Width - 1
			if next.Y < room.StartPos.Y || next.Y > brY || next.X < room.StartPos.X || next.X > brX {
				break
			}
		}

		e.Pos = next
		moved = true
		if manhattan(e.Pos, playerPos) <= 1 {
			break
		}
	}
	return moved
}

func (g *Game) defaultMove(e *Enemy, steps int, restrictToRoom bool) {
	playerPos := g.Player.Pos
	dist := manhattan(e.Pos, playerPos)

	if !e.IsAwake {
		if dist <= e.AggroRadius() && g.HasLineOfSight(e.Pos, playerPos) {
			e.IsAwake = true
		}
	}

	if e.IsAwake {
		if g.chasePlayer(e, steps, restrictToRoom) {
			return
		}
		if dist > e.AggroRadius()*2 {
			e.IsAwake = false
		}
	}

	for s := 0; s < steps; s++ {
		dirs := []entity.Point{
			{Y: -1, X: 0},
			{Y: 1, X: 0},
			{Y: 0, X: -1},
			{Y: 0, X: 1},
		}
		rand.Shuffle(len(dirs), func(i, j int) {
			dirs[i], dirs[j] = dirs[j], dirs[i]
		})

		for _, d := range dirs {
			np := entity.Point{Y: e.Pos.Y + d.Y, X: e.Pos.X + d.X}
			if g.walkableForEnemy(np, e) {
				if restrictToRoom {
					room, inRoom := g.getRoomAt(e.Pos)
					if inRoom {
						brY := room.StartPos.Y + room.Height - 1
						brX := room.StartPos.X + room.Width - 1
						if np.Y < room.StartPos.Y || np.Y > brY || np.X < room.StartPos.X || np.X > brX {
							continue
						}
					}
				}
				e.Pos = np
				break
			}
		}
	}
}

func (g *Game) stepEnemy(e *Enemy) {
	if !e.Alive || g.PlayerDead {
		return
	}

	playerPos := g.Player.Pos
	if manhattan(e.Pos, playerPos) == 1 {
		g.resolveEnemyHitOnPlayer(e)
		return
	}

	e.Behavior.TakeTurn(g, e)
}

func (g *Game) StepAllEnemies() {
	for i := range g.CurrentLevelMap.Enemies {
		g.stepEnemy(g.CurrentLevelMap.Enemies[i])
	}
}
