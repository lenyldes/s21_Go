package domain

import (
	"math/rand"

	"rogue/internal/domain/entity"
)

const (
	WIDTH         = 80
	HEIGHT        = 26
	ROOM_OFFSET_X = (WIDTH - 2) / 3
	ROOM_OFFSET_Y = (HEIGHT - 2) / 3
	ROOMS_COUNT   = 9
	EnvWall       = '#'
	EnvRoom       = '.'
	EnvDoor       = '+'
	EnvEmpty      = ' '
	EnvCorridor   = '='
	EnvExit       = '%'
	EnvPerson     = '@'
	EnvKey        = '*'
	EnvFood       = 'F'
	EnvElixir     = 'E'
	EnvRoll       = 'R'
	EnvWeapon     = 'W'
)

const (
	RoomMinW = 4
	RoomMinH = 4
	RoomMaxW = 26
	RoomMaxH = 8
)

func randNum(a int, b int) int {
	return rand.Intn(b-a+1) + a
}

type LevelMap struct {
	Level   int
	Field   [HEIGHT][WIDTH]entity.Cell
	Enemies []*Enemy
	Spawn   entity.Point
	Exit    entity.Point

	Rooms       []entity.Room
	Coridors    []Coridor
	IdStartRoom int
	IdEndRoom   int
	Key         Key
	Food        []Food
	Elixir      []Elixir
	Roll        []Roll
	Weapon      []Weapon
}

func NewLevelMap(level int) *LevelMap {
	lm := &LevelMap{Level: level}

	GenerateRooms(lm)
	GenerateNeighbor(&lm.Rooms)
	GenerateDoors(&lm.Rooms)
	GenerateCoridors(lm)

	startPoint, idStartRoom := GeneratePlayerStartPosition(lm.Rooms)
	lm.Spawn = startPoint
	lm.IdStartRoom = idStartRoom
	CreateNewKey(lm)

	GenerateFoodOnMap(lm)
	GenerateElixirOnMap(lm)
	GenerateRollOnMap(lm)
	GenerateWeaponOnMap(lm)

	initFuild(lm)
	UpdateFOW(lm, lm.Spawn)

	return lm
}

func initEmptyCell(level *LevelMap) {
	for y := range level.Field {
		for x := range level.Field[y] {
			level.Field[y][x].EnvType = EnvEmpty
		}
	}
}

func initFloor(level *LevelMap) {
	for _, room := range level.Rooms {
		for x := room.StartPos.X; x < room.StartPos.X+room.Width; x++ {
			for y := room.StartPos.Y; y < room.StartPos.Y+room.Height; y++ {
				if y == room.StartPos.Y || y == room.StartPos.Y+room.Height-1 ||
					x == room.StartPos.X || x == room.StartPos.X+room.Width-1 {
					level.Field[y][x].EnvType = EnvWall
				} else {
					level.Field[y][x].EnvType = EnvRoom
				}
			}
		}
	}
}

func initDoors(level *LevelMap) {
	for _, room := range level.Rooms {
		for _, door := range room.Doors {
			level.Field[door.PosOnMap.Y][door.PosOnMap.X].EnvType = EnvDoor
		}
	}
}

func initCorridors(level *LevelMap) {
	for _, cor := range level.Coridors {
		for _, p := range cor.Points {
			level.Field[p.Y][p.X].EnvType = EnvCorridor
		}
	}
}

func initFood(level *LevelMap) {
	for _, p := range GetAllFoodPositionOnMap(*level) {
		level.Field[p.Y][p.X].EnvType = EnvFood
	}
}

func initElexir(level *LevelMap) {
	for _, p := range GetAllElixirPositionOnMap(*level) {
		level.Field[p.Y][p.X].EnvType = EnvElixir
	}
}

func initRoll(level *LevelMap) {
	for _, p := range GetAllRollPositionOnMap(*level) {
		level.Field[p.Y][p.X].EnvType = EnvRoll
	}
}

func initWeapon(level *LevelMap) {
	for _, p := range GetAllWeaponPositionOnMap(*level) {
		level.Field[p.Y][p.X].EnvType = EnvWeapon
	}
}

func initFuild(level *LevelMap) {
	initEmptyCell(level)
	initFloor(level)
	initDoors(level)
	initCorridors(level)
	initFood(level)
	initElexir(level)
	initRoll(level)
	initWeapon(level)
}
