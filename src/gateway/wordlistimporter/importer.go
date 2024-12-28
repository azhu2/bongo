package wordlistimporter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"

	"go.uber.org/fx"
)

const path = "../../../../wordlist/wordlist-20210729.txt"

var wordRegex = regexp.MustCompile(`^\"\w{1,5}\"$`)

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
	filtered := slices.DeleteFunc(rows,
		func(word string) bool { return !wordRegex.MatchString(word) },
	)
	fmt.Printf("found %d words from %s\n", len(filtered), path)
	return filtered, nil
}
