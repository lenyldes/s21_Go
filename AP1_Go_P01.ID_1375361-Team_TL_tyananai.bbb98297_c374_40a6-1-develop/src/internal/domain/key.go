package domain

import "rogue/internal/domain/entity"

type Key struct {
	PosOnMap entity.Point
	Ch       rune
}

func CreateNewKey(level *LevelMap) {
	endPoint, idEndRoom := GenerateKeyStartPosition(level.IdStartRoom, level.Rooms)

	level.Key.PosOnMap.X = endPoint.X
	level.Key.PosOnMap.Y = endPoint.Y
	level.Key.Ch = EnvRoom

	level.IdEndRoom = idEndRoom
}

func GenerateKeyStartPosition(idStartRoom int, rooms []entity.Room) (endPoint entity.Point, idEndRoom int) {
	for {
		idEndRoom = randNum(0, 8)
		if idStartRoom != idEndRoom {
			break
		}
	}

	endRoom := rooms[idEndRoom]
	endPoint = GenerateRandomPointInRoom(endRoom)

	return endPoint, idEndRoom
}

func CheckFindKey(level LevelMap, player *Player) bool {
	return player.Pos == level.Key.PosOnMap
}
