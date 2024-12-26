package importer

import (
	"context"

	"github.com/azhu2/bongo/src/config/secrets"
	"github.com/machinebox/graphql"
	"go.uber.org/fx"
)

type Gateway interface {
	GetBongoBoard(ctx context.Context, date string) (string, error)
}

type Params struct {
	fx.In

	secrets.Secrets
	GraphqlClient *graphql.Client
}

type Result struct {
	fx.Out

	Gateway
}
