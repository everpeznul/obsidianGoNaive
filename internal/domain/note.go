package domain

import (
	"github.com/google/uuid"
	"time"
)

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
