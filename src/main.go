package main

import (
	"context"

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
				handler.Solve(ctx)
				shutdowner.Shutdown()
			}))
		}),
	).Run()
}
