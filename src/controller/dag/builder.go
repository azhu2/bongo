package dag

import (
	"context"

	"github.com/azhu2/bongo/src/entity"
	"github.com/azhu2/bongo/src/gateway/wordlistimporter"
	"go.uber.org/fx"
)

var Module = fx.Module("dagbuilder",
	wordlistimporter.Module,
	fx.Provide(New),
)

type Controller interface {
	BuildDAG(ctx context.Context, tiles map[rune]entity.Tile) (*entity.WordListDAG, error)
}

type Params struct {
	fx.In

	Importer wordlistimporter.Gateway
}

type Result struct {
	fx.Out

	Controller
}

type controller struct {
	importer wordlistimporter.Gateway
}

func New(p Params) (Result, error) {
	return Result{
		Controller: &controller{
			importer: p.Importer,
		},
	}, nil
}

func (c *controller) BuildDAG(ctx context.Context, tiles map[rune]entity.Tile) (*entity.WordListDAG, error) {

	return nil, nil
}
