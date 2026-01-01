package domain

import (
	"strings"
)

type Note_human struct {
	Note
}
type Note_periodic struct {
	Note
}
type Note_periodic_daily struct {
	Note_periodic
}

func (n *Note_periodic_daily) FindFounder() (string, error) {

	if Exists("0000-00-00") {
	}
	return "0000-00-00", nil
}

type Note_periodic_weekly struct {
	Note_periodic
}

func (n *Note_periodic_weekly) FindFounder() (string, error) {

	if Exists("0000-W00") {
	}
	return "0000-W00", nil
}

type Note_periodic_monthly struct {
	Note_periodic
}

func (n *Note_periodic_monthly) FindFounder() (string, error) {
	if Exists("0000-00") {
	}
	return "0000-00", nil
}

type Note_periodic_quarterly struct {
	Note_periodic
}

func (n *Note_periodic_quarterly) FindFounder() (string, error) {
	if Exists("0000-Q0") {
	}
	return "0000-Q0", nil
}

type Note_periodic_yearly struct {
	Note_periodic
}

func (n *Note_periodic_yearly) FindFounder() (string, error) {
	if Exists("0000") {
	}
	return "0000", nil
}

type Note_periodic_dream struct {
	Note_periodic
}

func (n *Note_periodic_dream) FindAncestor() (string, error) {

	ancestor := strings.Split(n.Title, ".")
	if len(ancestor) != 3 {
		return "", nil
	}
	if Exists(ancestor[1]) {
	}
	return ancestor[1], nil
}

type Note_periodic_thought struct {
	Note_periodic
}

func (n *Note_periodic_thought) FindAncestor() (string, error) {

	ancestor := strings.Split(n.Title, ".")
	if len(ancestor) != 3 {
		return "", nil
	}
	if Exists(ancestor[1]) {
	}
	return ancestor[1], nil
}
