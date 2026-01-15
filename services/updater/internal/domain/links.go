package domain

import (
	"context"
	"fmt"
	"log/slog"
)

func LinksFormat(ctx context.Context, n Noter, u Updater, obsiLog *slog.Logger) ([]string, error) {

	founder, err := n.FindFounder(ctx, u, obsiLog)
	if err != nil {

		return nil, fmt.Errorf("format FindFounder ERROR: %w", err)
	}

	ancestor, err := n.FindAncestor(ctx, u, obsiLog)
	if err != nil {

		return nil, fmt.Errorf("format FindAncestor ERROR: %w", err)
	}

	father, err := n.FindFather(ctx, u, obsiLog)
	if err != nil {

		return nil, fmt.Errorf("format FindFather ERROR: %w", err)
	}

	return []string{founder, ancestor, father}, nil
}
