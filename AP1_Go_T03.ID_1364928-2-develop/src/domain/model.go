package domain

const (
	Empty    = 0
	Computer = 1
	Human    = 2
)

type Board [3][3]int

type Game struct {
	UUID  string
	Board Board
}
