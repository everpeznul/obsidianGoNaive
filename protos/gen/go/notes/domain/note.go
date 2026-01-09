package domain

import (
	"context"
	"fmt"
	"strings"
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

// FindFounder поиск "основателя" ветки заметок
func (n *Note) FindFounder(ctx context.Context) (string, error) {

	founder := strings.Split(n.Title, ".")[0]

	ok, err := Exists(ctx, founder)
	if !ok {

		obsiLog.Error("Note FindFounder ERROR", "error", err)

		id, err := Create(ctx, founder)
		if err != nil {

			obsiLog.Error("Note FindFounder ERROR", "error", err)
			return "", fmt.Errorf("note FindFounder ERROR: %w", err)
		}

		obsiLog.Info("Note FindFounder Create note", "id", id)
	}

	obsiLog.Debug("Note FindFounder", "title", founder)
	return founder, nil
}

// FindAncestor поиск ближайшего "предка" заметки в ветке
func (n *Note) FindAncestor(ctx context.Context) (string, error) {

	a := strings.Split(n.Title, ".")
	ancestor := ""

	if len(a) == 1 {
		return n.FindFounder(ctx)
	}
	for i := len(a) - 2; i >= 0; i-- {

		if !strings.Contains(a[i], "%") {

			ancestor = strings.Join(a[:i+1], ".")
			break
		}
	}

	ok, err := Exists(ctx, ancestor)
	if !ok {

		obsiLog.Error("Note FindAncestor ERROR", "error", err)

		id, err := Create(ctx, ancestor)
		if err != nil {

			obsiLog.Error("Note FindAncestor ERROR", "error", err)
			return "", fmt.Errorf("note FindAncestor ERROR: %w", err)
		}

		obsiLog.Info("Note FindAncestor Create note", "id", id)
	}

	obsiLog.Debug("Note FindAncestor", "title", ancestor)
	return ancestor, nil
}

// FindFather поиск  "отца" заметки в ветке
func (n *Note) FindFather(ctx context.Context) (string, error) {

	f := strings.Split(n.Title, ".")
	if len(f) == 1 {
		return n.FindFounder(ctx)
	}

	father := strings.Join(f[:len(f)-1], ".")

	ok, err := Exists(ctx, father)
	if !ok {

		obsiLog.Error("Note FindFather ERROR", "error", err)

		id, err := Create(ctx, father)
		if err != nil {

			obsiLog.Error("Note FindFather ERROR", "error", err)
			return "", fmt.Errorf("note FindFather ERROR: %w", err)
		}

		obsiLog.Info("Note FindFather Create note", "id", id)
	}

	obsiLog.Debug("Note FindFather", "title", father)
	return father, nil
}
