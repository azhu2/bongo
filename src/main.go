package main

import (
	"context"
	"log/slog"

	"github.com/machinebox/graphql"
	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/config/secrets"
	"github.com/azhu2/bongo/src/controller/parser"
	"github.com/azhu2/bongo/src/controller/scorer"
	"github.com/azhu2/bongo/src/controller/solver"
	"github.com/azhu2/bongo/src/controller/wordlist"
	"github.com/azhu2/bongo/src/entity"
	"github.com/azhu2/bongo/src/gateway/gameimporter"
	"github.com/azhu2/bongo/src/handler"
)

const (
	date = "2024-12-29"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	fx.New(
		wordlist.Module,
		handler.Module,
		gameimporter.GraphqlModule,
		parser.Module,
		scorer.Module,
		secrets.Module,
		solver.Module,
		fx.Supply(
			graphql.NewClient(gameimporter.GraphqlEndpoint),
		),
		fx.Provide(func(c wordlist.Controller) (*entity.WordList, error) {
			return c.BuildWordList(context.Background())
		}),
		fx.Invoke(func(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, handler handler.Handler) {
			lifecycle.Append(fx.StartHook(func(_ context.Context) {
				go func() {
					solution, score, err := handler.Solve(context.Background(), date)
					if err != nil {
						slog.Error("error in solver",
							"err", err,
						)
					}
					slog.Info("solution found", "solution", solution, "score", score)
					shutdowner.Shutdown()
				}()
			}))
		}),
	).Run()
}
