package solver

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"math"
	"slices"
	"strings"
	"sync"

	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/controller/scorer"
	"github.com/azhu2/bongo/src/entity"
)

const (
	// Only consider bonus words at least this percent as good as the best one
	bonusCandidateMultiplier = 0.6
)

var Module = fx.Module("solver",
	fx.Provide(New),
)

type Controller interface {
	Solve(context.Context, *entity.Board) ([]entity.Solution, error)
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
	scorer    scorer.Controller
	wordList  *entity.WordList
	bestScore int
}

func New(p Params) (Result, error) {
	return Result{
		Controller: &solver{
			scorer:   p.Scorer,
			wordList: p.WordList,
		},
	}, nil
}

type candidateSolution struct {
	solution entity.Solution
	score    int
}

func (s *solver) Solve(ctx context.Context, board *entity.Board) ([]entity.Solution, error) {
	// Start by generating bonus words
	candidates := s.generateBonusCandidates(ctx, board)

	availableLetters := map[rune]int{}
	for letter, tile := range board.Tiles {
		availableLetters[letter] = tile.Count
	}

	var wg sync.WaitGroup
	best := []entity.Solution{}

	// Then seed the recursive row-by-row solver with bonus words already set in grid
	for _, candidate := range candidates {
		remainingLetters := maps.Clone(availableLetters)
		for _, letter := range candidate {
			if letter != ' ' {
				remainingLetters[letter]--
			}
		}
		wg.Add(1)
		solutionChan := make(chan candidateSolution)
		go func() {
			defer wg.Done()
			s.evaluateRow(ctx, board, partialSolution{
				solution:         candidate,
				availableLetters: remainingLetters,
				wildcardCount:    0,
				curRow:           0,
			}, solutionChan)
			close(solutionChan)
		}()
		for solution := range solutionChan {
			if solution.score == s.bestScore {
				slog.Debug("new best board (tied)", "board", solution.solution, "score", solution.score)
				best = append(best, solution.solution)
			} else if solution.score > s.bestScore {
				slog.Debug("new best board", "board", solution.solution, "score", solution.score)
				s.bestScore = solution.score
				best = []entity.Solution{solution.solution}
			}
		}
	}

	wg.Wait()
	return best, nil
}

func (s *solver) generateBonusCandidates(ctx context.Context, board *entity.Board) []entity.Solution {
	candidates := []candidateSolution{}

	maxValue := 0
	nodes := entity.Stack[*entity.DAGNode]{}
	nodes.Push(s.wordList.Root)
	for !nodes.IsEmpty() {
		cur := nodes.Pop()
		for letter, child := range cur.Children {
			// Assume all bonus tiles used
			if letter != ' ' {
				nodes.Push(child)
			}
		}
		if cur.IsWord && len(cur.Fragment) == len(board.BonusWord) {
			candidate := entity.EmptySolution()

			// Assume no wildcards in bonus (may not be true)
			letters := map[rune]int{}
			isWildCard := false
			for _, letter := range cur.Fragment {
				letters[letter]++
				if letters[letter] > board.Tiles[letter].Count {
					isWildCard = true
					break
				}
			}
			if isWildCard {
				continue
			}

			for i, b := range board.BonusWord {
				candidate.Set(b[0], b[1], cur.Fragment[i])
			}
			score, err := s.scorer.Score(ctx, board, candidate)
			if err != nil {
				if !errors.Is(err, scorer.InvalidLetterError{}) {
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
				candidates = append(candidates, candidateSolution{solution: candidate, score: score})
			}
		}
	}

	slices.SortFunc(candidates, func(a, b candidateSolution) int {
		return b.score - a.score
	})

	// Filter once at the end to avoid repeatedly rebalancing and trimming a tree
	// Another option is a priority queue with a fixed size
	bonusBoards := []entity.Solution{}
	logMsg := ""
	for _, candidate := range candidates {
		if candidate.score >= int(bonusCandidateMultiplier*float64(maxValue)) {
			bonusBoards = append(bonusBoards, candidate.solution)
			logMsg += strings.ReplaceAll(strings.ReplaceAll(string(candidate.solution), "|", ""), " ", "") + "|"
		}
	}

	slog.Debug("generated bonus word candidates",
		"count", len(bonusBoards),
		"words", logMsg,
	)
	return bonusBoards
}

type partialSolution struct {
	solution         entity.Solution
	availableLetters map[rune]int
	wildcardCount    int
	curRow           int
}

type partialRow struct {
	node             *entity.DAGNode
	availableLetters map[rune]int
	wildcardCount    int
}

func (s *solver) evaluateRow(ctx context.Context, board *entity.Board, partial partialSolution, solutions chan<- candidateSolution) []entity.Solution {
	// Base case
	if partial.curRow == entity.BoardSize {
		return []entity.Solution{partial.solution}
	}

	// Short-circuit if not possible to beat current max
	max := s.getTheoreticalMax(ctx, board, partial)
	if max < s.bestScore {
		return []entity.Solution{partial.solution}
	}

	best := []entity.Solution{partial.solution}
	bestScore := 0

	rowCandidates := entity.Stack[partialRow]{}
	filledCol := -1
	for col, letter := range partial.solution.GetRow(partial.curRow) {
		if letter != ' ' {
			filledCol = col
		}
	}
	if filledCol != -1 {
		// If there are tiles already filled in this row, seed from node map in word list
		filledLetter := partial.solution.Get(partial.curRow, filledCol)
		filledCandidates := s.wordList.NodeMap[filledCol][filledLetter]
		for _, candidate := range filledCandidates {
			filledSolution := slices.Clone(partial.solution)
			remainingLetters := maps.Clone(partial.availableLetters)
			wildcardCount := partial.wildcardCount
			// Backfill the earlier letters before this node
			for col, letter := range candidate.Fragment {
				if col == filledCol {
					continue
				}
				filledSolution.Set(partial.curRow, col, letter)
				if remainingLetters[letter] > 0 {
					remainingLetters[letter]--
				} else {
					wildcardCount++
				}
			}
			if wildcardCount > entity.MaxWildcards {
				continue
			}
			rowCandidates.Push(partialRow{
				node:             candidate,
				availableLetters: remainingLetters,
				wildcardCount:    wildcardCount,
			})
		}
	} else {
		// Start with blank row and root of word list
		rowCandidates.Push(partialRow{
			node:             s.wordList.Root,
			availableLetters: partial.availableLetters,
			wildcardCount:    partial.wildcardCount,
		})
	}

	for !rowCandidates.IsEmpty() {
		cur := rowCandidates.Pop()
		for nextLetter, childNode := range cur.node.Children {
			// Add valid children nodes
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

		// Only score 5-letter fragments to avoid recounting the same candidate with/without spaces.
		// This words because the DAG is padded with leading and trailing spaces.
		if cur.node.IsWord && len(cur.node.Fragment) == entity.BoardSize {
			nextPartial := slices.Clone(partial.solution)
			nextPartial.SetRow(partial.curRow, cur.node.Fragment)
			remainingLetters := maps.Clone(partial.availableLetters)
			for j, letter := range cur.node.Fragment {
				if partial.solution.Get(partial.curRow, j) != letter {
					// Don't deduct if already set (from partial)
					if remainingLetters[letter] > 0 {
						remainingLetters[letter]--
					}
				}
			}
			candidates := s.evaluateRow(ctx, board, partialSolution{
				solution:         nextPartial,
				availableLetters: remainingLetters,
				wildcardCount:    cur.wildcardCount,
				curRow:           partial.curRow + 1,
			}, solutions)
			score, err := s.scorer.Score(ctx, board, candidates[0])
			if err != nil {
				// swallow error and continue
				slog.Error("invalid board generated", "board", candidates, "err", err)
				continue
			}
			if score == bestScore {
				best = append(best, candidates...)
			} else if score >= bestScore {
				best = candidates
				bestScore = score
				if partial.curRow == entity.BoardSize-1 {
					for _, candidate := range candidates {
						solutions <- candidateSolution{
							solution: candidate,
							score:    score,
						}
					}
				}
			}
		}
	}

	return best
}

func (s *solver) getTheoreticalMax(ctx context.Context, board *entity.Board, partial partialSolution) int {
	potentialSolution := slices.Clone(partial.solution)

	// Count what's already set in partial
	score, err := s.scorer.Score(ctx, board, potentialSolution)
	if err != nil {
		// Swallow and ignore this branch
		slog.Error("Unable to score partial board",
			"board", partial.solution,
			"err", err,
		)
		return 0
	}

	// Extract sorted coords of all non-blank multpliers
	type multiplierCoord struct {
		coord      []int
		multiplier int
	}
	multiplierCoords := []multiplierCoord{}
	for row, rowData := range board.Multipliers {
		if row < partial.curRow {
			// Already counted
			continue
		}
		for col, multipler := range rowData {
			if multipler > 1 {
				multiplierCoords = append(multiplierCoords, multiplierCoord{
					coord:      []int{row, col},
					multiplier: multipler,
				})
			}
		}
	}
	slices.SortFunc(multiplierCoords, func(i, j multiplierCoord) int {
		return i.multiplier - j.multiplier
	})

	// Place letters in multplier slots
	sortedTiles := board.SortedTiles()
	tileIdx := 0
	remainingLetters := maps.Clone(partial.availableLetters)
	for _, multiplier := range multiplierCoords {
		for ; ; tileIdx++ {
			if remainingLetters[sortedTiles[tileIdx]] < 1 {
				continue
			}
			letter := sortedTiles[tileIdx]
			remainingLetters[letter]--
			potentialSolution.Set(multiplier.coord[0], multiplier.coord[1], letter)
			if remainingLetters[letter] < 1 {
				tileIdx++
			}
			break
		}
	}

	// Place remaining letters
	for row := partial.curRow; row < entity.BoardSize; row++ {
		for col := 0; col < entity.BoardSize; col++ {
			if potentialSolution.Get(row, col) != ' ' {
				continue
			}
			for ; ; tileIdx++ {
				if tileIdx >= len(sortedTiles) {
					break
				}
				if remainingLetters[sortedTiles[tileIdx]] < 1 {
					continue
				}
				letter := sortedTiles[tileIdx]
				remainingLetters[letter]--
				potentialSolution.Set(row, col, letter)
				if remainingLetters[letter] < 1 {
					tileIdx++
				}
				break
			}
		}
	}

	// Count rows
	for row := partial.curRow; row < entity.BoardSize; row++ {
		rowScore := 0
		for col, letter := range potentialSolution.GetRow(row) {
			rowScore += board.Tiles[letter].Value * board.Multipliers[row][col]
		}
		// Fudge by 1 for rounding
		score += int(math.Ceil(float64(rowScore)*entity.CommonMultiplier)) + 1
	}

	// Count bonus word if needed
	if partial.curRow <= board.BonusWord[len(board.BonusWord)-1][0] {
		bonusScore := 0
		for _, bonusCoord := range board.BonusWord {
			letter := potentialSolution.Get(bonusCoord[0], bonusCoord[1])
			bonusScore += board.Tiles[letter].Value * board.Multipliers[bonusCoord[0]][bonusCoord[1]]
		}
		// Fudge by 1 for rounding
		score += int(math.Ceil(float64(bonusScore)*entity.CommonMultiplier)) + 1
	}

	return score
}
