package handler

import (
	"context"

	"github.com/azhu2/bongo/src/controller/solver"
	"github.com/azhu2/bongo/src/gateway/importer"
	"go.uber.org/fx"
)

var Module = fx.Module("handler",
	fx.Provide(New),
)

type Handler interface {
	Solve(context.Context) error
}

type Params struct {
	fx.In

	Importer importer.Gateway
	Solver   solver.Controller
}

type Results struct {
	fx.Out

	Handler
}

type handler struct {
	importer importer.Gateway
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
	board, err := h.importer.ImportBoard(ctx)

	if err != nil {
		return err
	}

	return h.solver.Solve(ctx, board)
}
