package wordlistimporter

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"go.uber.org/fx"
)

const path = "../../../../wordlist/wordlist-20210729.txt"

<<<<<<< HEAD
var wordRegex = regexp.MustCompile(`^\"(\w{1,5})\"$`)
=======
var wordRegex = regexp.MustCompile(`^(\w{3,5})$`)
>>>>>>> ad9ae6d (One more try to handle wildcards in bonus words)

var Module = fx.Module("wordimporter",
	fx.Provide(New),
)

type Gateway interface {
	ImportWordList(ctx context.Context) ([]string, error)
}

type Result struct {
	fx.Out

	Gateway
}

type gateway struct {
}

func New() (Result, error) {
	return Result{
		Gateway: &gateway{},
	}, nil
}

func (g *gateway) ImportWordList(ctx context.Context) ([]string, error) {
	_, file, _, _ := runtime.Caller(0)
	path := filepath.Join(file, path)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rows := strings.Split(strings.ToUpper(string(raw)), "\n")
	filtered := []string{}
	for _, word := range rows {
		if match := wordRegex.FindStringSubmatch(word); len(match) > 0 {
			filtered = append(filtered, match[1])
		}
	}
	slog.Debug("loaded word list",
		"path", path,
		"word_count", len(filtered),
	)
	return filtered, nil
}
