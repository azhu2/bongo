package handler

import (
	"context"

	"github.com/azhu2/bongo/src/controller/parser"
	"github.com/azhu2/bongo/src/controller/solver"
	"github.com/azhu2/bongo/src/gateway/gameimporter"
	"go.uber.org/fx"
)

var Module = fx.Module("handler",
	fx.Provide(New),
)

type Handler interface {
	Solve(ctx context.Context, date string) error
}

type Params struct {
	fx.In

	GameImporter gameimporter.Gateway

	Parser parser.Controller
	Solver solver.Controller
}

type Result struct {
	fx.Out

	Handler
}

type handler struct {
	gameImporter gameimporter.Gateway

	parser parser.Controller
	solver solver.Controller
}

func New(p Params) (Result, error) {
	return Result{
		Handler: &handler{
			gameImporter: p.GameImporter,

			parser: p.Parser,
			solver: p.Solver,
		},
	}, nil
}

func (h *handler) Solve(ctx context.Context, date string) error {
	boardData, err := h.gameImporter.ImportBoard(ctx, date)
	if err != nil {
		return err
	}

	board, err := h.parser.ParseBoard(ctx, boardData)

	if err != nil {
		return err
	}

	return h.solver.Solve(ctx, board)
}
