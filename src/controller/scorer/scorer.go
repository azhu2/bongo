package scorer

import (
	"context"
	"fmt"
	"math"

	"github.com/azhu2/bongo/src/entity"
	"go.uber.org/fx"
)

const (
	commonMultiplier = 1.3
)

var Module = fx.Module("scorer",
	fx.Provide(New),
)

type Controller interface {
	Score(context.Context, entity.Board, entity.Solution) (int, error)
}

type Results struct {
	fx.Out

	Controller
}

type scorer struct{}

func New() (Results, error) {
	return Results{
		Controller: &scorer{},
	}, nil
}

func (s *scorer) Score(ctx context.Context, board entity.Board, solution entity.Solution) (int, error) {
	score := 0

	for rowIdx, row := range solution.Board {
		wordScore := 0
		for colIdx, letter := range row {
			letterScore, err := scoreLetter(ctx, board, rowIdx, colIdx, letter)
			if err != nil {
				return 0, err
			}
			wordScore += letterScore
		}
		score += (int)(math.Ceil(s.wordMultiplier(ctx, string(row)) * float64(wordScore)))
	}

	bonusLetters := make([]rune, len(board.BonusWord))
	bonusScore := 0
	for _, coords := range board.BonusWord {
		rowIdx := coords[0]
		colIdx := coords[1]
		letter := solution.Board[rowIdx][colIdx]
		bonusLetters = append(bonusLetters, letter)
		letterScore, err := scoreLetter(ctx, board, rowIdx, colIdx, letter)
		if err != nil {
			return 0, err
		}
		bonusScore += letterScore
	}
	score += (int)(math.Ceil(s.wordMultiplier(ctx, string(bonusLetters)) * float64(bonusScore)))

	return score, nil
}

func scoreLetter(_ context.Context, board entity.Board, row, col int, letter rune) (int, error) {
	tile, ok := board.Tiles[letter]
	if !ok {
		return 0, fmt.Errorf("solution has invalid letter: %c", letter)
	}
	return board.Multipliers[row][col] * tile.Value, nil
}

func (s *scorer) wordMultiplier(ctx context.Context, word string) float64 {
	if !s.isWord(ctx, word) {
		return 0
	}
	if s.isCommon(ctx, word) {
		return commonMultiplier
	}
	return 1
}

func (s *scorer) isWord(_ context.Context, _ string) bool {
	// TODO Implement
	return true
}

func (s *scorer) isCommon(_ context.Context, _ string) bool {
	// TODO Implement
	return true
}
