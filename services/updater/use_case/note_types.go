package use_case

import (
	"context"
	"fmt"
	"strings"
)

type NoteHuman struct {
	Note
}
type NotePeriodic struct {
	Note
}
type NotePeriodicDaily struct {
	NotePeriodic
}

func (n *NotePeriodicDaily) FindFounder(ctx context.Context, upd UpdaterService) (string, error) {

	ok, err := upd.Exists(ctx, "0000-00-00")
	if !ok {

		obsiLog.Error("Daily FindFounder ERROR", "error", err)
		return "", fmt.Errorf("daily FindFounder not find: %w", err)
	}

	obsiLog.Debug("Daily FindFounder", "title", "0000-00-00")
	return "0000-00-00", nil
}

type NotePeriodicWeekly struct {
	NotePeriodic
}

func (n *NotePeriodicWeekly) FindFounder(ctx context.Context, upd UpdaterService) (string, error) {

	ok, err := upd.Exists(ctx, "0000-W00")
	if !ok {

		obsiLog.Error("Weekly FindFounder ERROR", "error", err)
		return "", fmt.Errorf("weekly FindFounder not find: %w", err)
	}

	obsiLog.Debug("Weekly FindFounder", "title", "0000-W00")
	return "0000-W00", nil
}

type NotePeriodicMonthly struct {
	NotePeriodic
}

func (n *NotePeriodicMonthly) FindFounder(ctx context.Context, upd UpdaterService) (string, error) {

	ok, err := upd.Exists(ctx, "0000-00")
	if !ok {

		obsiLog.Error("Monthly FindFounder ERROR", "error", err)
		return "", fmt.Errorf("monthly FindFounder not find: %w", err)
	}

	obsiLog.Debug("Monthly FindFounder", "title", "0000-00")
	return "0000-00", nil
}

type NotePeriodicQuarterly struct {
	NotePeriodic
}

func (n *NotePeriodicQuarterly) FindFounder(ctx context.Context, upd UpdaterService) (string, error) {

	ok, err := upd.Exists(ctx, "0000-Q0")
	if !ok {

		obsiLog.Error("Quarterly FindFounder ERROR", "error", err)
		return "", fmt.Errorf("quarterly FindFounder not find: %w", err)
	}

	obsiLog.Debug("Quarterly FindFounder", "title", "0000-Q0")
	return "0000-Q0", nil
}

type NotePeriodicYearly struct {
	NotePeriodic
}

func (n *NotePeriodicYearly) FindFounder(ctx context.Context, upd UpdaterService) (string, error) {

	ok, err := upd.Exists(ctx, "0000")
	if !ok {

		obsiLog.Error("Yearly FindFounder ERROR", "error", err)
		return "", fmt.Errorf("yearly FindFounder not find: %w", err)
	}

	obsiLog.Debug("Yearly FindFounder", "title", "0000")
	return "0000", nil
}

type NotePeriodicDream struct {
	NotePeriodic
}

func (n *NotePeriodicDream) FindAncestor(ctx context.Context, upd UpdaterService) (string, error) {

	a := strings.Split(n.Title, ".")
	if len(a) != 3 {
		return "", fmt.Errorf("dream not valid")
	}
	ancestor := a[1]

	ok, err := upd.Exists(ctx, ancestor)
	if !ok {

		obsiLog.Error("Dream FindAncestor ERROR", "error", err)

		id, err := upd.Create(ctx, ancestor)
		if err != nil {

			obsiLog.Error("Dream FindAncestor ERROR", "error", err)
			return "", fmt.Errorf("dream FindAncestor ERROR: %w", err)
		}

		obsiLog.Info("Dream FindAncestor Create note", "id", id)
	}

	obsiLog.Debug("Dream FindAncestor", "title", ancestor)
	return ancestor, nil
}

type NotePeriodicThought struct {
	NotePeriodic
}

func (n *NotePeriodicThought) FindAncestor(ctx context.Context, upd UpdaterService) (string, error) {

	a := strings.Split(n.Title, ".")
	if len(a) != 3 {
		return "", fmt.Errorf("thought not valid")
	}
	ancestor := a[1]

	ok, err := upd.Exists(ctx, ancestor)
	if !ok {

		obsiLog.Error("Thought FindAncestor ERROR", "error", err)

		id, err := upd.Create(ctx, ancestor)
		if err != nil {

			obsiLog.Error("Thought FindAncestor ERROR", "error", err)
			return "", fmt.Errorf("thought FindAncestor ERROR: %w", err)
		}

		obsiLog.Info("Thought FindAncestor Create note", "id", id)
	}

	obsiLog.Debug("Thought FindAncestor", "title", ancestor)
	return ancestor, nil
}
