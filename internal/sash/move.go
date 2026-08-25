package sash

import (
	"fmt"
	"time"

	"labvent/internal/store"
)

type MoveRecord struct {
	SashID string    `json:"sash_id"`
	Height int       `json:"height"`
	At     time.Time `json:"at"`
}

type MoveService struct {
	blobs      store.Blob
	airflow    func(hoodID string, flow int) error
	notifyFlow func(hoodID string, threshold float64) error
}

func NewMoveService(
	blobs store.Blob,
	airflow func(hoodID string, flow int) error,
	notifyFlow func(hoodID string, threshold float64) error,
) *MoveService {
	return &MoveService{blobs: blobs, airflow: airflow, notifyFlow: notifyFlow}
}

func (s *MoveService) Move(sashID string, height int) error {
	var item Sash
	if err := s.blobs.Load("sash", sashID, &item); err != nil {
		return err
	}
	if height <= 0 || height > item.FullOpen {
		return fmt.Errorf("sash height out of range")
	}
	record := MoveRecord{SashID: sashID, Height: height, At: time.Now()}
	if err := s.blobs.Save("sash-move", sashID, record); err != nil {
		return err
	}
	flow := flowFor(item, height)
	if err := s.airflow(item.HoodID, flow); err != nil {
		return err
	}
	threshold := ThresholdFor(item.FullOpen, height)
	if err := s.notifyFlow(item.HoodID, threshold); err != nil {
		return err
	}
	item.Height = height
	return s.blobs.Save("sash", sashID, item)
}

func flowFor(item Sash, height int) int {
	return height * 40
}
