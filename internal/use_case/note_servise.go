package use_case

import "context"

type Noter interface {
	FindFather(context.Context) (string, error)
	FindAncestor(context.Context) (string, error)
	FindFounder(context.Context) (string, error)
}
