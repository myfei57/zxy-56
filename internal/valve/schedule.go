package valve

import (
	"fmt"
)

type Schedule struct {
	HoodID string         `json:"hood_id"`
	Steps  []SequenceStep `json:"steps"`
}

func BuildSchedule(hoodID string, supplyName string, exhaustName string) Schedule {
	return Schedule{
		HoodID: hoodID,
		Steps: []SequenceStep{
			{ValveName: supplyName, Position: 30},
			{ValveName: exhaustName, Position: 100},
			{ValveName: supplyName, Position: 80},
		},
	}
}

func (s *ValveService) RunSchedule(schedule Schedule) error {
	if schedule.HoodID == "" || len(schedule.Steps) == 0 {
		return fmt.Errorf("schedule requires hood and steps")
	}
	return s.RunSequence(schedule.HoodID, schedule.Steps)
}
