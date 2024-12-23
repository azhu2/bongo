package entity

type Board struct {
	Tiles       []Tile
	Multipliers [][]int
	BonusWord   [][]int // Slice of [x,y] coords
}

type Tile struct {
	Letter rune
	Value  int
}
