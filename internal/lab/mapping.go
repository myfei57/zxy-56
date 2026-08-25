package lab

import (
	"fmt"

	"labvent/internal/store"
)

type ZoneMove struct {
	SensorID string `json:"sensor_id"`
	ToZone   string `json:"to_zone"`
}

type PartitionPlan struct {
	RoomID string     `json:"room_id"`
	Moves  []ZoneMove `json:"moves"`
}

type MappingService struct {
	blobs store.Blob
}

func NewMappingService(blobs store.Blob) *MappingService {
	return &MappingService{blobs: blobs}
}

func (s *MappingService) RegisterSensor(sensorID string, zoneID string) error {
	if sensorID == "" || zoneID == "" {
		return fmt.Errorf("sensor and zone are required")
	}
	return s.blobs.Save("sensor-zone", sensorID, ZoneBinding{SensorID: sensorID, ZoneID: zoneID})
}

func (s *MappingService) CurrentZone(sensorID string) (string, error) {
	var binding ZoneBinding
	if err := s.blobs.Load("sensor-zone", sensorID, &binding); err != nil {
		return "", err
	}
	return binding.ZoneID, nil
}

func (s *MappingService) ApplyPartition(plan PartitionPlan) error {
	if plan.RoomID == "" {
		return fmt.Errorf("partition room is required")
	}
	for _, move := range plan.Moves {
		if move.SensorID == "" || move.ToZone == "" {
			return fmt.Errorf("partition move requires sensor and zone")
		}
		if err := s.blobs.Save("sensor-zone", move.SensorID, ZoneBinding{SensorID: move.SensorID, ZoneID: move.ToZone}); err != nil {
			return err
		}
	}
	return nil
}

type ZoneBinding struct {
	SensorID string `json:"sensor_id"`
	ZoneID   string `json:"zone_id"`
}
