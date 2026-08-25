package valve

import (
	"fmt"
)

func (s *ValveService) SetStroke(id string, position int) error {
	if position < 0 || position > 100 {
		return fmt.Errorf("valve position must be between 0 and 100")
	}
	item, err := s.Get(id)
	if err != nil {
		return err
	}
	item.Position = position
	if position == 0 {
		item.Stroke = "closed"
	} else if position == 100 {
		item.Stroke = "open"
	} else {
		item.Stroke = "partial"
	}
	return s.blobs.Save("valve", item.ID, item)
}

type SequenceStep struct {
	ValveName string `json:"valve_name"`
	Position  int    `json:"position"`
}

func (s *ValveService) RunSequence(hoodID string, steps []SequenceStep) error {
	if len(steps) == 0 {
		return fmt.Errorf("valve sequence is empty")
	}
	items, err := s.List(hoodID)
	if err != nil {
		return err
	}
	byName := map[string]Valve{}
	for _, item := range items {
		byName[item.Name] = item
	}
	for _, step := range steps {
		item, ok := byName[step.ValveName]
		if !ok {
			return fmt.Errorf("valve %s not found for hood %s", step.ValveName, hoodID)
		}
		if err := s.SetStroke(item.ID, step.Position); err != nil {
			return err
		}
	}
	return nil
}
