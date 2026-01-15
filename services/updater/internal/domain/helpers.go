package domain

import (
	"context"
	"log/slog"
	"regexp"
	"strings"
)

type Noter interface {
	FindFounder(context.Context, Updater, *slog.Logger) (string, error)
	FindAncestor(context.Context, Updater, *slog.Logger) (string, error)
	FindFather(context.Context, Updater, *slog.Logger) (string, error)
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

func ReturnTypesNote(n *Note) Noter {

	switch {

	case strings.HasPrefix(n.Title, "мысль"):

		return &NotePeriodicThought{&NotePeriodic{n}}

	case strings.HasPrefix(n.Title, "сон"):

		return &NotePeriodicDream{&NotePeriodic{n}}

	case strings.HasPrefix(n.Title, "человек"):

		return &NoteHuman{n}

	case IsDay(n.Title):

		return &NotePeriodicDaily{&NotePeriodic{n}}

	case IsWeek(n.Title):

		return &NotePeriodicWeekly{&NotePeriodic{n}}

	case IsMonth(n.Title):

		return &NotePeriodicMonthly{&NotePeriodic{n}}

	case IsQuarter(n.Title):

		return &NotePeriodicQuarterly{&NotePeriodic{n}}

	case IsYear(n.Title):

		return &NotePeriodicYearly{&NotePeriodic{n}}

	default:

		return n

	}
}
