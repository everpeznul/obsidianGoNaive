package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Noter interface {
	FindFather(context.Context) (string, error)
	FindAncestor(context.Context) (string, error)
	FindFounder(context.Context) (string, error)
}

type Note struct {
	Id         uuid.UUID
	Title      string
	Path       string
	Class      string
	Tags       []string
	Links      []string
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}
