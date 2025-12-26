package domain

import "strings"
import "regexp"

type Noter interface {
	FindFather() string
	FindAncestor() string
	FindFounder() string
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

type Has struct {
}

func ReturnTypesNote(n Note) Noter {

	switch {

	case strings.HasPrefix(n.Title, "мысль"):
		return &Note_periodic_thought{Note_periodic{n}}

	case strings.HasPrefix(n.Title, "сон"):
		return &Note_periodic_dream{Note_periodic{n}}

	case strings.HasPrefix(n.Title, "человек"):
		return &Note_human{n}

	case IsDay(n.Title):
		return &Note_periodic_daily{Note_periodic{n}}

	case IsWeek(n.Title):
		return &Note_periodic_weekly{Note_periodic{n}}

	case IsMonth(n.Title):
		return &Note_periodic_monthly{Note_periodic{n}}

	case IsQuarter(n.Title):
		return &Note_periodic_quarterly{Note_periodic{n}}

	case IsYear(n.Title):
		return &Note_periodic_yearly{Note_periodic{n}}

	default:
		return &n
	}
}

func Exists(title string) bool {
	_, err := Repo.FindByName(title)

	if err != nil {
		//если заметка не существует, то создать заметку и обновить её содержимое
	}
	/*
		if len(note) > 1 {
		//такого быть не должно, нужно разрешать

		    return true
		}
	*/
	return true
}
