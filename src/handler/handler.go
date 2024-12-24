package handler

import (
	"context"

	"github.com/azhu2/bongo/src/controller/importer"
	"github.com/azhu2/bongo/src/controller/solver"
	"go.uber.org/fx"
)

const sourceFile = "testdata/example.txt"

var Module = fx.Module("handler",
	fx.Provide(New),
)

type Handler interface {
	Solve(context.Context) error
}

type Params struct {
	fx.In

	Importer importer.Controller
	Solver   solver.Controller
}

type Results struct {
	fx.Out

	Handler
}

type handler struct {
	importer importer.Controller
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

func (h *handler) Solve(ctx context.Context) error {
	board, err := h.importer.ImportBoard(ctx, sourceFile)

	if err != nil {
		return err
	}

	return h.solver.Solve(ctx, board)
}
