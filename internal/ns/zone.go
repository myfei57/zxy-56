package ns

import (
	"fmt"

	"github.com/google/uuid"

	"labvent/internal/store"
)

type Zone struct {
	ID          string  `json:"id"`
	NamespaceID string  `json:"namespace_id"`
	Name        string  `json:"name"`
	Area        float64 `json:"area"`
}

type ZoneService struct {
	blobs store.Blob
}

func NewZoneService(blobs store.Blob) *ZoneService {
	return &ZoneService{blobs: blobs}
}

func (s *ZoneService) Register(namespaceID string, name string, area float64) (Zone, error) {
	if namespaceID == "" || name == "" {
		return Zone{}, fmt.Errorf("zone namespace and name are required")
	}
	if area <= 0 {
		return Zone{}, fmt.Errorf("zone area must be positive")
	}
	item := Zone{ID: uuid.NewString(), NamespaceID: namespaceID, Name: name, Area: area}
	if err := s.blobs.Save("zone", item.ID, item); err != nil {
		return Zone{}, err
	}
	return item, nil
}

func (s *ZoneService) Get(id string) (Zone, error) {
	var item Zone
	if err := s.blobs.Load("zone", id, &item); err != nil {
		return Zone{}, err
	}
	return item, nil
}

func (s *ZoneService) List(namespaceID string) ([]Zone, error) {
	ids, err := s.blobs.List("zone")
	if err != nil {
		return nil, err
	}
	items := []Zone{}
	for _, id := range ids {
		item, err := s.Get(id)
		if err != nil {
			return nil, err
		}
		if item.NamespaceID == namespaceID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (s *ZoneService) Area(namespaceID string) (float64, error) {
	items, err := s.List(namespaceID)
	if err != nil {
		return 0, err
	}
	total := 0.0
	for _, item := range items {
		total += item.Area
	}
	return total, nil
}
