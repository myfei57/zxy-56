package store

import (
	"fmt"
	"time"
)

type IssueGate struct {
	effective func(chemID string) (time.Time, error)
}

func NewIssueGate(effective func(chemID string) (time.Time, error)) *IssueGate {
	return &IssueGate{effective: effective}
}

func (g *IssueGate) Check(chemID string, now time.Time) error {
	date, err := g.effective(chemID)
	if err != nil {
		return err
	}
	if !now.Before(date) {
		return fmt.Errorf("reagent %s expired at %s", chemID, date.Format(time.RFC3339))
	}
	return nil
}

func (g *IssueGate) Issue(chemID string, quantity float64, now time.Time) error {
	if quantity <= 0 {
		return fmt.Errorf("issue quantity must be positive")
	}
	return g.Check(chemID, now)
}
