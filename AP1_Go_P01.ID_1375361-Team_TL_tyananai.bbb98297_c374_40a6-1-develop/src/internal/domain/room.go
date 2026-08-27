package domain

import (
	"math/rand"
	"slices"

	"rogue/internal/domain/entity"
)

func CreateNewRoom(start entity.Point, w int, h int, id int) entity.Room {
	var room entity.Room
	room.StartPos = start
	room.Width = w
	room.Height = h
	room.ID = id
	return room
}

func GenerateRooms(level *LevelMap) {
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			xFrom := x * (ROOM_OFFSET_X + 1)
			xTo := xFrom + (ROOM_OFFSET_X - RoomMinW)

			yFrom := y * (ROOM_OFFSET_Y + 1)
			yTo := yFrom + (ROOM_OFFSET_Y - RoomMinH)

			xRand := randNum(xFrom, xTo)
			yRand := randNum(yFrom, yTo)

			wMax := ROOM_OFFSET_X*(x+1) + x - xRand
			hMax := ROOM_OFFSET_Y*(y+1) + y - yRand

			width := randNum(RoomMinW, wMax)
			height := randNum(RoomMinH, hMax)

			id := y*3 + x
			room := CreateNewRoom(entity.Point{X: xRand, Y: yRand}, width, height, id)
			level.Rooms = append(level.Rooms, room)
		}
	}
}

func findGroupRepresentative(parent []int, x int) int {
	if parent[x] != x {
		parent[x] = findGroupRepresentative(parent, parent[x])
	}
	return parent[x]
}

func unionGroup(parent []int, a int, b int) bool {
	ra, rb := findGroupRepresentative(parent, a), findGroupRepresentative(parent, b)
	if ra == rb {
		return false
	}
	parent[ra] = rb
	return true
}

func GenerateNeighbor(rooms *[]entity.Room) {
	allEdges := [][2]int{
		{0, 1}, {1, 2}, {3, 4}, {4, 5}, {6, 7}, {7, 8},
		{0, 3}, {3, 6}, {1, 4}, {4, 7}, {2, 5}, {5, 8},
	}

	rand.Shuffle(len(allEdges), func(i, j int) {
		allEdges[i], allEdges[j] = allEdges[j], allEdges[i]
	})

	parent := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

	for _, edge := range allEdges {
		a, b := edge[0], edge[1]
		if unionGroup(parent, a, b) {
			(*rooms)[a].Neighbor = append((*rooms)[a].Neighbor, b)
			(*rooms)[b].Neighbor = append((*rooms)[b].Neighbor, a)
		}
	}
}

func GenerateDoor(room entity.Room, doorType rune) entity.Room {
	door := entity.Door{}
	switch doorType {
	case 'u':
		door.PosOnMap.X = randNum(room.StartPos.X+1, room.StartPos.X+room.Width-2)
		door.PosOnMap.Y = room.StartPos.Y
	case 'd':
		door.PosOnMap.X = randNum(room.StartPos.X+1, room.StartPos.X+room.Width-2)
		door.PosOnMap.Y = room.StartPos.Y + room.Height - 1
	case 'l':
		door.PosOnMap.X = room.StartPos.X
		door.PosOnMap.Y = randNum(room.StartPos.Y+1, room.StartPos.Y+room.Height-2)
	case 'r':
		door.PosOnMap.X = room.StartPos.X + room.Width - 1
		door.PosOnMap.Y = randNum(room.StartPos.Y+1, room.StartPos.Y+room.Height-2)
	}
	room.Doors = append(room.Doors, door)
	return room
}

func GenerateDoors(rooms *[]entity.Room) {
	for i, room := range *rooms {
		for _, neighborID := range room.Neighbor {
			diff := neighborID - room.ID
			switch diff {
			case 1:
				(*rooms)[i] = GenerateDoor((*rooms)[i], 'r')
			case -1:
				(*rooms)[i] = GenerateDoor((*rooms)[i], 'l')
			case 3:
				(*rooms)[i] = GenerateDoor((*rooms)[i], 'd')
			case -3:
				(*rooms)[i] = GenerateDoor((*rooms)[i], 'u')
			}
		}
	}
}

func GenerateRandomPointInRoom(room entity.Room) entity.Point {
	return entity.Point{
		X: randNum(room.StartPos.X+1, room.StartPos.X+room.Width-2),
		Y: randNum(room.StartPos.Y+1, room.StartPos.Y+room.Height-2),
	}
}

func GetFreeCellsInRoom(id int, level *LevelMap) []entity.Point {
	room := level.Rooms[id]
	busyCells := []entity.Point{}

	busyCells = append(busyCells, level.Spawn)
	busyCells = append(busyCells, level.Key.PosOnMap)

	for _, f := range level.Food {
		busyCells = append(busyCells, f.PosOnMap)
	}

	freeCells := []entity.Point{}
	for x := room.StartPos.X + 1; x < room.StartPos.X+room.Width-1; x++ {
		for y := room.StartPos.Y + 1; y < room.StartPos.Y+room.Height-1; y++ {
			p := entity.Point{X: x, Y: y}
			if !slices.Contains(busyCells, p) {
				freeCells = append(freeCells, p)
			}
		}
	}

	return freeCells
}

func GetRoomCoordinatesWithP(level *LevelMap, posP entity.Point) (leftX int, rightX int, topY int, bottomY int) {
	leftX = posP.X
	for leftX > 0 && level.Field[posP.Y][leftX].EnvType != EnvWall && level.Field[posP.Y][leftX].EnvType != EnvDoor {
		leftX--
	}

	rightX = posP.X
	for rightX < WIDTH-1 && level.Field[posP.Y][rightX].EnvType != EnvWall && level.Field[posP.Y][rightX].EnvType != EnvDoor {
		rightX++
	}

	topY = posP.Y
	for topY > 0 && level.Field[topY][posP.X].EnvType != EnvWall && level.Field[topY][posP.X].EnvType != EnvDoor {
		topY--
	}

	bottomY = posP.Y
	for bottomY < HEIGHT-1 && level.Field[bottomY][posP.X].EnvType != EnvWall && level.Field[bottomY][posP.X].EnvType != EnvDoor {
		bottomY++
	}

	return leftX, rightX, topY, bottomY
}
