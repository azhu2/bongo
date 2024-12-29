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
	letterValues := make([][]int, entity.BoardSize)
	for i := 0; i < entity.BoardSize; i++ {
		letterValues[i] = make([]int, entity.BoardSize)
		for j := 0; j < entity.BoardSize; j++ {
			letterValues[i][j] = -1
		}
	}

	// Count rows
	for rowIdx, row := range solution.Rows() {
		wordScore := 0
		for colIdx, letter := range row {
			letterScore, err := scoreLetter(ctx, board, availableLetters, rowIdx, colIdx, letter)
			if err != nil {
				if errors.Is(err, InvalidLetterError{}) {
					wildcardCount++
				}
				if wildcardCount > entity.MaxWildcards {
					return 0, err
				}
			} else {
				letterValues[rowIdx][colIdx] = letterScore
			}
			wordScore += letterScore
		}
		multiplier := s.wordMultiplier(ctx, string(row))
		score += (int)(math.Ceil(multiplier * float64(wordScore)))
		if multiplier == 0 {
			// Return letters to availability pool if word is invalid
			for _, letter := range row {
				if availableLetters[letter] > 0 {
					availableLetters[letter]++
				}
			}
		}
	}

	// Count bonus word
	bonusLetters := make([]rune, len(board.BonusWord))
	bonusScore := 0
	for i, coords := range board.BonusWord {
		rowIdx := coords[0]
		colIdx := coords[1]
		letter := solution.Get(rowIdx, colIdx)
		bonusLetters[i] = letter
		letterScore := letterValues[rowIdx][colIdx]
		bonusScore += letterScore
	}
	score += (int)(math.Ceil(s.wordMultiplier(ctx, string(bonusLetters)) * float64(bonusScore)))

	return score, nil
}

func scoreLetter(_ context.Context, board *entity.Board, availableLetters map[rune]int, row, col int, letter rune) (int, error) {
	if letter == ' ' {
		return 0, nil
	}
	tile, ok := board.Tiles[letter]
	if !ok || availableLetters[letter] < 1 {
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
