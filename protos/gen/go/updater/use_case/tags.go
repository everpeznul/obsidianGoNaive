package use_case

import "context"

type Tager struct{}

func (t *Tager) Format(ctx context.Context, n Noter) ([]string, error) {

	return []string{"tempTag"}, nil
}
