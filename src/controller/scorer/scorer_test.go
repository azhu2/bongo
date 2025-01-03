package scorer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/azhu2/bongo/testdata"
)

func TestScore(t *testing.T) {
	for _, tt := range testdata.TestData {
		t.Run(tt.Date, func(t *testing.T) {
			result, _ := New()
			s := result.Controller
			score, err := s.Score(context.Background(), tt.Board, tt.Solution)
			assert.NoError(t, err)
			assert.Equal(t, tt.Score, score)
		})
	}
}
