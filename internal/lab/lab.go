package lab

import (
	"fmt"

	"github.com/google/uuid"

	"labvent/internal/ns"
	"labvent/internal/store"
)

type Lab struct {
	ID          string `json:"id"`
	NamespaceID string `json:"namespace_id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
}

type LabService struct {
	blobs store.Blob
	zones *ns.ZoneService
}

func NewLabService(blobs store.Blob, zones *ns.ZoneService) *LabService {
	return &LabService{blobs: blobs, zones: zones}
}

func (s *LabService) Register(namespaceID string, name string, code string) (Lab, error) {
	if namespaceID == "" || name == "" || code == "" {
		return Lab{}, fmt.Errorf("lab namespace, name and code are required")
	}
	if _, err := s.zones.Area(namespaceID); err != nil {
		return Lab{}, err
	}
	item := Lab{ID: uuid.NewString(), NamespaceID: namespaceID, Name: name, Code: code}
	if err := s.blobs.Save("lab", item.ID, item); err != nil {
		return Lab{}, err
	}
	return item, nil
}

func (s *LabService) Get(id string) (Lab, error) {
	var item Lab
	if err := s.blobs.Load("lab", id, &item); err != nil {
		return Lab{}, err
	}
	return item, nil
}

func (s *LabService) List(namespaceID string) ([]Lab, error) {
	ids, err := s.blobs.List("lab")
	if err != nil {
		return nil, err
	}
	items := []Lab{}
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
