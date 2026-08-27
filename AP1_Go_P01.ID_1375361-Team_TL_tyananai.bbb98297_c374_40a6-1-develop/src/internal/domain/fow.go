package domain

import "rogue/internal/domain/entity"

func UpdateFOW(level *LevelMap, pos entity.Point) {
	ch := level.Field[pos.Y][pos.X].EnvType

	if ch == EnvRoom {
		UpdateFOWRoom(level, pos)
	}
	if ch == EnvCorridor || ch == EnvDoor {
		UpdateFOWCorridorAndDoor(level, pos)
		UpdateFOWHideItems(level)
	}
}

func UpdateFOWRoom(level *LevelMap, pos entity.Point) {
	leftX, rightX, topY, bottomY := GetRoomCoordinatesWithP(level, pos)

	for y := topY; y <= bottomY; y++ {
		for x := leftX; x <= rightX; x++ {
			level.Field[y][x].Visibility = true

			kP := level.Key.PosOnMap
			if kP.X == x && kP.Y == y {
				level.Key.Ch = EnvKey
			}
		}
	}
}

func UpdateFOWCorridorAndDoor(level *LevelMap, pos entity.Point) {
	for y := pos.Y - 1; y <= pos.Y+1; y++ {
		for x := pos.X - 1; x <= pos.X+1; x++ {
			if y >= 0 && y < HEIGHT && x >= 0 && x < WIDTH {
				if level.Field[y][x].EnvType == EnvCorridor || level.Field[y][x].EnvType == EnvDoor {
					level.Field[y][x].Visibility = true
				}
			}
		}
	}
}

func UpdateFOWHideItems(level *LevelMap) {
	level.Key.Ch = EnvRoom
}
