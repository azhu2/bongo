package main

import (
	"github.com/azhu2/bongo/src/gateway/importer"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		importer.Module,
	).Run()
}
