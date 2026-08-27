package entity

type Door struct {
	PosOnMap Point
}

type Room struct {
	StartPos Point
	Width    int
	Height   int
	ID       int
	Neighbor []int
	Doors    []Door
}
