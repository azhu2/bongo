package importer

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	"go.uber.org/fx"
)

const (
	GraphqlEndpoint = "https://www.puzzmo.com/_api/prod/graphql"
	bongoSlug       = "today:/%s/bongo"

	// Env variables
	envAuthToken = "AUTH_TOKEN"
	envUserID    = "USER_ID"
)

var GraphqlModule = fx.Module("graphqlimporter",
	fx.Provide(NewGraphql),
)

type graphqlGateway struct {
	graphqlClient *graphql.Client
}

func NewGraphql(p Params) (Result, error) {
	godotenv.Load()
	return Result{
		Gateway: &graphqlGateway{
			graphqlClient: p.GraphqlClient,
		},
	}, nil
}

func (g *graphqlGateway) GetBongoBoard(ctx context.Context, date string) (string, error) {
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
	req.Var("finderKey", fmt.Sprintf(bongoSlug, date))
	req.Var("gameContext", map[string]any{
		"partnerSlug":             nil,
		"pingOwnerForMultiplayer": true,
	})
	req.Header.Set("context-type", "application/json")
	req.Header.Set("authorization", os.Getenv(envAuthToken))
	req.Header.Set("auth-provider", "custom")
	req.Header.Set("puzzmo-gameplay-id", os.Getenv(envUserID))

	var resp graphqlBoardResponse
	err := g.graphqlClient.Run(ctx, req, &resp)
	board := resp.StartOrFindGameplay.GamePlayed.Puzzle.Puzzle
	if err != nil || len(board) == 0 {
		return "", fmt.Errorf("unable to fetch Bongo board from Puzzmo %w", err)
	}

	return board, err
}
