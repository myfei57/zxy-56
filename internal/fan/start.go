package fan

import (
	"fmt"
)

func (c *Controller) StartRow(rowID string) error {
	fans, err := c.List(rowID)
	if err != nil {
		return err
	}
	if len(fans) == 0 {
		return fmt.Errorf("row %s has no fans", rowID)
	}
	if c.reserve != nil {
		if err := c.reserve(fans[0].HoodID, 2); err != nil {
			return err
		}
	}
	for _, fan := range fans {
		if fan.Role != "end" {
			continue
		}
		if err := c.IssueStart(fan.ID); err != nil {
			return err
		}
	}
	if err := c.doors.Open(rowID); err != nil {
		return err
	}
	for _, fan := range fans {
		if fan.Role != "end" {
			continue
		}
		if err := c.ConfirmRunning(fan.ID); err != nil {
			return err
		}
	}
	for _, fan := range fans {
		if fan.Role != "end" {
			continue
		}
		current, err := c.Get(fan.ID)
		if err != nil {
			return err
		}
		if !current.Running() {
			return fmt.Errorf("end fan %s not running before door release", fan.ID)
		}
	}
	if c.audit != nil {
		_ = c.audit.Record("fan", rowID, "row-start", "morning")
	}
	return nil
}
