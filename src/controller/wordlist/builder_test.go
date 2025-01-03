package wordlist

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/azhu2/bongo/src/gateway/wordlistimporter"
)

func TestScore(t *testing.T) {
	ctx := context.Background()
	importerGateway, err := wordlistimporter.New()
	require.NoError(t, err)
	wordlistBuilder, err := New(Params{Importer: importerGateway.Gateway})
	require.NoError(t, err)

	list, err := wordlistBuilder.BuildWordList(ctx)
	assert.NoError(t, err)
	assert.Contains(t, list.Root.Children, 'A')
}

func TestScore_TraverseWord(t *testing.T) {
	ctx := context.Background()
	importerGateway, err := wordlistimporter.New()
	require.NoError(t, err)
	wordlistBuilder, err := New(Params{Importer: importerGateway.Gateway})
	require.NoError(t, err)

	list, err := wordlistBuilder.BuildWordList(ctx)
	require.NoError(t, err)

	// Check a specific word CRAB
	node := list.Root

	node = node.Children['C']
	assert.NotEmpty(t, node)
	assert.Equal(t, []rune{'C'}, node.Fragment)
	assert.False(t, node.IsWord)
	assert.Contains(t, list.NodeMap[0]['C'], node)

	node = node.Children['R']
	assert.NotEmpty(t, node)
	assert.Equal(t, []rune{'C', 'R'}, node.Fragment)
	assert.False(t, node.IsWord)
	assert.Contains(t, list.NodeMap[1]['R'], node)

	node = node.Children['A']
	assert.NotEmpty(t, node)
	assert.Equal(t, []rune{'C', 'R', 'A'}, node.Fragment)
	assert.False(t, node.IsWord)
	assert.Contains(t, list.NodeMap[2]['A'], node)

	node = node.Children['B']
	assert.NotEmpty(t, node)
	assert.Equal(t, []rune{'C', 'R', 'A', 'B'}, node.Fragment)
	assert.True(t, node.IsWord)
	assert.Contains(t, list.NodeMap[3]['B'], node)

	node = node.Children[' ']
	assert.NotEmpty(t, node, "should have trailing empty node")
	assert.Equal(t, []rune{'C', 'R', 'A', 'B', ' '}, node.Fragment)
	assert.True(t, node.IsWord)
}

func TestScore_TraverseWordWithLeadingSpace(t *testing.T) {
	ctx := context.Background()
	importerGateway, err := wordlistimporter.New()
	require.NoError(t, err)
	wordlistBuilder, err := New(Params{Importer: importerGateway.Gateway})
	require.NoError(t, err)

	list, err := wordlistBuilder.BuildWordList(ctx)
	require.NoError(t, err)

	// Check CRAB again with leading space
	node := list.Root
	node = node.Children[' ']
	assert.NotEmpty(t, node)
	assert.Equal(t, []rune{' '}, node.Fragment)
	assert.False(t, node.IsWord)

	node = node.Children['C']
	assert.NotEmpty(t, node)
	assert.Equal(t, []rune{'C'}, node.Fragment)
	assert.False(t, node.IsWord)

	node = node.Children['R']
	assert.NotEmpty(t, node)
	assert.Equal(t, []rune{'C', 'R'}, node.Fragment)
	assert.False(t, node.IsWord)

	node = node.Children['A']
	assert.NotEmpty(t, node)
	assert.Equal(t, []rune{'C', 'R', 'A'}, node.Fragment)
	assert.False(t, node.IsWord)

	node = node.Children['B']
	assert.NotEmpty(t, node)
	assert.Equal(t, []rune{'C', 'R', 'A', 'B'}, node.Fragment)
	assert.True(t, node.IsWord)
}
