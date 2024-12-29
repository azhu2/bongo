package scorer

import (
	"context"
	"errors"
	"math"
	"strings"

	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/entity"
)

const (
	commonMultiplier = 1.3
)

var Module = fx.Module("scorer",
	fx.Provide(New),
)

type Controller interface {
	Score(context.Context, *entity.Board, entity.Solution) (int, error)
}

type Params struct {
	fx.In

	WordList *entity.WordList
}

type Result struct {
	fx.Out

	Controller
}

type scorer struct {
	wordList *entity.WordList
}

func New(p Params) (Result, error) {
	return Result{
		Controller: &scorer{
			wordList: p.WordList,
		},
	}, nil
}

func (s *scorer) Score(ctx context.Context, board *entity.Board, solution entity.Solution) (int, error) {
	score := 0
	wildcardCount := 0

	availableLetters := map[rune]int{}
	for letter, tile := range board.Tiles {
		availableLetters[letter] = tile.Count
	}
	// Avoid double-counting availability for bonus words
	bonusRowsCounted := map[int]bool{}

	for rowIdx, row := range solution.Rows() {
		wordScore := 0
		for colIdx, letter := range row {
			letterScore, err := scoreLetter(ctx, board, availableLetters, rowIdx, colIdx, letter, false)
			if err != nil {
				if errors.Is(err, InvalidLetterError{}) {
					wildcardCount++
				}
				if wildcardCount > entity.MaxWildcards {
					return 0, err
				}
			}
			wordScore += letterScore
		}
		multiplier := s.wordMultiplier(ctx, string(row))
		score += (int)(math.Ceil(multiplier * float64(wordScore)))
		if multiplier != 0 {
			bonusRowsCounted[rowIdx] = true
		} else {
			// Return invalid word letters to availability pool
			for _, letter := range row {
				if availableLetters[letter] > 0 {
					availableLetters[letter]++
				}
			}
		}
	}

	bonusLetters := make([]rune, len(board.BonusWord))
	bonusScore := 0
	for i, coords := range board.BonusWord {
		rowIdx := coords[0]
		colIdx := coords[1]
		letter := solution.Get(rowIdx, colIdx)
		bonusLetters[i] = letter
		shouldSkipAvailabilityCheck := bonusRowsCounted[rowIdx]
		letterScore, err := scoreLetter(ctx, board, availableLetters, rowIdx, colIdx, letter, shouldSkipAvailabilityCheck)
		if err != nil {
			if !errors.Is(err, InvalidLetterError{}) {
				return 0, err
			} else {
				continue
			}
		}
		bonusScore += letterScore
	}
	score += (int)(math.Ceil(s.wordMultiplier(ctx, string(bonusLetters)) * float64(bonusScore)))

	return score, nil
}

func scoreLetter(_ context.Context, board *entity.Board, availableLetters map[rune]int, row, col int, letter rune, shouldSkipAvailabilityCheck bool) (int, error) {
	if letter == ' ' {
		return 0, nil
	}
	tile, ok := board.Tiles[letter]
	if !shouldSkipAvailabilityCheck && (!ok || availableLetters[letter] < 1) {
		return 0, InvalidLetterError{letter: letter}
	}
	availableLetters[letter]--
	return board.Multipliers[row][col] * tile.Value, nil
}

func (s *scorer) wordMultiplier(ctx context.Context, word string) float64 {
	// Assume no gaps. I'm pretty sure only one word per line is counted.
	trimmed := strings.TrimSpace(word)
	if !s.isWord(ctx, trimmed) {
		return 0
	}
	if s.isCommon(ctx, trimmed) {
		return commonMultiplier
	}
	return 1
}

func (s *scorer) isWord(_ context.Context, word string) bool {
	node := s.wordList.Root
	for _, letter := range word {
		if child := node.Children[letter]; child != nil {
			node = child
			continue
		}
		return false
	}
	return node.IsWord
}

func (s *scorer) isCommon(_ context.Context, _ string) bool {
	// TODO Implement
	return true
}
