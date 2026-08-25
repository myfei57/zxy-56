package chem

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"labvent/internal/store"
)

type Chemical struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Class          string    `json:"class"`
	ExpiryAt       time.Time `json:"expiry_at"`
	EffectiveUntil time.Time `json:"effective_until"`
	Quantity       float64   `json:"quantity"`
	CabinetType    string    `json:"cabinet_type"`
}

type ChemService struct {
	blobs    store.Blob
	expiry   *ExpiryService
	classify *ClassifyService
}

func NewChemService(blobs store.Blob, expiry *ExpiryService, classify *ClassifyService) *ChemService {
	return &ChemService{blobs: blobs, expiry: expiry, classify: classify}
}

func (s *ChemService) Register(name string, class string, expiryAt time.Time, quantity float64) (Chemical, error) {
	if name == "" || class == "" {
		return Chemical{}, fmt.Errorf("chemical name and class are required")
	}
	item := Chemical{
		ID:             uuid.NewString(),
		Name:           name,
		Class:          class,
		ExpiryAt:       expiryAt,
		EffectiveUntil: expiryAt,
		Quantity:       quantity,
		CabinetType:    "ordinary",
	}
	if err := s.blobs.Save("chem", item.ID, item); err != nil {
		return Chemical{}, err
	}
	if err := s.expiry.Register(item.ID, expiryAt); err != nil {
		return Chemical{}, err
	}
	if err := s.classify.Register(item.ID, class); err != nil {
		return Chemical{}, err
	}
	return item, nil
}

func (s *ChemService) Get(id string) (Chemical, error) {
	var item Chemical
	if err := s.blobs.Load("chem", id, &item); err != nil {
		return Chemical{}, err
	}
	return item, nil
}

func (s *ChemService) List() ([]Chemical, error) {
	ids, err := s.blobs.List("chem")
	if err != nil {
		return nil, err
	}
	items := []Chemical{}
	for _, id := range ids {
		item, err := s.Get(id)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *ChemService) UpdateQuantity(id string, quantity float64) (Chemical, error) {
	item, err := s.Get(id)
	if err != nil {
		return Chemical{}, err
	}
	item.Quantity = quantity
	if err := s.blobs.Save("chem", item.ID, item); err != nil {
		return Chemical{}, err
	}
	return item, nil
}
