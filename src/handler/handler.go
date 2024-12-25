package handler

import (
	"context"

	"github.com/azhu2/bongo/src/controller/parser"
	"github.com/azhu2/bongo/src/controller/solver"
	"go.uber.org/fx"
)

var Module = fx.Module("handler",
	fx.Provide(New),
)

type Handler interface {
	Solve(ctx context.Context, sourceFile string) error
}

type Params struct {
	fx.In

	Importer parser.Controller
	Solver   solver.Controller
}

type Results struct {
	fx.Out

	Handler
}

type handler struct {
	importer parser.Controller
	solver   solver.Controller
}

func New(p Params) (Results, error) {
	return Results{
		Handler: &handler{
			importer: p.Importer,
			solver:   p.Solver,
		},
	}, nil
}

func (h *handler) Solve(ctx context.Context, sourceFile string) error {
	board, err := h.importer.ParseBoard(ctx, sourceFile)

	if err != nil {
		return err
	}

	return h.solver.Solve(ctx, board)
}
