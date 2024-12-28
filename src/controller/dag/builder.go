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
	wordList, err := c.importer.ImportWordList(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not import word list %w", err)
	}

	root := entity.WordListDAG{
		Fragment: []rune{},
		Children: make(map[rune]*entity.WordListDAG),
		MaxValue: 0,
	}

	for _, word := range wordList {
		node := &root
		stack := entity.Stack[*entity.WordListDAG]{}

		// Add words to DAG
		for _, letter := range word {
			if child, ok := node.Children[letter]; ok {
				node = child
			} else {
				child := &entity.WordListDAG{
					Fragment: append(node.Fragment, letter),
					Children: make(map[rune]*entity.WordListDAG),
					MaxValue: 0,
				}
				node.Children[letter] = child
				node = child
			}
			stack.Push(node)
		}

		// Process max value per node in reverse order
		node.MaxValue = tiles[node.Fragment[len(word)-1:][0]].Value
		for !stack.IsEmpty() {
			node = stack.Pop()
			for _, child := range node.Children {
				if child.MaxValue > node.MaxValue {
					node.MaxValue = child.MaxValue
				}
			}
		}
	}

	// Also process root max. Probably not needed but nice to make it clean.
	for _, child := range root.Children {
		if child.MaxValue > root.MaxValue {
			root.MaxValue = child.MaxValue
		}
	}

	slog.Debug("processed words into DAG")

	return &root, nil
}
