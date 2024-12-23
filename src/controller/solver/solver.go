package solver

import (
	"context"

	"github.com/azhu2/bongo/src/entity"
	"go.uber.org/fx"
)

var Module = fx.Module("solver",
	fx.Provide(New),
)

type Controller interface {
	Solve(context.Context, entity.Board) error
}

type Results struct {
	fx.Out

	Controller
}

type solver struct{}

func New() (Results, error) {
	return Results{
		Controller: &solver{},
	}, nil
}

func (s *solver) Solve(ctx context.Context, board entity.Board) error {
	return nil
}
