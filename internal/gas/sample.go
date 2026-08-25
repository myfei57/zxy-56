package gas

import (
	"time"

	"labvent/internal/store"
)

type Reading struct {
	SensorID string    `json:"sensor_id"`
	PPM      float64   `json:"ppm"`
	At       time.Time `json:"at"`
}

type SampleService struct {
	blobs store.Blob
}

func NewSampleService(blobs store.Blob) *SampleService {
	return &SampleService{blobs: blobs}
}

func (s *SampleService) Record(sensorID string, ppm float64) error {
	reading := Reading{SensorID: sensorID, PPM: ppm, At: time.Now()}
	return s.blobs.Save("gas-sample", sensorID, reading)
}

func (s *SampleService) Latest(sensorID string) (Reading, error) {
	var reading Reading
	if err := s.blobs.Load("gas-sample", sensorID, &reading); err != nil {
		return Reading{}, err
	}
	return reading, nil
}

func (s *SampleService) History(sensorID string) ([]Reading, error) {
	ids, err := s.blobs.List("gas-sample")
	if err != nil {
		return nil, err
	}
	items := []Reading{}
	for _, id := range ids {
		if id != sensorID {
			continue
		}
		reading, err := s.Latest(id)
		if err != nil {
			return nil, err
		}
		items = append(items, reading)
	}
	return items, nil
}
