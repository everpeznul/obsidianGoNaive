package http

import (
	"github.com/google/uuid"
	"time"
)

type note struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Title      string    `json:"title"`
	Path       string    `json:"path"`
	Class      string    `json:"class"`
	Tags       []string  `json:"tags"`
	Links      []string  `json:"links"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time,omitempty"`
	UpdateTime time.Time `json:"update_time,omitempty"`
}
