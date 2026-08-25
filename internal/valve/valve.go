package valve

import (
	"fmt"

	"github.com/google/uuid"

	"labvent/internal/store"
)

type Valve struct {
	ID       string `json:"id"`
	HoodID   string `json:"hood_id"`
	Name     string `json:"name"`
	Position int    `json:"position"`
	Stroke   string `json:"stroke"`
}

type ValveService struct {
	blobs store.Blob
}

func NewValveService(blobs store.Blob) *ValveService {
	return &ValveService{blobs: blobs}
}

func (s *ValveService) Register(hoodID string, name string) (Valve, error) {
	if hoodID == "" || name == "" {
		return Valve{}, fmt.Errorf("valve hood and name are required")
	}
	item := Valve{ID: uuid.NewString(), HoodID: hoodID, Name: name, Position: 0, Stroke: "closed"}
	if err := s.blobs.Save("valve", item.ID, item); err != nil {
		return Valve{}, err
	}
	return item, nil
}

func (s *ValveService) Get(id string) (Valve, error) {
	var item Valve
	if err := s.blobs.Load("valve", id, &item); err != nil {
		return Valve{}, err
	}
	return item, nil
}

func (s *ValveService) List(hoodID string) ([]Valve, error) {
	ids, err := s.blobs.List("valve")
	if err != nil {
		return nil, err
	}
	items := []Valve{}
	for _, id := range ids {
		item, err := s.Get(id)
		if err != nil {
			return nil, err
		}
		if item.HoodID == hoodID {
			items = append(items, item)
		}
	}
	return items, nil
}
