package lab

import (
	"fmt"

	"labvent/internal/store"
)

type PartitionService struct {
	blobs   store.Blob
	mapping *MappingService
}

func NewPartitionService(blobs store.Blob, mapping *MappingService) *PartitionService {
	return &PartitionService{blobs: blobs, mapping: mapping}
}

func (s *PartitionService) Apply(plan PartitionPlan) error {
	if err := s.mapping.ApplyPartition(plan); err != nil {
		return err
	}
	record := PartitionRecord{
		ID:     fmt.Sprintf("partition-%s", plan.RoomID),
		RoomID: plan.RoomID,
		Count:  len(plan.Moves),
	}
	return s.blobs.Save("partition", record.ID, record)
}

func (s *PartitionService) Count(roomID string) (int, error) {
	var record PartitionRecord
	if err := s.blobs.Load("partition", fmt.Sprintf("partition-%s", roomID), &record); err != nil {
		return 0, err
	}
	return record.Count, nil
}

type PartitionRecord struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
	Count  int    `json:"count"`
}
