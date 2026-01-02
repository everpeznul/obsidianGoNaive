package domain

import (
	"context"
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

func (n *Note_periodic_daily) FindFounder(ctx context.Context) (string, error) {

	if Exists(ctx, "0000-00-00") {
	}

	obsiLog.Debug("Daily FindFounder", "founder", "0000-00-00")
	return "0000-00-00", nil
}

type Note_periodic_weekly struct {
	Note_periodic
}

func (n *Note_periodic_weekly) FindFounder(ctx context.Context) (string, error) {

	if Exists(ctx, "0000-W00") {
	}

	obsiLog.Debug("Weekly FindFounder", "founder", "0000-W00")
	return "0000-W00", nil
}

type Note_periodic_monthly struct {
	Note_periodic
}

func (n *Note_periodic_monthly) FindFounder(ctx context.Context) (string, error) {
	if Exists(ctx, "0000-00") {
	}

	obsiLog.Debug("Monthly FindFounder", "founder", "0000-00")
	return "0000-00", nil
}

type Note_periodic_quarterly struct {
	Note_periodic
}

func (n *Note_periodic_quarterly) FindFounder(ctx context.Context) (string, error) {
	if Exists(ctx, "0000-Q0") {
	}

	obsiLog.Debug("Quarterly FindFounder", "founder", "0000-Q0")
	return "0000-Q0", nil
}

type Note_periodic_yearly struct {
	Note_periodic
}

func (n *Note_periodic_yearly) FindFounder(ctx context.Context) (string, error) {
	if Exists(ctx, "0000") {
	}

	obsiLog.Debug("Yearly FindFounder", "founder", "0000")
	return "0000", nil
}

type Note_periodic_dream struct {
	Note_periodic
}

func (n *Note_periodic_dream) FindAncestor(ctx context.Context) (string, error) {

	ancestor := strings.Split(n.Title, ".")
	if len(ancestor) != 3 {
		return "", nil
	}
	if Exists(ctx, ancestor[1]) {
	}

	obsiLog.Debug("Dream FindAncestor", "ancestor", ancestor[1])
	return ancestor[1], nil
}

type Note_periodic_thought struct {
	Note_periodic
}

func (n *Note_periodic_thought) FindAncestor(ctx context.Context) (string, error) {

	ancestor := strings.Split(n.Title, ".")
	if len(ancestor) != 3 {
		return "", nil
	}
	if Exists(ctx, ancestor[1]) {
	}

	obsiLog.Debug("Thought", "ancestor", ancestor[1])
	return ancestor[1], nil
}
