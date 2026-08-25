package hood

import (
	"time"

	"labvent/internal/store"
)

type Baseline struct {
	RoomID    string    `json:"room_id"`
	Value     float64   `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BaselineService struct {
	blobs store.Blob
}

func NewBaselineService(blobs store.Blob) *BaselineService {
	return &BaselineService{blobs: blobs}
}

func (s *BaselineService) Set(roomID string, value float64) error {
	baseline := Baseline{RoomID: roomID, Value: value, UpdatedAt: time.Now()}
	return s.blobs.Save("baseline", roomID, baseline)
}

func (s *BaselineService) Current(roomID string) (float64, error) {
	var baseline Baseline
	if err := s.blobs.Load("baseline", roomID, &baseline); err != nil {
		return 0, err
	}
	return baseline.Value, nil
}

func (s *BaselineService) Get(roomID string) (Baseline, error) {
	var baseline Baseline
	if err := s.blobs.Load("baseline", roomID, &baseline); err != nil {
		return Baseline{}, err
	}
	return baseline, nil
}
