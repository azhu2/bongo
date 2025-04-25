package gameimporter

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/machinebox/graphql"
	"go.uber.org/fx"
)

const (
	GraphqlEndpoint = "https://www.puzzmo.com/_api/prod/graphql"
	gameKey         = "today:/%s/bongo"
)

var GraphqlModule = fx.Module("graphqlimporter",
	fx.Provide(NewGraphql),
)

type graphqlGateway struct {
	userID    string
	authToken string

	graphqlClient *graphql.Client
}

func NewGraphql(p Params) (Result, error) {
	return Result{
		Gateway: &graphqlGateway{
			userID:    p.Secrets.UserID,
			authToken: p.Secrets.AuthToken,

			graphqlClient: p.GraphqlClient,
		},
	}, nil
}

func (g *graphqlGateway) ImportBoard(ctx context.Context, date string) (string, error) {
	return g.importBoardFromDailyScreen(ctx, date)
}

// (legacy) importing directly from game screen. Assumes game slug
func (g *graphqlGateway) importBoardFromGameScreen(ctx context.Context, date string) (string, error) {
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
	req.Var("finderKey", fmt.Sprintf(gameKey, date))
	req.Var("gameContext", map[string]any{
		"partnerSlug":             nil,
		"pingOwnerForMultiplayer": true,
	})
	req.Header.Set("context-type", "application/json")
	req.Header.Set("authorization", g.authToken)
	req.Header.Set("auth-provider", "custom")
	req.Header.Set("puzzmo-gameplay-id", g.userID)

	var resp graphqlBoardResponse
	err := g.graphqlClient.Run(ctx, req, &resp)
	board := resp.StartOrFindGameplay.GamePlayed.Puzzle.Puzzle
	if err != nil || len(board) == 0 {
		return "", fmt.Errorf("unable to fetch Bongo board from Puzzmo %w", err)
	}

	slog.Debug("loaded game board",
		"source", "graphql",
		"date", date,
	)

	return board, err
}

// Import board from daily game screen. Grabs the first game with "bongo" in its slug
func (g *graphqlGateway) importBoardFromDailyScreen(ctx context.Context, date string) (string, error) {
	req := graphql.NewRequest(`
		query TodayScreenQuery(
			$day: String
		) {
			todayPage(day: $day) {
				daily {
					puzzles {
						urlPath
						puzzle {
							game {
								slug
							}
							puzzle
						}
					}
				}
			}
		}
	`)
	req.Var("day", date)
	req.Header.Set("context-type", "application/json")
	req.Header.Set("authorization", g.authToken)
	req.Header.Set("auth-provider", "custom")
	req.Header.Set("puzzmo-gameplay-id", g.userID)

	var resp graphqlTodayScreenResponse
	err := g.graphqlClient.Run(ctx, req, &resp)

	if err != nil {
		return "", fmt.Errorf("unable to fetch daily puzzles from Puzzmo %w", err)
	}

	for _, puzzle := range resp.TodayPage.Daily.Puzzles {
		if strings.Contains(puzzle.Puzzle.Game.Slug, "bongo") {
			board := puzzle.Puzzle.Puzzle
			if len(board) > 0 {
				slog.Debug("loaded game board",
					"source", "graphql",
					"slug", puzzle.Puzzle.Game.Slug,
				)
				return board, nil
			}
			slog.Warn("found empty board",
				"slug", puzzle.Puzzle.Game.Slug,
			)
		}
	}

	return "", fmt.Errorf("no bongo boards found")
}
