package lab

import (
	"fmt"

	"github.com/google/uuid"

	"labvent/internal/store"
)

type Room struct {
	ID        string  `json:"id"`
	LabID     string  `json:"lab_id"`
	Name      string  `json:"name"`
	Cleanroom bool    `json:"cleanroom"`
	Section   string  `json:"section"`
	Pressure  float64 `json:"pressure"`
}

type RoomService struct {
	blobs store.Blob
}

func NewRoomService(blobs store.Blob) *RoomService {
	return &RoomService{blobs: blobs}
}

func (s *RoomService) Register(labID string, name string, section string, cleanroom bool) (Room, error) {
	if labID == "" || name == "" {
		return Room{}, fmt.Errorf("room lab and name are required")
	}
	item := Room{ID: uuid.NewString(), LabID: labID, Name: name, Section: section, Cleanroom: cleanroom}
	if err := s.blobs.Save("room", item.ID, item); err != nil {
		return Room{}, err
	}
	return item, nil
}

func (s *RoomService) Get(id string) (Room, error) {
	var item Room
	if err := s.blobs.Load("room", id, &item); err != nil {
		return Room{}, err
	}
	return item, nil
}

func (s *RoomService) List(labID string) ([]Room, error) {
	ids, err := s.blobs.List("room")
	if err != nil {
		return nil, err
	}
	items := []Room{}
	for _, id := range ids {
		item, err := s.Get(id)
		if err != nil {
			return nil, err
		}
		if item.LabID == labID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (s *RoomService) UpdatePressure(id string, pressure float64) (Room, error) {
	item, err := s.Get(id)
	if err != nil {
		return Room{}, err
	}
	item.Pressure = pressure
	if err := s.blobs.Save("room", item.ID, item); err != nil {
		return Room{}, err
	}
	return item, nil
}
