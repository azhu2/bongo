package puzzmo

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	"go.uber.org/fx"
)

const (
	endpoint  = "https://www.puzzmo.com/_api/prod/graphql"
	bongoSlug = "play:/bongo/%s"
)

var Module = fx.Module("puzzmo",
	fx.Provide(New),
)

type Gateway interface {
	GetBongoBoard(ctx context.Context, gameSlug string) (string, error)
}

type Results struct {
	fx.Out

	Gateway
}

type gateway struct {
	graphqlClient *graphql.Client
}

func New() (Results, error) {
	return Results{
		Gateway: &gateway{
			graphqlClient: graphql.NewClient(endpoint),
		},
	}, nil
}

func (g *gateway) GetBongoBoard(ctx context.Context, gameSlug string) (string, error) {
	req := graphql.NewRequest(`
		query PlayGameScreenQuery(
			$finderKey: String!
			$gameContext: StartGameContext!
		) {
			startOrFindGameplay(finderKey: $finderKey, context: $gameContext) {
				__typename
				... on ErrorableResponse {
				message
				failed
				success
				}
				... on HasGamePlayed {
				gamePlayed {
					puzzle {
					puzzle
					}
				}
				}
			}
		}
	`)
	req.Var("finderKey", fmt.Sprintf(bongoSlug, gameSlug))
	req.Var("gameContext", map[string]any{
		"partnerSlug":             nil,
		"pingOwnerForMultiplayer": true,
	})
	req.Header.Set("context-type", "application/json")

	var resp bongoResponse
	err := g.graphqlClient.Run(ctx, req, &resp)
	board := resp.StartOrFindGameplay.GamePlayed.Puzzle.Puzzle
	if err != nil || len(board) == 0 {
		return "", fmt.Errorf("unable to fetch Bongo board from Puzzmo %w", err)
	}

	return board, err
}
