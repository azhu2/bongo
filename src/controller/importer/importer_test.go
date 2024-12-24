package importer

import (
	"context"
	"testing"

	"github.com/azhu2/bongo/src/testdata"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	t.Run("2024-12-23 board", func(t *testing.T) {
		results, _ := New()
		c := results.Controller
		board, err := c.ImportBoard(context.Background(), "../../testdata/example.txt")
		assert.NoError(t, err)
		assert.NotNil(t, board)
		assert.Equal(t, testdata.Board.Tiles, board.Tiles, "tiles should match")
		assert.Equal(t, testdata.Board.Multipliers, board.Multipliers, "multipliers should match")
		assert.Equal(t, testdata.Board.BonusWord, board.BonusWord, "bonus word should match")
	})
}
