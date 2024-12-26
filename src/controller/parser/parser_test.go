package parser

import (
	"context"
	"testing"

	"github.com/azhu2/bongo/src/gateway/gameimporter"
	"github.com/azhu2/bongo/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBoard(t *testing.T) {
	for _, tt := range testdata.TestData {
		t.Run(tt.Date, func(t *testing.T) {
			data := extractBoardData(t, tt.Date)

			result, _ := New()
			c := result.Controller
			board, err := c.ParseBoard(context.Background(), data)
			assert.NoError(t, err)
			assert.NotNil(t, board)
			assert.Equal(t, tt.Board.Tiles, board.Tiles, "tiles should match")
			assert.Equal(t, tt.Board.Multipliers, board.Multipliers, "multipliers should match")
			assert.Equal(t, tt.Board.BonusWord, board.BonusWord, "bonus word should match")
		})
	}
}

func extractBoardData(t *testing.T, date string) string {
	importer, err := gameimporter.NewFile(gameimporter.Params{})
	require.NoError(t, err)
	data, err := importer.Gateway.ImportBoard(context.Background(), date)
	require.NoError(t, err)
	return data
}
