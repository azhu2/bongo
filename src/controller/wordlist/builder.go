package wordlist

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/entity"
	"github.com/azhu2/bongo/src/gateway/wordlistimporter"
)

var Module = fx.Module("wordlist",
	wordlistimporter.Module,
	fx.Provide(New),
)

type Controller interface {
	BuildWordList(ctx context.Context) (*entity.WordList, error)
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

func (c *controller) BuildWordList(ctx context.Context) (*entity.WordList, error) {
	wordList, err := c.importer.ImportWordList(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not import word list %w", err)
	}

	root := entity.DAGNode{
		Fragment: []rune{},
		Children: make(map[rune]*entity.DAGNode),
	}
	nodeMap := make(map[int]map[rune][]*entity.DAGNode)

	for _, word := range wordList {
		node := &root
		stack := entity.Stack[*entity.DAGNode]{}

		for i, letter := range word {
			if child, ok := node.Children[letter]; ok {
				// Move current pointer to existing child node
				node = child
			} else {
				// Create new child node
				child := entity.DAGNode{
					Fragment: append(slices.Clone(node.Fragment), letter),
					Children: make(map[rune]*entity.DAGNode),
				}
				node.Children[letter] = &child
				node = &child

				// Add to node map
				var mapEntry map[rune][]*entity.DAGNode
				if mapEntry = nodeMap[i]; mapEntry == nil {
					mapEntry = map[rune][]*entity.DAGNode{}
				}
				mapEntry[letter] = append(mapEntry[letter], &child)
				nodeMap[i] = mapEntry
			}
			stack.Push(node)
		}
		node.IsWord = true
	}

	slog.Debug("processed words into DAG")

	return &entity.WordList{
		Root:    &root,
		NodeMap: nodeMap,
	}, nil
}
