package entity

const (
	BoardSize    = 5
	MaxWildcards = 1
)

type Board struct {
	Tiles       map[rune]Tile
	Multipliers [][]int // Grid of multipliers
	BonusWord   [][]int // Slice of [row,col] coords
}

type Tile struct {
	Value int
	Count int
}
