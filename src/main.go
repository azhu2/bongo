package main

import (
	"context"
	"fmt"

	"github.com/azhu2/bongo/src/controller/parser"
	"github.com/azhu2/bongo/src/controller/solver"
	"github.com/azhu2/bongo/src/gateway/importer"
	"github.com/azhu2/bongo/src/handler"
	"github.com/machinebox/graphql"
	"go.uber.org/fx"
)

const sourceFile = "testdata/2024-12-23.txt"

func main() {
	fx.New(
		handler.Module,
		importer.GraphqlModule,
		parser.Module,
		solver.Module,
		fx.Provide(
			func() *graphql.Client { return graphql.NewClient(importer.GraphqlEndpoint) },
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
