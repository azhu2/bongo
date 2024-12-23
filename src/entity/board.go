package entity

type Board struct {
	Tiles       []Tile
	Multipliers [][]int // Grid of multipliers - (0, 0) is bottom-left
	BonusWord   [][]int // Slice of [x,y] coords
}

type Tile struct {
	Letter rune
	Value  int
}
