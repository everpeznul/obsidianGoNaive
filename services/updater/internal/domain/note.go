package domain

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Updater interface {
	Exist(context.Context, string) (bool, error)
	Create(context.Context, string) (uuid.UUID, error)
}

type Notes struct {
	Log *slog.Logger
}

func NewNotes(log *slog.Logger) *Notes {

	return &Notes{Log: log}
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
func (n *Note) FindFounder(ctx context.Context, upd Updater, obsiLog *slog.Logger) (string, error) {

	founder := strings.Split(n.Title, ".")[0]

	ok, err := upd.Exist(ctx, founder)
	if !ok {

		obsiLog.Error("Note FindFounder ERROR", "error", err)

		id, err := upd.Create(ctx, founder)
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
func (n *Note) FindAncestor(ctx context.Context, upd Updater, obsiLog *slog.Logger) (string, error) {

	a := strings.Split(n.Title, ".")
	ancestor := ""

	if len(a) == 1 {
		return n.FindFounder(ctx, upd, obsiLog)
	}
	for i := len(a) - 2; i >= 0; i-- {

		if !strings.Contains(a[i], "%") {

			ancestor = strings.Join(a[:i+1], ".")
			break
		}
	}

	ok, err := upd.Exist(ctx, ancestor)
	if !ok {

		obsiLog.Error("Note FindAncestor ERROR", "error", err)

		id, err := upd.Create(ctx, ancestor)
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
func (n *Note) FindFather(ctx context.Context, upd Updater, obsiLog *slog.Logger) (string, error) {

	f := strings.Split(n.Title, ".")
	if len(f) == 1 {
		return n.FindFounder(ctx, upd, obsiLog)
	}

	father := strings.Join(f[:len(f)-1], ".")

	ok, err := upd.Exist(ctx, father)
	if !ok {

		obsiLog.Error("Note FindFather ERROR", "error", err)

		id, err := upd.Create(ctx, father)
		if err != nil {

			obsiLog.Error("Note FindFather ERROR", "error", err)
			return "", fmt.Errorf("note FindFather ERROR: %w", err)
		}

		obsiLog.Info("Note FindFather Create note", "id", id)
	}

	obsiLog.Debug("Note FindFather", "title", father)
	return father, nil
}
