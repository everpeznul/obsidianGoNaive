package use_case

import "context"

type Noter interface {
	FindFather(context.Context, UpdaterService) (string, error)
	FindAncestor(context.Context, UpdaterService) (string, error)
	FindFounder(context.Context, UpdaterService) (string, error)
}
