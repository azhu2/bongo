package dag

import (
	"context"
	"fmt"
	"log/slog"

	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/entity"
	"github.com/azhu2/bongo/src/gateway/wordlistimporter"
)

var Module = fx.Module("dagbuilder",
	wordlistimporter.Module,
	fx.Provide(New),
)

type Controller interface {
	BuildDAG(ctx context.Context) (*entity.WordListDAG, error)
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

func (c *controller) BuildDAG(ctx context.Context) (*entity.WordListDAG, error) {
	wordList, err := c.importer.ImportWordList(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not import word list %w", err)
	}

	root := entity.WordListDAG{
		Fragment: []rune{},
		Children: make(map[rune]*entity.WordListDAG),
	}

	for _, word := range wordList {
		node := &root
		stack := entity.Stack[*entity.WordListDAG]{}

		for _, letter := range word {
			if child, ok := node.Children[letter]; ok {
				node = child
			} else {
				child := &entity.WordListDAG{
					Fragment: append(node.Fragment, letter),
					Children: make(map[rune]*entity.WordListDAG),
				}
				node.Children[letter] = child
				node = child
			}
			stack.Push(node)
		}
		node.IsWord = true
	}

	slog.Debug("processed words into DAG")

	return &root, nil
}
