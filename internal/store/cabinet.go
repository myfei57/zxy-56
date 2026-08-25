package store

import (
	"fmt"
)

type CabinetService struct {
	classOf func(chemID string) (string, error)
}

func NewCabinetService(classOf func(chemID string) (string, error)) *CabinetService {
	return &CabinetService{classOf: classOf}
}

func (c *CabinetService) Required(chemID string) (string, error) {
	cls, err := c.classOf(chemID)
	if err != nil {
		return "", err
	}
	if cls == "controlled" {
		return "controlled", nil
	}
	return "ordinary", nil
}

func (c *CabinetService) Verify(chemID string, cabinetType string) error {
	want, err := c.Required(chemID)
	if err != nil {
		return err
	}
	if cabinetType != want {
		return fmt.Errorf("reagent %s requires %s cabinet, found %s", chemID, want, cabinetType)
	}
	return nil
}
