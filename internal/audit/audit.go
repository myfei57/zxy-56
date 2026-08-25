package audit

import (
	"time"

	"labvent/internal/store"
)

type Entry struct {
	ID      string    `json:"id"`
	At      time.Time `json:"at"`
	Source  string    `json:"source"`
	Subject string    `json:"subject"`
	Action  string    `json:"action"`
	Detail  string    `json:"detail"`
}

type Service struct {
	blobs     store.Blob
	onConfirm func(subject string) error
	onReset   func(subject string) error
}

func NewService(blobs store.Blob) *Service {
	return &Service{blobs: blobs}
}

func (s *Service) SetConfirmCallback(fn func(subject string) error) {
	s.onConfirm = fn
}

func (s *Service) SetResetCallback(fn func(subject string) error) {
	s.onReset = fn
}

func (s *Service) Count(subject string) (int, error) {
	entries, err := s.List(subject)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}
