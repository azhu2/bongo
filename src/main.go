package main

import (
	"context"
	"log/slog"

	"github.com/machinebox/graphql"
	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/config/secrets"
	"github.com/azhu2/bongo/src/controller/dag"
	"github.com/azhu2/bongo/src/controller/parser"
	"github.com/azhu2/bongo/src/controller/solver"
	"github.com/azhu2/bongo/src/gateway/gameimporter"
	"github.com/azhu2/bongo/src/handler"
)

const date = "2024-12-28"

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	fx.New(
		dag.Module,
		handler.Module,
		gameimporter.GraphqlModule,
		parser.Module,
		secrets.Module,
		solver.Module,
		fx.Supply(
			graphql.NewClient(gameimporter.GraphqlEndpoint),
		),
		fx.Invoke(func(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, handler handler.Handler) {
			lifecycle.Append(fx.StartHook(func(ctx context.Context) {
				err := handler.Solve(ctx, date)
				if err != nil {
					slog.Error("error in solver",
						"err", err,
					)
				}
				shutdowner.Shutdown()
			}))
		}),
	).Run()
}
