package chem

import (
	"time"

	"labvent/internal/store"
)

type ExpiryService struct {
	blobs store.Blob
	cache map[string]time.Time
}

func NewExpiryService(blobs store.Blob) *ExpiryService {
	return &ExpiryService{blobs: blobs, cache: map[string]time.Time{}}
}

func (s *ExpiryService) Register(chemID string, expiryAt time.Time) error {
	s.cache[chemID] = expiryAt
	item := ExpiryRecord{ChemID: chemID, EffectiveUntil: expiryAt}
	return s.blobs.Save("chem-expiry", chemID, item)
}

func (s *ExpiryService) Revise(chemID string, effectiveUntil time.Time) error {
	var chem Chemical
	if err := s.blobs.Load("chem", chemID, &chem); err != nil {
		return err
	}
	chem.EffectiveUntil = effectiveUntil
	if err := s.blobs.Save("chem", chemID, chem); err != nil {
		return err
	}
	if err := s.blobs.Save("chem-expiry", chemID, ExpiryRecord{ChemID: chemID, EffectiveUntil: effectiveUntil}); err != nil {
		return err
	}
	// Keep the in-memory cache in sync with the persisted record so that
	// EffectiveDate (the source the issue gate reads from) reflects the
	// current expiry. Without this, a downward revision would still be
	// shadowed by the stale cached value and an expired reagent could be
	// issued.
	s.cache[chemID] = effectiveUntil
	return nil
}

func (s *ExpiryService) EffectiveDate(chemID string) (time.Time, error) {
	if date, ok := s.cache[chemID]; ok {
		return date, nil
	}
	var record ExpiryRecord
	if err := s.blobs.Load("chem-expiry", chemID, &record); err != nil {
		return time.Time{}, err
	}
	return record.EffectiveUntil, nil
}

type ExpiryRecord struct {
	ChemID         string    `json:"chem_id"`
	EffectiveUntil time.Time `json:"effective_until"`
}
