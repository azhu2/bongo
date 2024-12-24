package importer

import (
	"context"
	"fmt"
	"testing"

	"github.com/azhu2/bongo/src/testdata"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	for _, tt := range testdata.TestData {
		t.Run(tt.Name, func(t *testing.T) {
			results, _ := New()
			c := results.Controller
			board, err := c.ImportBoard(context.Background(), fmt.Sprintf("../../testdata/%s", tt.Filename))
			assert.NoError(t, err)
			assert.NotNil(t, board)
			assert.Equal(t, tt.Board.Tiles, board.Tiles, "tiles should match")
			assert.Equal(t, tt.Board.Multipliers, board.Multipliers, "multipliers should match")
			assert.Equal(t, tt.Board.BonusWord, board.BonusWord, "bonus word should match")
		})
	}
}
