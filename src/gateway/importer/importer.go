package importer

import (
	"context"
	"os"

	"github.com/azhu2/bongo/src/entity"
	"go.uber.org/fx"
)

var Module = fx.Module("importer",
	fx.Provide(New),
)

type Gateway interface {
	ImportBoard(context.Context) (entity.Board, error)
}

type Results struct {
	fx.Out

	Gateway
}

type importer struct{}

func New() (Results, error) {
	return Results{
		Gateway: &importer{},
	}, nil
}

func (i *importer) ImportBoard(ctx context.Context) (entity.Board, error) {
	i.loadData(ctx)

	return entity.Board{}, nil
}

func (i *importer) loadData(_ context.Context) (string, error) {
	// TODO Figure out how to make query for data, but graphql might complicate this.
	os.ReadFile("example.txt")

	return "", nil
}
