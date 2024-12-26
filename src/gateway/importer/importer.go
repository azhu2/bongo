package importer

import (
	"context"

	"github.com/machinebox/graphql"
	"go.uber.org/fx"
)

type Gateway interface {
	GetBongoBoard(ctx context.Context, date string) (string, error)
}

type Params struct {
	fx.In

	GraphqlClient *graphql.Client
}

type Result struct {
	fx.Out

	Gateway
}
