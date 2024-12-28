package solver

import (
	"context"
	"log/slog"
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
	Solve(context.Context, *entity.Board, *entity.WordListDAG) error
}

type Params struct {
	fx.In

	Scorer scorer.Controller
}

type Result struct {
	fx.Out

	Controller
}

type solver struct {
	scorer scorer.Controller
}

func New(p Params) (Result, error) {
	return Result{
		Controller: &solver{
			scorer: p.Scorer,
		},
	}, nil
}

func (s *solver) Solve(ctx context.Context, board *entity.Board, words *entity.WordListDAG) error {
	s.generateBonusCandidates(ctx, board, words)

	return nil
}

type bonusCandidate struct {
	solution entity.Solution
	score    int
}

func (s *solver) generateBonusCandidates(ctx context.Context, board *entity.Board, words *entity.WordListDAG) []entity.Solution {
	candidates := []bonusCandidate{}

	maxValue := 0
	nodes := entity.Stack[*entity.WordListDAG]{}
	nodes.Push(words)
	for !nodes.IsEmpty() {
		cur := nodes.Pop()
		for _, child := range cur.Children {
			nodes.Push(child)
		}
		// Assume all bonus word letters must be filled (sometimes false)
		if cur.IsWord && len(cur.Fragment) == len(board.BonusWord) {
			candidate := entity.EmptySolution()
			for i, b := range board.BonusWord {
				candidate[b[0]][b[1]] = cur.Fragment[i]
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
