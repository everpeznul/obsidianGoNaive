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

func (n *Note_periodic_daily) FindFounder() string {

	if Exists("0000-00-00") {
	}
	return "0000-00-00"
}

type Note_periodic_weekly struct {
	Note_periodic
}

func (n *Note_periodic_weekly) FindFounder() string {

	if Exists("0000-W00") {
	}
	return "0000-W00"
}

type Note_periodic_monthly struct {
	Note_periodic
}

func (n *Note_periodic_monthly) FindFounder() string {
	if Exists("0000-00") {
	}
	return "0000-00"
}

type Note_periodic_quarterly struct {
	Note_periodic
}

func (n *Note_periodic_quarterly) FindFounder() string {
	if Exists("0000-Q0") {
	}
	return "0000-Q0"
}

type Note_periodic_yearly struct {
	Note_periodic
}

func (n *Note_periodic_yearly) FindFounder() string {
	if Exists("0000") {
	}
	return "0000"
}

type Note_periodic_dream struct {
	Note_periodic
}

func (n *Note_periodic_dream) FindAncestor() string {

	ancestor := strings.Split(n.Title, ".")[1]
	if Exists(ancestor) {
	}
	return ancestor
}

type Note_periodic_thought struct {
	Note_periodic
}

func (n *Note_periodic_thought) FindAncestor() string {

	ancestor := strings.Split(n.Title, ".")[1]
	if Exists(ancestor) {
	}
	return ancestor
}
