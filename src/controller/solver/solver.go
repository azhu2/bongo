package solver

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"strings"

	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/controller/scorer"
	"github.com/azhu2/bongo/src/entity"
)

const (
	// Only consider bonus words at least this percent as good as the best one found so far
	bonusCandidateMultiplier = 0.75
)

var Module = fx.Module("solver",
	fx.Provide(New),
)

type Controller interface {
	Solve(context.Context, *entity.Board) (entity.Solution, error)
}

type Params struct {
	fx.In

	Scorer   scorer.Controller
	WordList *entity.WordList
}

type Result struct {
	fx.Out

	Controller
}

type solver struct {
	scorer   scorer.Controller
	wordList *entity.WordList
}

func New(p Params) (Result, error) {
	return Result{
		Controller: &solver{
			scorer:   p.Scorer,
			wordList: p.WordList,
		},
	}, nil
}

func (s *solver) Solve(ctx context.Context, board *entity.Board) (entity.Solution, error) {
	// The DAG can't really work backwards with letters in the middle, so this is a problem...
	// Ignore for now, but I'd really like to use it to narrow the solution space.
	// Maybe make the DAG doubly-linked?
	// s.generateBonusCandidates(ctx, board, words)

	availableLetters := map[rune]int{}
	for letter, tile := range board.Tiles {
		availableLetters[letter] = tile.Count
	}
	globalBest := 0
	solution := s.evaluateRow(ctx, board, partialSolution{
		solution:         entity.EmptySolution(),
		availableLetters: availableLetters,
		wildcardCount:    0,
		rowsFilled:       0,
	}, &globalBest)

	return solution, nil
}

type bonusCandidate struct {
	solution entity.Solution
	score    int
}

func (s *solver) generateBonusCandidates(ctx context.Context, board *entity.Board) []entity.Solution {
	candidates := []bonusCandidate{}

	maxValue := 0
	nodes := entity.Stack[*entity.DAGNode]{}
	nodes.Push(s.wordList.Root)
	for !nodes.IsEmpty() {
		cur := nodes.Pop()
		for _, child := range cur.Children {
			nodes.Push(child)
		}
		// Assume all bonus word letters must be filled (sometimes false)
		if cur.IsWord && len(cur.Fragment) == len(board.BonusWord) {
			candidate := entity.EmptySolution()
			for i, b := range board.BonusWord {
				candidate.Set(b[0], b[1], cur.Fragment[i])
			}
			score, err := s.scorer.Score(ctx, board, candidate)
			if err != nil {
				// TODO Type this
				if !strings.Contains(err.Error(), "solution has invalid letter") {
					slog.Error("unable to score bonus word candidate; continuing",
						"candidate", candidate,
						"err", err,
					)
				}
				continue
			}
			if score > maxValue {
				maxValue = score
			}
			if score >= int(bonusCandidateMultiplier*float64(maxValue)) {
				candidates = append(candidates, bonusCandidate{solution: candidate, score: score})
			}
		}
	}

	// Filter once more to avoid repeatedly rebalancing and trimming a tree
	// Another option is a priority queue with a fixed size
	bonusBoards := []entity.Solution{}
	for _, candidate := range candidates {
		if candidate.score >= int(bonusCandidateMultiplier*float64(maxValue)) {
			bonusBoards = append(bonusBoards, candidate.solution)
		}
	}

	slog.Debug("generated bonus word candidates", "count", len(bonusBoards))
	return bonusBoards
}

type partialSolution struct {
	solution         entity.Solution
	availableLetters map[rune]int
	wildcardCount    int
	rowsFilled       int
}

type partialRow struct {
	node             *entity.DAGNode
	availableLetters map[rune]int
	wildcardCount    int
}

func (s *solver) evaluateRow(ctx context.Context, board *entity.Board, partial partialSolution, globalBest *int) entity.Solution {
	// Base case
	if partial.rowsFilled == entity.BoardSize {
		return partial.solution
	}

	best := partial.solution
	bestScore := 0

	rowCandidates := entity.Stack[partialRow]{}
	rowCandidates.Push(partialRow{
		node:             s.wordList.Root,
		availableLetters: partial.availableLetters,
		wildcardCount:    partial.wildcardCount,
	})
	for !rowCandidates.IsEmpty() {
		cur := rowCandidates.Pop()
		for nextLetter, childNode := range cur.node.Children {
			isLetterAvailable := cur.availableLetters[nextLetter] > 0
			if !isLetterAvailable && cur.wildcardCount >= entity.MaxWildcards {
				continue
			}
			remainingLetters := maps.Clone(cur.availableLetters)
			wildcardCount := cur.wildcardCount
			if isLetterAvailable {
				remainingLetters[nextLetter]--
			} else {
				wildcardCount++
			}
			rowCandidates.Push(partialRow{
				node:             childNode,
				availableLetters: remainingLetters,
				wildcardCount:    wildcardCount,
			})
		}
		// Start with assuming only 5-letter words to save on search space
		if cur.node.IsWord && len(cur.node.Fragment) == entity.BoardSize {
			nextPartial := slices.Clone(partial.solution)
			nextPartial.SetRow(partial.rowsFilled, cur.node.Fragment)
			remainingLetters := maps.Clone(partial.availableLetters)
			for _, letter := range cur.node.Fragment {
				remainingLetters[letter]--
			}
			candidate := s.evaluateRow(ctx, board, partialSolution{
				solution:         nextPartial,
				availableLetters: remainingLetters,
				wildcardCount:    cur.wildcardCount,
				rowsFilled:       partial.rowsFilled + 1,
			}, globalBest)
			score, err := s.scorer.Score(ctx, board, candidate)
			if err != nil {
				// swallow error and continue
				slog.Error("invalid board generated", "board", candidate, "err", err)
				continue
			}
			if score > bestScore {
				best = candidate
				bestScore = score
				if score > *globalBest {
					slog.Debug("new best board", "board", candidate, "score", score)
					*globalBest = score
				}
			}
		}
	}

	return best
}
