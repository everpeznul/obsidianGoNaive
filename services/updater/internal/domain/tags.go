package domain

import (
	"context"
)

type Tager struct{}

func TagsFormat(ctx context.Context, n Noter, u Updater) ([]string, error) {

	return []string{"tempTag"}, nil
}
