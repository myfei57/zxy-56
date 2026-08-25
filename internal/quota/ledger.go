package quota

import (
	"time"

	"labvent/internal/store"
)

type Ledger struct {
	blobs store.Blob
}

func NewLedger(blobs store.Blob) *Ledger {
	return &Ledger{blobs: blobs}
}

func (l *Ledger) path(labID string) string {
	return store.JoinPath(l.blobs.Root(), "quota-ledger", labID+".json")
}

func (l *Ledger) Current(labID string) (float64, error) {
	path := l.path(labID)
	if !l.blobs.Exists("quota-ledger", labID) {
		return 0, nil
	}
	var entry LedgerEntry
	if err := store.ReadJSON(path, &entry); err != nil {
		return 0, err
	}
	if entry.Date != time.Now().Format("2006-01-02") {
		return 0, nil
	}
	return entry.Hours, nil
}

func (l *Ledger) Add(labID string, hours float64) error {
	current, err := l.Current(labID)
	if err != nil {
		return err
	}
	entry := LedgerEntry{LabID: labID, Date: time.Now().Format("2006-01-02"), Hours: current + hours}
	return store.WriteJSON(l.path(labID), entry)
}

func (l *Ledger) Subtract(labID string, hours float64) error {
	current, err := l.Current(labID)
	if err != nil {
		return err
	}
	entry := LedgerEntry{LabID: labID, Date: time.Now().Format("2006-01-02"), Hours: current - hours}
	return store.WriteJSON(l.path(labID), entry)
}

type LedgerEntry struct {
	LabID string  `json:"lab_id"`
	Date  string  `json:"date"`
	Hours float64 `json:"hours"`
}

func (l *Ledger) Reset(labID string) error {
	entry := LedgerEntry{LabID: labID, Date: time.Now().Format("2006-01-02"), Hours: 0}
	return store.WriteJSON(l.path(labID), entry)
}
