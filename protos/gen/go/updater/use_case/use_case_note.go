package use_case

import (
	"context"
	"fmt"
	"regexp"
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

func Exists(ctx context.Context, title string) (bool, error) {

	_, err := Repo.FindByName(ctx, title)

	if err != nil {

		return false, fmt.Errorf("note=%s not found: %w", title, err)
	}

	return true, nil
}

func Create(ctx context.Context, title string) (uuid.UUID, error) {

	id, err := Repo.Insert(ctx, Note{})

	if err != nil {

		return uuid.Nil, fmt.Errorf("create ERROR: %w", err)
	}

	return id, nil
}

var (
	ReYear    = regexp.MustCompile(`^\d{4}$`)
	ReMonth   = regexp.MustCompile(`^(?P<year>\d{4})-(?P<month>0[0-9]|1[0-2])$`)
	ReQuarter = regexp.MustCompile(`^(?P<year>\d{4})-Q(?P<quarter>[0-4])$`)
	ReWeek    = regexp.MustCompile(`^(?P<year>\d{4})-W(?P<week>[0-4]\d|5[0-3])$`)
	ReDay     = regexp.MustCompile(`^(?P<year>\d{4})-(?P<month>0[0-9]|1[0-2])-(?P<day>0[0-9]|[12]\d|3[01])$`)
)

func IsYear(s string) bool    { return ReYear.MatchString(s) }
func IsMonth(s string) bool   { return ReMonth.MatchString(s) }
func IsQuarter(s string) bool { return ReQuarter.MatchString(s) }
func IsWeek(s string) bool    { return ReWeek.MatchString(s) }
func IsDay(s string) bool     { return ReDay.MatchString(s) }

func ReturnTypesNote(n Note) Noter {

	switch {

	case strings.HasPrefix(n.Title, "мысль"):

		return &NotePeriodicThought{NotePeriodic{n}}

	case strings.HasPrefix(n.Title, "сон"):

		return &NotePeriodicDream{NotePeriodic{n}}

	case strings.HasPrefix(n.Title, "человек"):

		return &NoteHuman{n}

	case IsDay(n.Title):

		return &NotePeriodicDaily{NotePeriodic{n}}

	case IsWeek(n.Title):

		return &NotePeriodicWeekly{NotePeriodic{n}}

	case IsMonth(n.Title):

		return &NotePeriodicMonthly{NotePeriodic{n}}

	case IsQuarter(n.Title):

		return &NotePeriodicQuarterly{NotePeriodic{n}}

	case IsYear(n.Title):

		return &NotePeriodicYearly{NotePeriodic{n}}

	default:

		return &n

	}
}
