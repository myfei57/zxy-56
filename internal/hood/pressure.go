package hood

import (
	"time"

	"labvent/internal/store"
)

type PressureState struct {
	HoodID string    `json:"hood_id"`
	Value  float64   `json:"value"`
	At     time.Time `json:"at"`
}

type PressureService struct {
	blobs store.Blob
}

func NewPressureService(blobs store.Blob) *PressureService {
	return &PressureService{blobs: blobs}
}

func (s *PressureService) Observe(hoodID string, value float64) error {
	state := PressureState{HoodID: hoodID, Value: value, At: time.Now()}
	return s.blobs.Save("pressure", hoodID, state)
}

func (s *PressureService) State(hoodID string) (PressureState, error) {
	var state PressureState
	if err := s.blobs.Load("pressure", hoodID, &state); err != nil {
		return PressureState{}, err
	}
	return state, nil
}

func (s *PressureService) Healthy(hoodID string) (bool, error) {
	state, err := s.State(hoodID)
	if err != nil {
		return false, err
	}
	return state.Value <= -5, nil
}
