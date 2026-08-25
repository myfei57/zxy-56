package ns

import (
	"fmt"

	"github.com/google/uuid"

	"labvent/internal/store"
)

type Namespace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type NamespaceService struct {
	blobs store.Blob
}

func NewNamespaceService(blobs store.Blob) *NamespaceService {
	return &NamespaceService{blobs: blobs}
}

func (s *NamespaceService) Register(name string, code string) (Namespace, error) {
	if name == "" || code == "" {
		return Namespace{}, fmt.Errorf("namespace name and code are required")
	}
	item := Namespace{ID: uuid.NewString(), Name: name, Code: code}
	if err := s.blobs.Save("namespace", item.ID, item); err != nil {
		return Namespace{}, err
	}
	return item, nil
}

func (s *NamespaceService) Get(id string) (Namespace, error) {
	var item Namespace
	if err := s.blobs.Load("namespace", id, &item); err != nil {
		return Namespace{}, err
	}
	return item, nil
}

func (s *NamespaceService) List() ([]Namespace, error) {
	ids, err := s.blobs.List("namespace")
	if err != nil {
		return nil, err
	}
	items := []Namespace{}
	for _, id := range ids {
		item, err := s.Get(id)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
