package scorer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/azhu2/bongo/src/controller/wordlist"
	"github.com/azhu2/bongo/src/gateway/wordlistimporter"
	"github.com/azhu2/bongo/testdata"
)

func TestScore(t *testing.T) {
	for _, tt := range testdata.TestData {
		t.Run(tt.Date, func(t *testing.T) {
			ctx := context.Background()
			importerGateway, err := wordlistimporter.New()
			require.NoError(t, err)
			wordlistBuilder, err := wordlist.New(wordlist.Params{Importer: importerGateway.Gateway})
			require.NoError(t, err)
			wordlist, err := wordlistBuilder.BuildWordList(ctx)
			require.NoError(t, err)
			result, _ := New(Params{WordList: wordlist})
			s := result.Controller

			score, err := s.Score(context.Background(), tt.Board, tt.Solution)
			assert.NoError(t, err)
			assert.Equal(t, tt.Score, score)
		})
	}
}
