package gas

import (
	"fmt"

	"labvent/internal/audit"
)

type AlarmService struct {
	gas    *GasService
	latch  *Latch
	audit  *audit.Service
	router *ZoneRouter
}

func NewAlarmService(gas *GasService, latch *Latch, auditService *audit.Service, router *ZoneRouter) *AlarmService {
	return &AlarmService{gas: gas, latch: latch, audit: auditService, router: router}
}

func (s *AlarmService) Sample(sensorID string, ppm float64) error {
	sensor, err := s.gas.Get(sensorID)
	if err != nil {
		return err
	}
	state, err := s.latch.State(sensorID)
	if err != nil {
		return err
	}
	if ppm >= sensor.Threshold {
		return s.latch.Raise(sensorID)
	}
	if state.Latched && state.Confirmed {
		return s.latch.Clear(sensorID)
	}
	return nil
}

func (s *AlarmService) Confirm(sensorID string) error {
	return s.audit.Confirm("gas", sensorID, "alarm confirm")
}

func (s *AlarmService) State(sensorID string) (AlarmState, error) {
	return s.latch.State(sensorID)
}

func (s *AlarmService) Route(sensorID string) (string, error) {
	if s.router == nil {
		return "", fmt.Errorf("zone router is not configured")
	}
	return s.router.Route(sensorID)
}
