package scorer

import (
	"context"
	"testing"

	"github.com/azhu2/bongo/src/testdata"
	"github.com/stretchr/testify/assert"
)

func TestScore(t *testing.T) {
	t.Run("2024-12-23 solution", func(t *testing.T) {
		results, _ := New()
		s := results.Controller
		score, err := s.Score(context.Background(), testdata.Board, testdata.Solution)
		assert.NoError(t, err)
		assert.Equal(t, testdata.Score, score)
	})
}
