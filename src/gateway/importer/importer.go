package importer

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/azhu2/bongo/src/entity"
	"go.uber.org/fx"
)

const sourceFile = "../example.txt"

var Module = fx.Module("importer",
	fx.Provide(New),
)

type Gateway interface {
	ImportBoard(context.Context) (entity.Board, error)
}

type Results struct {
	fx.Out

	Gateway
}

type importer struct{}

func New() (Results, error) {
	return Results{
		Gateway: &importer{},
	}, nil
}

func (i *importer) ImportBoard(ctx context.Context) (entity.Board, error) {
	data, err := i.loadData(ctx)
	if err != nil {
		return entity.Board{}, err
	}

	lines := strings.Split(data, "\n")

	return parseData(lines)
}

func (i *importer) loadData(_ context.Context) (string, error) {
	// TODO Figure out how to make query for data, but graphql might complicate this.
	raw, err := os.ReadFile(sourceFile)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func parseData(lines []string) (entity.Board, error) {
	board := entity.Board{}

	idx := 0

	// Ignore line 0
	idx++

	// Double-check board size
	if lines[idx] != "5x5" {
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
	bonus, err := parseBonusWord(lines[idx])
	if err != nil {
		return entity.Board{}, err
	}
	board.BonusWord = bonus
	idx++

	// Parse multipliers
	multipliers, err := parseMultipliers(lines[idx])
	if err != nil {
		return entity.Board{}, err
	}
	board.Multipliers = multipliers
	idx++

	// Parse tiles
	for lines[idx] != "" {
		tiles, err := parseTile(lines[idx])
		if err != nil {
			return entity.Board{}, err
		}
		board.Tiles = append(board.Tiles, tiles...)
		idx++
	}
	if len(board.Tiles) < 25 {
		return entity.Board{}, fmt.Errorf("incorrect number of tiles found: %d", len(board.Tiles))
	}

	return board, nil
}

func parseBonusWord(line string) ([][]int, error) {
	coords := strings.Split(line, " ")
	bonus := make([][]int, len(coords))
	if len(coords) == 0 {
		return nil, fmt.Errorf("no bonus word coordinates found")
	}
	for _, coord := range coords {
		x, y, err := parseCoordinate(coord)
		if err != nil {
			return nil, fmt.Errorf("unable to parse bonus word coordinate %w", err)
		}
		bonus = append(bonus, []int{x, y})
	}
	return bonus, nil
}

func parseMultipliers(line string) ([][]int, error) {
	// Default to 1s everywhere
	multipliers := make([][]int, 5)
	for i := 0; i < 5; i++ {
		multipliers[i] = []int{1, 1, 1, 1, 1}
	}

	entries := strings.Split(line, " ")
	for _, entry := range entries {
		vals := strings.Split(entry, "x")
		if len(vals) != 2 {
			return nil, fmt.Errorf("unable to parse multiplier: %s", entry)
		}
		x, y, err := parseCoordinate(vals[0])
		if err != nil {
			return nil, fmt.Errorf("unable to parse multiplier coordinate: %s %w", vals[0], err)
		}
		multiplier, err := strconv.Atoi(vals[1])
		if err != nil {
			return nil, fmt.Errorf("unalbe to parse multiplier value: %s %w", vals[1], err)
		}
		multipliers[x][y] = multiplier
	}
	return multipliers, nil
}

func parseTile(line string) ([]entity.Tile, error) {
	parts := strings.Split(line, "x")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unable to parse tile: %s", line)
	}
	letter := ([]rune)(parts[0])
	if len(letter) != 1 {
		return nil, fmt.Errorf("unable to parse tile letter: %s", line)
	}
	parts = strings.Split(parts[1], ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unable to parse tile: %s", line)
	}
	count, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("unable to parse tile count: %s %w", line, err)
	}
	// Trim occasional parenthetical number. Doesn't seem to serve a purpose
	valueStr := parts[1]
	if strings.Contains(valueStr, "(") {
		valueStr = parts[1][:strings.Index(parts[1], "(")]
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse tile value: %s %w", line, err)
	}
	tiles := make([]entity.Tile, count)
	for i := 0; i < count; i++ {
		tiles[i] = entity.Tile{
			Letter: letter[0],
			Value:  value,
		}
	}
	return tiles, nil
}

func parseCoordinate(coord string) (int, int, error) {
	stripped := strings.Replace(strings.Replace(coord, "(", "", 1), ")", "", 1)
	values := strings.Split(stripped, ",")
	if len(values) != 2 {
		return 0, 0, fmt.Errorf("unable to parse coordinate: %s", coord)
	}
	x, perr := strconv.Atoi(values[0])
	if perr != nil {
		return 0, 0, fmt.Errorf("unable to parse coordinate: %s %w", coord, perr)
	}
	y, perr := strconv.Atoi(values[1])
	if perr != nil {
		return 0, 0, fmt.Errorf("unable to parse coordinate: %s %w", coord, perr)
	}
	return x, y, nil
}
