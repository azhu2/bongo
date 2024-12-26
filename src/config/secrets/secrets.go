package secrets

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

const (
	envAuthToken = "AUTH_TOKEN"
	envUserID    = "USER_ID"
)

type Secrets struct {
	UserID    string
	AuthToken string
}

var Module = fx.Module("secrets",
	fx.Provide(New),
)

type Result struct {
	fx.Out

	Secrets
}

func New() (Result, error) {
	godotenv.Load("../.env")
	userID := os.Getenv(envUserID)
	if len(userID) == 0 {
		return Result{}, fmt.Errorf("user ID not set in env")
	}
	authToken := os.Getenv(envAuthToken)
	if len(authToken) == 0 {
		return Result{}, fmt.Errorf("auth token not set in env")
	}
	return Result{
		Secrets: Secrets{
			UserID:    userID,
			AuthToken: authToken,
		},
	}, nil
}
