package handler

import (
	"context"

	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/controller/dag"
	"github.com/azhu2/bongo/src/controller/parser"
	"github.com/azhu2/bongo/src/controller/scorer"
	"github.com/azhu2/bongo/src/controller/solver"
	"github.com/azhu2/bongo/src/entity"
	"github.com/azhu2/bongo/src/gateway/gameimporter"
)

var Module = fx.Module("handler",
	fx.Provide(New),
)

type Handler interface {
	Solve(ctx context.Context, date string) (entity.Solution, int, error)
}

type Params struct {
	fx.In

	GameImporter gameimporter.Gateway

	DAGBuilder dag.Controller
	Parser     parser.Controller
	Scorer     scorer.Controller
	Solver     solver.Controller
}

type Result struct {
	fx.Out

	Handler
}

type handler struct {
	gameImporter gameimporter.Gateway

	dagBuilder dag.Controller
	parser     parser.Controller
	scorer     scorer.Controller
	solver     solver.Controller
}

func New(p Params) (Result, error) {
	return Result{
		Handler: &handler{
			gameImporter: p.GameImporter,

			dagBuilder: p.DAGBuilder,
			parser:     p.Parser,
			scorer:     p.Scorer,
			solver:     p.Solver,
		},
	}, nil
}

func (h *handler) Solve(ctx context.Context, date string) (entity.Solution, int, error) {
	boardData, err := h.gameImporter.ImportBoard(ctx, date)
	if err != nil {
		return nil, 0, err
	}

	board, err := h.parser.ParseBoard(ctx, boardData)

	if err != nil {
		return nil, 0, err
	}

	solution, err := h.solver.Solve(ctx, board)
	if err != nil {
		return nil, 0, err
	}

	score, err := h.scorer.Score(ctx, board, solution)
	if err != nil {
		return nil, 0, err
	}

	return solution, score, err
}
