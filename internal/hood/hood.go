package hood

import (
	"fmt"

	"github.com/google/uuid"

	"labvent/internal/store"
)

type Hood struct {
	ID     string  `json:"id"`
	LabID  string  `json:"lab_id"`
	RoomID string  `json:"room_id"`
	Name   string  `json:"name"`
	Width  float64 `json:"width"`
	Status string  `json:"status"`
}

type HoodService struct {
	blobs store.Blob
}

func NewHoodService(blobs store.Blob) *HoodService {
	return &HoodService{blobs: blobs}
}

func (s *HoodService) Register(labID string, roomID string, name string, width float64) (Hood, error) {
	if labID == "" || roomID == "" || name == "" {
		return Hood{}, fmt.Errorf("hood lab, room and name are required")
	}
	if width <= 0 {
		return Hood{}, fmt.Errorf("hood width must be positive")
	}
	item := Hood{ID: uuid.NewString(), LabID: labID, RoomID: roomID, Name: name, Width: width, Status: "idle"}
	if err := s.blobs.Save("hood", item.ID, item); err != nil {
		return Hood{}, err
	}
	return item, nil
}

func (s *HoodService) Get(id string) (Hood, error) {
	var item Hood
	if err := s.blobs.Load("hood", id, &item); err != nil {
		return Hood{}, err
	}
	return item, nil
}

func (s *HoodService) List(roomID string) ([]Hood, error) {
	ids, err := s.blobs.List("hood")
	if err != nil {
		return nil, err
	}
	items := []Hood{}
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

func (s *HoodService) SetStatus(id string, status string) (Hood, error) {
	item, err := s.Get(id)
	if err != nil {
		return Hood{}, err
	}
	item.Status = status
	if err := s.blobs.Save("hood", item.ID, item); err != nil {
		return Hood{}, err
	}
	return item, nil
}
