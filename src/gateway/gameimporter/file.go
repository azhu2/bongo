package gameimporter

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"go.uber.org/fx"
)

const (
	fileFormat = "../../../../testdata/%s.txt"
)

var FileModule = fx.Module("importer",
	fx.Provide(NewFile),
)

type fileImporter struct{}

func NewFile(p Params) (Result, error) {
	return Result{
		Gateway: &fileImporter{},
	}, nil
}

func (f *fileImporter) ImportBoard(ctx context.Context, date string) (string, error) {
	_, file, _, _ := runtime.Caller(0)
	path := filepath.Join(file, fmt.Sprintf(fileFormat, date))
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	slog.Debug("loaded game board",
		"source", "file",
		"path", path,
	)
	return string(raw), nil
}
