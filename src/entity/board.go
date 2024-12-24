package entity

const (
	BoardSize = 5
)

type Board struct {
	Tiles       map[rune]Tile
	Multipliers [][]int // Grid of multipliers - (0, 0) is bottom-left
	BonusWord   [][]int // Slice of [x,y] coords
}

type Tile struct {
	Value int
	Count int
}
