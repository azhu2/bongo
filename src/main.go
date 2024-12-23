package main

import (
	"context"
	"fmt"

	"github.com/azhu2/bongo/src/controller/solver"
	"github.com/azhu2/bongo/src/gateway/importer"
	"github.com/azhu2/bongo/src/handler"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		handler.Module,
		importer.Module,
		solver.Module,
		fx.Invoke(func(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, handler handler.Handler) {
			lifecycle.Append(fx.StartHook(func(ctx context.Context) {
				err := handler.Solve(ctx)
				if err != nil {
					fmt.Println(err)
				}
				shutdowner.Shutdown()
			}))
		}),
	).Run()
}
