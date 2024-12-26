package solver

import (
	"context"

	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/entity"
)

var Module = fx.Module("solver",
	fx.Provide(New),
)

type Controller interface {
	Solve(context.Context, entity.Board) error
}

type Result struct {
	fx.Out

	Controller
}

type solver struct{}

func New() (Result, error) {
	return Result{
		Controller: &solver{},
	}, nil
}

func (s *solver) Solve(ctx context.Context, board entity.Board) error {
	return nil
}
