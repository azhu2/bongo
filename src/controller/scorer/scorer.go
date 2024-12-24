package scorer

import (
	"context"
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
	Score(context.Context, entity.Board, entity.Solution) int
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

func (s *scorer) Score(ctx context.Context, board entity.Board, solution entity.Solution) int {
	score := 0

	for _, row := range solution.Board {
		score += s.scoreWord(ctx, row)
	}

	bonusTiles := make([]entity.Tile, len(board.BonusWord))
	for _, coords := range board.BonusWord {
		bonusTiles = append(bonusTiles, solution.Board[coords[0]][coords[1]])
	}
	score += s.scoreWord(ctx, bonusTiles)

	return score
}

func (s *scorer) scoreWord(ctx context.Context, tiles []entity.Tile) int {
	letters := make([]rune, len(tiles))
	value := 0
	for _, tile := range tiles {
		letters = append(letters, tile.Letter)
		value += tile.Value
	}
	word := string(letters)
	if !s.isWord(ctx, word) {
		return 0
	}
	if s.isCommon(ctx, word) {
		value = (int)(math.Ceil(commonMultiplier * float64(value)))
	}
	return value
}

func (s *scorer) isWord(_ context.Context, _ string) bool {
	// TODO Implement
	return true
}

func (s *scorer) isCommon(_ context.Context, _ string) bool {
	// TODO Implement
	return true
}
