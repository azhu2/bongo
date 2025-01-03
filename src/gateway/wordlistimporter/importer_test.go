package wordlistimporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	ctx := context.Background()
	importerGateway, err := New()
	require.NoError(t, err)
	words, err := importerGateway.ImportWordList(ctx)
	assert.NoError(t, err)
	assert.Contains(t, words, "LAMBS")
	assert.Contains(t, words, "SPEAK")
	assert.Contains(t, words, "BACK")
}
