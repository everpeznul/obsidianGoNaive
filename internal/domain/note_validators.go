package domain

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

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

	return id, err
}
