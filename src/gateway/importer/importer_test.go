package importer

import (
	"context"
	"testing"

	"github.com/azhu2/bongo/src/config/secrets"
	"github.com/azhu2/bongo/src/testdata"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Compare results of graphql and file importer
func TestImportBoard(t *testing.T) {
	for _, tt := range testdata.TestData {
		t.Run(tt.Date, func(t *testing.T) {
			godotenv.Load("../../.env")
			secretResult, err := secrets.New()
			require.NoError(t, err)
			params := Params{
				Secrets:       secretResult.Secrets,
				GraphqlClient: graphql.NewClient(GraphqlEndpoint),
			}
			graphqlResult, err := NewGraphql(params)
			require.NoError(t, err)
			fileResult, err := NewFile(params)
			require.NoError(t, err)

			graphqlImport, err := graphqlResult.Gateway.ImportBoard(context.Background(), tt.Date)
			assert.NoError(t, err)
			fileImport, err := fileResult.Gateway.ImportBoard(context.Background(), tt.Date)
			assert.NoError(t, err)

			assert.Equal(t, graphqlImport, fileImport)
		})
	}
}
