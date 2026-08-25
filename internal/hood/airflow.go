package hood

import (
	"time"

	"labvent/internal/store"
)

type AirflowReading struct {
	HoodID string    `json:"hood_id"`
	Flow   int       `json:"flow"`
	At     time.Time `json:"at"`
}

type AirflowService struct {
	blobs store.Blob
}

func NewAirflowService(blobs store.Blob) *AirflowService {
	return &AirflowService{blobs: blobs}
}

func (s *AirflowService) Persist(hoodID string, flow int) error {
	reading := AirflowReading{HoodID: hoodID, Flow: flow, At: time.Now()}
	return s.blobs.Save("airflow", hoodID, reading)
}

func (s *AirflowService) Reading(hoodID string) (AirflowReading, error) {
	var reading AirflowReading
	if err := s.blobs.Load("airflow", hoodID, &reading); err != nil {
		return AirflowReading{}, err
	}
	return reading, nil
}

func (s *AirflowService) Makeup(hoodID string) (int, error) {
	reading, err := s.Reading(hoodID)
	if err != nil {
		return 0, err
	}
	return reading.Flow + 120, nil
}

func (s *AirflowService) Latest(hoodID string) (AirflowReading, error) {
	return s.Reading(hoodID)
}
