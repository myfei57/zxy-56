package gas

import (
	"fmt"

	"github.com/google/uuid"

	"labvent/internal/store"
)

type GasSensor struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	RoomID      string  `json:"room_id"`
	Threshold   float64 `json:"threshold"`
}

type GasService struct {
	blobs store.Blob
}

func NewGasService(blobs store.Blob) *GasService {
	return &GasService{blobs: blobs}
}

func (s *GasService) Register(name string, roomID string, threshold float64) (GasSensor, error) {
	if name == "" || roomID == "" {
		return GasSensor{}, fmt.Errorf("sensor name and room are required")
	}
	if threshold <= 0 {
		return GasSensor{}, fmt.Errorf("sensor threshold must be positive")
	}
	item := GasSensor{ID: uuid.NewString(), Name: name, RoomID: roomID, Threshold: threshold}
	if err := s.blobs.Save("gas", item.ID, item); err != nil {
		return GasSensor{}, err
	}
	return item, nil
}

func (s *GasService) Get(id string) (GasSensor, error) {
	var item GasSensor
	if err := s.blobs.Load("gas", id, &item); err != nil {
		return GasSensor{}, err
	}
	return item, nil
}

func (s *GasService) List(roomID string) ([]GasSensor, error) {
	ids, err := s.blobs.List("gas")
	if err != nil {
		return nil, err
	}
	items := []GasSensor{}
	for _, id := range ids {
		item, err := s.Get(id)
		if err != nil {
			return nil, err
		}
		if item.RoomID == roomID {
			items = append(items, item)
		}
	}
	return items, nil
}
