package entity

import (
	"maps"
	"slices"
)

const (
	BoardSize    = 5
	MaxWildcards = 1
)

type Board struct {
	Tiles       map[rune]Tile
	Multipliers [][]int // Grid of multipliers
	BonusWord   [][]int // Slice of [row,col] coords

	sortedTiles []rune
}

type Tile struct {
	Value int
	Count int
}

func Less(a, b Tile) bool {
	return a.Value < b.Value
}

func (b Board) SortedTiles() []rune {
	if b.sortedTiles != nil {
		return b.sortedTiles
	}

	tiles := make([]rune, 0, len(b.Tiles))
	for letter := range maps.Keys(b.Tiles) {
		tiles = append(tiles, letter)
	}
	slices.SortFunc(tiles, func(x, y rune) int {
		return b.Tiles[y].Value - b.Tiles[x].Value
	})

	b.sortedTiles = tiles
	return tiles
}
