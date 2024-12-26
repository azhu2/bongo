package importer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/fx"
)

const (
	fileFormat = "testdata/%s.txt"
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

func (f *fileImporter) GetBongoBoard(ctx context.Context, date string) (string, error) {
	base, _ := os.Getwd()
	path := filepath.Join(base, fmt.Sprintf(fileFormat, date))
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
