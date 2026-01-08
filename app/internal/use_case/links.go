package use_case

import (
	"context"
	"fmt"
)

type Linker struct {
}

func (l *Linker) Format(ctx context.Context, n Noter) ([]string, error) {

	founder, err := n.FindFounder(ctx)
	if err != nil {

		return nil, fmt.Errorf("format FindFounder ERROR: %w", err)
	}

	ancestor, err := n.FindAncestor(ctx)
	if err != nil {

		return nil, fmt.Errorf("format FindAncestor ERROR: %w", err)
	}

	father, err := n.FindFather(ctx)
	if err != nil {

		return nil, fmt.Errorf("format FindFather ERROR: %w", err)
	}

	return []string{founder, ancestor, father}, nil
}
