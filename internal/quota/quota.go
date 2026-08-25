package quota

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"labvent/internal/store"
)

type Quota struct {
	ID       string  `json:"id"`
	LabID    string  `json:"lab_id"`
	Date     string  `json:"date"`
	MaxHours float64 `json:"max_hours"`
}

type QuotaService struct {
	blobs  store.Blob
	ledger *Ledger
}

func NewQuotaService(blobs store.Blob) *QuotaService {
	return &QuotaService{blobs: blobs, ledger: NewLedger(blobs)}
}

func (s *QuotaService) Register(labID string, maxHours float64) (Quota, error) {
	if labID == "" || maxHours <= 0 {
		return Quota{}, fmt.Errorf("quota lab and max hours are required")
	}
	item := Quota{ID: uuid.NewString(), LabID: labID, Date: time.Now().Format("2006-01-02"), MaxHours: maxHours}
	if err := s.blobs.Save("quota", item.ID, item); err != nil {
		return Quota{}, err
	}
	return item, nil
}

func (s *QuotaService) Get(id string) (Quota, error) {
	var item Quota
	if err := s.blobs.Load("quota", id, &item); err != nil {
		return Quota{}, err
	}
	return item, nil
}

func (s *QuotaService) List(labID string) ([]Quota, error) {
	ids, err := s.blobs.List("quota")
	if err != nil {
		return nil, err
	}
	items := []Quota{}
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

func (s *QuotaService) Reserve(labID string, hours float64) error {
	items, err := s.List(labID)
	if err != nil {
		return err
	}
	maxHours := 0.0
	for _, item := range items {
		if item.MaxHours > maxHours {
			maxHours = item.MaxHours
		}
	}
	if maxHours <= 0 {
		return fmt.Errorf("no ventilation quota configured for lab %s", labID)
	}
	used, err := s.ledger.Current(labID)
	if err != nil {
		return err
	}
	if used+hours > maxHours {
		return fmt.Errorf("ventilation quota exceeded for lab %s", labID)
	}
	return s.ledger.Add(labID, hours)
}

func (s *QuotaService) Release(labID string, hours float64) error {
	used, err := s.ledger.Current(labID)
	if err != nil {
		return err
	}
	if hours > used {
		return fmt.Errorf("cannot release more than used for lab %s", labID)
	}
	return s.ledger.Subtract(labID, hours)
}

func (s *QuotaService) Usage(labID string) (float64, error) {
	return s.ledger.Current(labID)
}

func (s *QuotaService) Reset(labID string) error {
	return s.ledger.Reset(labID)
}
