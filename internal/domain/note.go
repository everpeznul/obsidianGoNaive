package domain

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
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

func (n *Note) FindFounder(ctx context.Context) (string, error) {

	founder := strings.Split(n.Title, ".")[0]
	if Exists(ctx, founder) {

	}

	obsiLog.Debug("Note FindFounder", "founder", founder)

	return founder, nil
}
func (n *Note) FindAncestor(ctx context.Context) (string, error) {

	a := strings.Split(n.Title, ".")
	ancestor := ""

	if len(a) == 1 {
		return a[0], nil
	}
	for i := len(a) - 2; i >= 0; i-- {

		if !strings.Contains(a[i], "%") {
			ancestor = strings.Join(a[:i+1], ".")
			break
		}
	}

	if Exists(ctx, ancestor) {
	}

	obsiLog.Debug("Note FindAncestor", "ancestor", ancestor)
	return ancestor, nil
}
func (n *Note) FindFather(ctx context.Context) (string, error) {

	f := strings.Split(n.Title, ".")
	if len(f) == 1 {
		return f[0], nil
	}
	father := strings.Join(f[:len(f)-1], ".")
	if Exists(ctx, father) {

	}

	obsiLog.Debug("Note FindFather", "father", father)
	return father, nil
}
