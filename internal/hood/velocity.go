package hood

import (
	"fmt"

	"labvent/internal/store"
)

type VelocityState struct {
	HoodID    string  `json:"hood_id"`
	Threshold float64 `json:"threshold"`
}

type VelocityService struct {
	blobs store.Blob
}

func NewVelocityService(blobs store.Blob) *VelocityService {
	return &VelocityService{blobs: blobs}
}

func (s *VelocityService) SetFaceThreshold(hoodID string, threshold float64) error {
	if threshold <= 0 {
		return fmt.Errorf("face velocity threshold must be positive")
	}
	return s.blobs.Save("velocity", hoodID, VelocityState{HoodID: hoodID, Threshold: fullOpenThreshold})
}

func (s *VelocityService) Threshold(hoodID string) (float64, error) {
	var state VelocityState
	if err := s.blobs.Load("velocity", hoodID, &state); err != nil {
		return 0, err
	}
	return state.Threshold, nil
}

func (s *VelocityService) Verdict(hoodID string, measured float64) error {
	threshold, err := s.Threshold(hoodID)
	if err != nil {
		return err
	}
	if measured > threshold {
		return fmt.Errorf("face velocity %.2f above threshold %.2f", measured, threshold)
	}
	return nil
}

const fullOpenThreshold = 0.5
