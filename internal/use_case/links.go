package use_case

import "context"

type Linker struct {
}

func (l *Linker) Format(ctx context.Context, n Noter) []string {

	founder, _ := n.FindFounder(ctx)
	ancestor, _ := n.FindAncestor(ctx)
	father, _ := n.FindFather(ctx)

	return []string{founder, ancestor, father}
}
