package handler

import (
	"context"

	"github.com/azhu2/bongo/src/controller/parser"
	"github.com/azhu2/bongo/src/controller/solver"
	"github.com/azhu2/bongo/src/gateway/puzzmo"
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
	Puzzmo   puzzmo.Gateway
	Solver   solver.Controller
}

type Results struct {
	fx.Out

	Handler
}

type handler struct {
	importer parser.Controller
	puzzmo   puzzmo.Gateway
	solver   solver.Controller
}

func New(p Params) (Results, error) {
	return Results{
		Handler: &handler{
			importer: p.Importer,
			puzzmo:   p.Puzzmo,
			solver:   p.Solver,
		},
	}, nil
}

func (h *handler) Solve(ctx context.Context, sourceFile string) error {
	_, err := h.puzzmo.GetBongoBoard(ctx, "b72thg31tf")
	if err != nil {
		return err
	}

	board, err := h.importer.ParseBoard(ctx, sourceFile)

	if err != nil {
		return err
	}

	return h.solver.Solve(ctx, board)
}
