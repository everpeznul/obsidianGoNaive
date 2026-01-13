package use_case

import (
	"context"
	"fmt"
)

func (upd UpdaterService) LinksFormat(ctx context.Context, n Noter) ([]string, error) {

	founder, err := n.FindFounder(ctx, upd)
	if err != nil {

		return nil, fmt.Errorf("format FindFounder ERROR: %w", err)
	}

	ancestor, err := n.FindAncestor(ctx, upd)
	if err != nil {

		return nil, fmt.Errorf("format FindAncestor ERROR: %w", err)
	}

	father, err := n.FindFather(ctx, upd)
	if err != nil {

		return nil, fmt.Errorf("format FindFather ERROR: %w", err)
	}

	return []string{founder, ancestor, father}, nil
}
