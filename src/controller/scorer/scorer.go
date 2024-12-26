package scorer

import (
	"context"
	"fmt"
	"math"
	"strings"

	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/entity"
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

type Result struct {
	fx.Out

	Controller
}

type scorer struct{}

func New() (Result, error) {
	return Result{
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
	if letter == ' ' {
		return 0, nil
	}
	tile, ok := board.Tiles[letter]
	if !ok {
		return 0, fmt.Errorf("solution has invalid letter: %c", letter)
	}
	return board.Multipliers[row][col] * tile.Value, nil
}

func (s *scorer) wordMultiplier(ctx context.Context, word string) float64 {
	// Assume no gaps. I'm pretty sure only one word per line is counted.
	trimmed := strings.TrimSpace(word)
	if !s.isWord(ctx, trimmed) {
		return 0
	}
	if s.isCommon(ctx, trimmed) {
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
