package domain

import (
	"math/rand/v2"

	"rogue/internal/domain/entity"
)

func (g *Game) bfsNextToward(self *Enemy, goal entity.Point) (entity.Point, bool) {
	from := self.Pos
	if from == goal {
		return entity.Point{}, false
	}

	queue := []entity.Point{from}
	parent := make(map[entity.Point]entity.Point)
	parent[from] = from

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, d := range cardinalDirs {
			nxt := entity.Point{Y: cur.Y + d.Y, X: cur.X + d.X}
			if nxt == goal {
				p := cur
				for parent[p] != from {
					p = parent[p]
				}
				return p, true
			}
			if !g.walkableForEnemy(nxt, self) {
				continue
			}
			if _, seen := parent[nxt]; seen {
				continue
			}
			parent[nxt] = cur
			queue = append(queue, nxt)
		}
	}
	return entity.Point{}, false
}

func (g *Game) pickSpreadSpawnPoints(n int) []entity.Point {
	var out []entity.Point
	rooms := g.CurrentLevelMap.Rooms

	playerRoomID := -1
	if pRoom, found := g.getRoomAt(g.Player.Pos); found {
		playerRoomID = pRoom.ID
	}

	var validRooms []entity.Room
	for _, r := range rooms {
		if r.ID == playerRoomID {
			continue
		}
		validRooms = append(validRooms, r)
	}

	if len(validRooms) == 0 {
		validRooms = rooms
	}

	for i := 0; i < n; i++ {
		r := validRooms[rand.IntN(len(validRooms))]
		heightInner := r.Height - 2
		widthInner := r.Width - 2
		if heightInner < 1 {
			heightInner = 1
		}
		if widthInner < 1 {
			widthInner = 1
		}
		ry := r.StartPos.Y + 1 + rand.IntN(heightInner)
		rx := r.StartPos.X + 1 + rand.IntN(widthInner)
		out = append(out, entity.Point{Y: ry, X: rx})
	}
	return out
}
