package parser

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/entity"
)

var (
	boardSizeRegex  = regexp.MustCompile(`(\d)x(\d)`)
	coordinateRegex = regexp.MustCompile(`\((\d),(\d)\)`)        // (1,4)
	multiplierRegex = regexp.MustCompile(`(\((\d),(\d)\))x(\d)`) // (1,4)x2
	tileRegex       = regexp.MustCompile(`(\w)x(\d):(\d+)`)      // Gx2:45(10) - final parenthetical part is ignored
)

var Module = fx.Module("parser",
	fx.Provide(New),
)

type Controller interface {
	ParseBoard(ctx context.Context, boardData string) (entity.Board, error)
}

type Result struct {
	fx.Out

	Controller
}

type parser struct{}

func New() (Result, error) {
	return Result{
		Controller: &parser{},
	}, nil
}

func (i *parser) ParseBoard(ctx context.Context, boardData string) (entity.Board, error) {
	board := entity.Board{}
	lines := strings.Split(boardData, "\n")

	idx := 0

	// Ignore line 0
	idx++

	// Double-check board size
	sizeMatch := boardSizeRegex.FindStringSubmatch(lines[idx])
	if len(sizeMatch) == 0 ||
		sizeMatch[1] != strconv.Itoa(entity.BoardSize) ||
		sizeMatch[2] != strconv.Itoa(entity.BoardSize) {
		return entity.Board{}, fmt.Errorf("unexpected board size: %s", lines[idx])
	}
	idx++

	// Skip seed words
	idx++

	// Skip empty line
	idx++

	// Skip par
	idx++

	// Parse bonus word
	bonus, err := parseBonusWord(strings.TrimSpace(lines[idx]))
	if err != nil {
		return entity.Board{}, err
	}
	board.BonusWord = bonus
	idx++

	// Parse multipliers
	multipliers, err := parseMultipliers(strings.TrimSpace(lines[idx]))
	if err != nil {
		return entity.Board{}, err
	}
	board.Multipliers = multipliers
	idx++

	// Parse tiles
	board.Tiles = make(map[rune]entity.Tile)
	tileCount := 0
	for lines[idx] != "" {
		letter, tile, err := parseTile(strings.TrimSpace(lines[idx]))
		if err != nil {
			return entity.Board{}, err
		}
		if existing, ok := board.Tiles[letter]; ok {
			// Not sure if multiple stacks in UI show up twice or not
			if existing.Value != tile.Value {
				return entity.Board{}, fmt.Errorf("duplicate tile with different value: %c", letter)
			}
			existing.Count += tile.Count
			board.Tiles[letter] = existing
		} else {
			board.Tiles[letter] = tile
		}
		tileCount += tile.Count
		idx++
	}
	if tileCount < entity.BoardSize*entity.BoardSize {
		return entity.Board{}, fmt.Errorf("incorrect number of tiles found: %d", len(board.Tiles))
	}

	return board, nil
}

func parseBonusWord(line string) ([][]int, error) {
	matches := coordinateRegex.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("unable to parse bonus word coordinates %s", line)
	}

	bonus := make([][]int, len(matches))
	for i, match := range matches {
		x, y, err := parseCoordinate(match[0])
		if err != nil {
			return nil, fmt.Errorf("unable to parse bonus word coordinate %w", err)
		}
		bonus[i] = []int{x, y}
	}
	return bonus, nil
}

func parseMultipliers(line string) ([][]int, error) {
	// Default to 1s everywhere
	multipliers := make([][]int, entity.BoardSize)
	for i := 0; i < entity.BoardSize; i++ {
		multipliers[i] = []int{1, 1, 1, 1, 1}
	}

	matches := multiplierRegex.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("unable to parse multipliers %s", line)
	}
	for _, match := range matches {
		coord := match[1]
		x, y, err := parseCoordinate(coord)
		if err != nil {
			return nil, fmt.Errorf("unable to parse multiplier coordinate: %s %w", coord, err)
		}
		multiplier, err := strconv.Atoi(match[4])
		if err != nil {
			return nil, fmt.Errorf("unable to parse multiplier value: %s %w", match[4], err)
		}
		multipliers[x][y] = multiplier
	}
	return multipliers, nil
}

func parseTile(line string) (rune, entity.Tile, error) {
	match := tileRegex.FindStringSubmatch(line)
	if len(match) == 0 {
		return 0, entity.Tile{}, fmt.Errorf("unable to parse tile: %s", line)
	}

	letter := match[1][0]
	count, err := strconv.Atoi(match[2])
	if err != nil {
		return 0, entity.Tile{}, fmt.Errorf("unable to parse tile count: %s %w", line, err)
	}
	value, err := strconv.Atoi(match[3])
	if err != nil {
		return 0, entity.Tile{}, fmt.Errorf("unable to parse tile value: %s %w", line, err)
	}
	return rune(letter), entity.Tile{
		Value: value,
		Count: count,
	}, nil
}

// Flip the coordinate grid
// Input data is (x,y)/(col,row) with (0,0) as bottom left
// We'll flip to (row,col) with (0,0) as top left
func parseCoordinate(coord string) (int, int, error) {
	match := coordinateRegex.FindStringSubmatch(coord)
	if len(match) == 0 {
		return 0, 0, fmt.Errorf("unable to parse coordinate: %s", coord)
	}
	x, perr := strconv.Atoi(match[1])
	if perr != nil {
		return 0, 0, fmt.Errorf("unable to parse coordinate: %s %w", coord, perr)
	}
	y, perr := strconv.Atoi(match[2])
	if perr != nil {
		return 0, 0, fmt.Errorf("unable to parse coordinate: %s %w", coord, perr)
	}
	return entity.BoardSize - y - 1, x, nil
}
