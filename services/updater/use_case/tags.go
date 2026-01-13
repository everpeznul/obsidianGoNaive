package use_case

import "context"

type Tager struct{}

func (upd *UpdaterService) TagsFormat(ctx context.Context, n Noter) ([]string, error) {

	return []string{"tempTag"}, nil
}
