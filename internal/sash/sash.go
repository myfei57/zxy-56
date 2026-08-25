package sash

import (
	"fmt"

	"github.com/google/uuid"

	"labvent/internal/store"
)

type Sash struct {
	ID       string `json:"id"`
	HoodID   string `json:"hood_id"`
	Height   int    `json:"height"`
	FullOpen int    `json:"full_open"`
}

type SashService struct {
	blobs store.Blob
}

func NewSashService(blobs store.Blob) *SashService {
	return &SashService{blobs: blobs}
}

func (s *SashService) Register(hoodID string, fullOpen int) (Sash, error) {
	if hoodID == "" || fullOpen <= 0 {
		return Sash{}, fmt.Errorf("sash hood and full-open height are required")
	}
	item := Sash{ID: uuid.NewString(), HoodID: hoodID, Height: fullOpen, FullOpen: fullOpen}
	if err := s.blobs.Save("sash", item.ID, item); err != nil {
		return Sash{}, err
	}
	return item, nil
}

func (s *SashService) Get(id string) (Sash, error) {
	var item Sash
	if err := s.blobs.Load("sash", id, &item); err != nil {
		return Sash{}, err
	}
	return item, nil
}

func (s *SashService) SetHeight(id string, height int) error {
	if height <= 0 {
		return fmt.Errorf("sash height must be positive")
	}
	item, err := s.Get(id)
	if err != nil {
		return err
	}
	item.Height = height
	return s.blobs.Save("sash", item.ID, item)
}
