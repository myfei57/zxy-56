package hood

import (
	"fmt"
	"time"

	"labvent/internal/store"
)

type DoorState struct {
	RowID  string    `json:"row_id"`
	Open   bool      `json:"open"`
	Opened time.Time `json:"opened"`
}

type DoorService struct {
	blobs store.Blob
}

func NewDoorService(blobs store.Blob) *DoorService {
	return &DoorService{blobs: blobs}
}

func (s *DoorService) Open(rowID string) error {
	if rowID == "" {
		return fmt.Errorf("door row is required")
	}
	state := DoorState{RowID: rowID, Open: true, Opened: time.Now()}
	return s.blobs.Save("door", rowID, state)
}

func (s *DoorService) Close(rowID string) error {
	state := DoorState{RowID: rowID, Open: false, Opened: time.Now()}
	return s.blobs.Save("door", rowID, state)
}

func (s *DoorService) IsOpen(rowID string) (bool, error) {
	var state DoorState
	if err := s.blobs.Load("door", rowID, &state); err != nil {
		return false, err
	}
	return state.Open, nil
}
