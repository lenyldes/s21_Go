package domain

import (
	"slices"

	"rogue/internal/domain/entity"
)

type Coridor struct {
	Points []entity.Point
}

func GenerateHorizontalCorridor(a entity.Door, b entity.Door) Coridor {
	cor := Coridor{}
	inflectionPointX := randNum(a.PosOnMap.X+1, b.PosOnMap.X-1)

	for x := a.PosOnMap.X + 1; x <= inflectionPointX; x++ {
		cor.Points = append(cor.Points, entity.Point{X: x, Y: a.PosOnMap.Y})
	}
	for x := inflectionPointX; x < b.PosOnMap.X; x++ {
		cor.Points = append(cor.Points, entity.Point{X: x, Y: b.PosOnMap.Y})
	}
	yStart := min(a.PosOnMap.Y, b.PosOnMap.Y)
	yEnd := max(a.PosOnMap.Y, b.PosOnMap.Y)
	for y := yStart; y <= yEnd; y++ {
		cor.Points = append(cor.Points, entity.Point{X: inflectionPointX, Y: y})
	}

	return cor
}

func GenerateVerticalCorridor(a entity.Door, b entity.Door) Coridor {
	cor := Coridor{}
	inflectionPointY := randNum(a.PosOnMap.Y+1, b.PosOnMap.Y-1)

	for y := a.PosOnMap.Y + 1; y <= inflectionPointY; y++ {
		cor.Points = append(cor.Points, entity.Point{X: a.PosOnMap.X, Y: y})
	}
	for y := inflectionPointY; y < b.PosOnMap.Y; y++ {
		cor.Points = append(cor.Points, entity.Point{X: b.PosOnMap.X, Y: y})
	}
	xStart := min(a.PosOnMap.X, b.PosOnMap.X)
	xEnd := max(a.PosOnMap.X, b.PosOnMap.X)
	for x := xStart; x <= xEnd; x++ {
		cor.Points = append(cor.Points, entity.Point{X: x, Y: inflectionPointY})
	}

	return cor
}

func GenerateCoridors(level *LevelMap) {
	rooms := &level.Rooms
	coridors := &level.Coridors

	horizontalEdges := [][2]int{
		{0, 1}, {1, 2}, {3, 4}, {4, 5}, {6, 7}, {7, 8},
	}
	for _, edge := range horizontalEdges {
		a, b := edge[0], edge[1]
		roomA, roomB := (*rooms)[a], (*rooms)[b]
		if slices.Contains(roomA.Neighbor, roomB.ID) && slices.Contains(roomB.Neighbor, roomA.ID) {
			aDoor := entity.Door{}
			for _, door := range roomA.Doors {
				if door.PosOnMap.X == roomA.StartPos.X+roomA.Width-1 {
					aDoor = door
				}
			}
			bDoor := entity.Door{}
			for _, door := range roomB.Doors {
				if door.PosOnMap.X == roomB.StartPos.X {
					bDoor = door
				}
			}
			cor := GenerateHorizontalCorridor(aDoor, bDoor)
			*coridors = append(*coridors, cor)
		}
	}

	verticalEdges := [][2]int{
		{0, 3}, {3, 6}, {1, 4}, {4, 7}, {2, 5}, {5, 8},
	}
	for _, edge := range verticalEdges {
		a, b := edge[0], edge[1]
		roomA, roomB := (*rooms)[a], (*rooms)[b]
		if slices.Contains(roomA.Neighbor, roomB.ID) && slices.Contains(roomB.Neighbor, roomA.ID) {
			aDoor := entity.Door{}
			for _, door := range roomA.Doors {
				if door.PosOnMap.Y == roomA.StartPos.Y+roomA.Height-1 {
					aDoor = door
				}
			}
			bDoor := entity.Door{}
			for _, door := range roomB.Doors {
				if door.PosOnMap.Y == roomB.StartPos.Y {
					bDoor = door
				}
			}
			cor := GenerateVerticalCorridor(aDoor, bDoor)
			*coridors = append(*coridors, cor)
		}
	}
}
