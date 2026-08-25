package fan

import (
	"time"
)

type VibrationState struct {
	FanID  string    `json:"fan_id"`
	Latched bool     `json:"latched"`
	At     time.Time `json:"at"`
}

func (c *Controller) TripVibration(id string) error {
	if err := c.blobs.Save("vibration", id, VibrationState{FanID: id, Latched: true, At: time.Now()}); err != nil {
		return err
	}
	if err := c.Stop(id); err != nil {
		return err
	}
	if c.audit != nil {
		_ = c.audit.Record("fan", id, "vibration-trip", "excessive vibration")
	}
	return nil
}

func (c *Controller) VibrationLatched(id string) (bool, error) {
	var state VibrationState
	if err := c.blobs.Load("vibration", id, &state); err != nil {
		if !c.blobs.Exists("vibration", id) {
			return false, nil
		}
		return false, err
	}
	return state.Latched, nil
}

func (c *Controller) ClearVibration(id string) error {
	return c.blobs.Save("vibration", id, VibrationState{FanID: id, Latched: false, At: time.Now()})
}
