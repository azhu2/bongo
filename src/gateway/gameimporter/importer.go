package gameimporter

import (
	"context"

	"github.com/machinebox/graphql"
	"go.uber.org/fx"

	"github.com/azhu2/bongo/src/config/secrets"
)

type Gateway interface {
	ImportBoard(ctx context.Context, date string) (string, error)
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
