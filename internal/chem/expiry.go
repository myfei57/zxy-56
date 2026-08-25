package chem

import (
	"time"

	"labvent/internal/store"
)

type ExpiryService struct {
	blobs store.Blob
}

func NewExpiryService(blobs store.Blob) *ExpiryService {
	return &ExpiryService{blobs: blobs}
}

func (s *ExpiryService) Register(chemID string, expiryAt time.Time) error {
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
	return s.blobs.Save("chem-expiry", chemID, ExpiryRecord{ChemID: chemID, EffectiveUntil: effectiveUntil})
}

func (s *ExpiryService) EffectiveDate(chemID string) (time.Time, error) {
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
