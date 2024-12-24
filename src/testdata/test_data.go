package testdata

import (
	"github.com/azhu2/bongo/src/entity"
)

type testCase struct {
	Name     string
	Filename string
	Board    entity.Board
	Solution entity.Solution
	Score    int
}

var TestData = []testCase{
	{
		Name:     "2024-12-23",
		Filename: "2024-12-23.txt",
		Board: entity.Board{
			Tiles: map[rune]entity.Tile{
				'W': entity.Tile{Value: 65, Count: 1},
				'H': entity.Tile{Value: 40, Count: 1},
				'P': entity.Tile{Value: 45, Count: 2},
				'M': entity.Tile{Value: 35, Count: 1},
				'Y': entity.Tile{Value: 35, Count: 1},
				'D': entity.Tile{Value: 35, Count: 1},
				'N': entity.Tile{Value: 20, Count: 2},
				'L': entity.Tile{Value: 10, Count: 1},
				'O': entity.Tile{Value: 7, Count: 1},
				'R': entity.Tile{Value: 7, Count: 2},
				'A': entity.Tile{Value: 5, Count: 2},
				'E': entity.Tile{Value: 5, Count: 6},
				'S': entity.Tile{Value: 5, Count: 4},
			},
			Multipliers: [][]int{
				[]int{1, 3, 1, 1, 1},
				[]int{1, 2, 1, 1, 1},
				[]int{2, 1, 1, 1, 1},
				[]int{1, 1, 1, 1, 1},
				[]int{1, 1, 1, 1, 1},
			},
			BonusWord: [][]int{
				[]int{0, 1},
				[]int{1, 1},
				[]int{2, 2},
				[]int{3, 3},
			},
		},
		Solution: entity.Solution{
			Board: [][]rune{
				[]rune("SWORN"),
				[]rune("SHAME"),
				[]rune("PLANE"),
				[]rune("SEEPS"),
				[]rune("REEDY"),
			},
		},
		Score: 1265,
	},
	{
		Name:     "2024-12-24",
		Filename: "2024-12-24.txt",
		Board: entity.Board{
			Tiles: map[rune]entity.Tile{
				'G': entity.Tile{Value: 45, Count: 2},
				'P': entity.Tile{Value: 35, Count: 1},
				'C': entity.Tile{Value: 35, Count: 3},
				'H': entity.Tile{Value: 50, Count: 1},
				'Y': entity.Tile{Value: 35, Count: 1},
				'N': entity.Tile{Value: 20, Count: 2},
				'T': entity.Tile{Value: 9, Count: 2},
				'I': entity.Tile{Value: 10, Count: 2},
				'L': entity.Tile{Value: 10, Count: 1},
				'O': entity.Tile{Value: 7, Count: 2},
				'R': entity.Tile{Value: 7, Count: 1},
				'A': entity.Tile{Value: 5, Count: 6},
				'S': entity.Tile{Value: 5, Count: 1},
			},
			Multipliers: [][]int{
				[]int{1, 1, 1, 1, 1},
				[]int{1, 1, 1, 1, 1},
				[]int{1, 1, 3, 1, 1},
				[]int{1, 1, 1, 1, 2},
				[]int{1, 1, 1, 1, 2},
			},
			BonusWord: [][]int{
				[]int{0, 1},
				[]int{1, 1},
				[]int{2, 1},
				[]int{3, 1},
			},
		},
		Solution: entity.Solution{
			Board: [][]rune{
				[]rune("OCTAL"),
				[]rune("CHAAT"),
				[]rune("CIGAR"),
				[]rune(" PONY"),
				[]rune(" SING"),
			},
		},
		Score: 976,
	},
}
