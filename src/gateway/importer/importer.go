package importer

import (
	"context"

	"github.com/azhu2/bongo/src/entity"
	"go.uber.org/fx"
)

var Module = fx.Module("importer",
	fx.Provide(New),
)

type Importer interface {
	ImportBoard(context.Context) (entity.Board, error)
}

type Results struct {
	fx.Out

	Importer
}

type importer struct{}

func New() (Results, error) {
	return Results{
		Importer: &importer{},
	}, nil
}

func (i *importer) ImportBoard(ctx context.Context) (entity.Board, error) {
	return entity.Board{}, nil
}
