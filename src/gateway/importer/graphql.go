package importer

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	"go.uber.org/fx"
)

const (
	GraphqlEndpoint = "https://www.puzzmo.com/_api/prod/graphql"
	bongoSlug       = "play:/bongo/%s"
)

var GraphqlModule = fx.Module("graphqlimporter",
	fx.Provide(NewGraphql),
)

type graphqlGateway struct {
	graphqlClient *graphql.Client
}

func NewGraphql(p Params) (Results, error) {
	return Results{
		Gateway: &graphqlGateway{
			graphqlClient: p.GraphqlClient,
		},
	}, nil
}

// TODO Map date -> slug
func (g *graphqlGateway) GetBongoBoard(ctx context.Context, gameSlug string) (string, error) {
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

	var resp graphqlBoardResponse
	err := g.graphqlClient.Run(ctx, req, &resp)
	board := resp.StartOrFindGameplay.GamePlayed.Puzzle.Puzzle
	if err != nil || len(board) == 0 {
		return "", fmt.Errorf("unable to fetch Bongo board from Puzzmo %w", err)
	}

	return board, err
}
