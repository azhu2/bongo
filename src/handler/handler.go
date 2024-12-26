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

	Puzzmo puzzmo.Gateway

	Parser parser.Controller
	Solver solver.Controller
}

type Results struct {
	fx.Out

	Handler
}

type handler struct {
	puzzmo puzzmo.Gateway

	parser parser.Controller
	solver solver.Controller
}

func New(p Params) (Results, error) {
	return Results{
		Handler: &handler{
			puzzmo: p.Puzzmo,

			parser: p.Parser,
			solver: p.Solver,
		},
	}, nil
}

func (h *handler) Solve(ctx context.Context, sourceFile string) error {
	boardData, err := h.puzzmo.GetBongoBoard(ctx, "b72thg31tf")
	if err != nil {
		return err
	}

	board, err := h.parser.ParseBoard(ctx, boardData)

	if err != nil {
		return err
	}

	return h.solver.Solve(ctx, board)
}
