package main

import (
	"context"
	"fmt"

	"github.com/azhu2/bongo/src/controller/parser"
	"github.com/azhu2/bongo/src/controller/solver"
	"github.com/azhu2/bongo/src/gateway/graphql"
	"github.com/azhu2/bongo/src/handler"
	graphqllib "github.com/machinebox/graphql"
	"go.uber.org/fx"
)

const sourceFile = "testdata/2024-12-23.txt"

func main() {
	fx.New(
		handler.Module,
		parser.Module,
		graphql.Module,
		solver.Module,
		fx.Provide(
			func() *graphqllib.Client { return graphqllib.NewClient(graphql.Endpoint) },
		),
		fx.Invoke(func(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, handler handler.Handler) {
			lifecycle.Append(fx.StartHook(func(ctx context.Context) {
				err := handler.Solve(ctx, sourceFile)
				if err != nil {
					fmt.Println(err)
				}
				shutdowner.Shutdown()
			}))
		}),
	).Run()
}
