package main

import (
	"context"
	"log/slog"
	"time"

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
					start := time.Now()
					serverTime, _ := time.LoadLocation("America/Chicago")
					solutions, score, err := handler.Solve(context.Background(), time.Now().In(serverTime).Format("2006-01-02"))
					if err != nil {
						slog.Error("error in solver",
							"err", err,
						)
					}
					slog.Info("solution found", "score", score, "time", time.Since(start))
					for _, solution := range solutions {
						slog.Info(solution.String())
					}
					shutdowner.Shutdown()
				}()
			}))
		}),
	).Run()
}
