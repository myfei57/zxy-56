package gas

import (
	"time"

	"labvent/internal/store"
)

type AlarmState struct {
	SensorID  string    `json:"sensor_id"`
	Latched   bool      `json:"latched"`
	Confirmed bool      `json:"confirmed"`
	Since     time.Time `json:"since"`
}

type Latch struct {
	blobs store.Blob
}

func NewLatch(blobs store.Blob) *Latch {
	return &Latch{blobs: blobs}
}

func (l *Latch) State(sensorID string) (AlarmState, error) {
	var state AlarmState
	if err := l.blobs.Load("gas-latch", sensorID, &state); err != nil {
		if !l.blobs.Exists("gas-latch", sensorID) {
			return AlarmState{SensorID: sensorID}, nil
		}
		return AlarmState{}, err
	}
	return state, nil
}

func (l *Latch) Raise(sensorID string) error {
	state, err := l.State(sensorID)
	if err != nil {
		return err
	}
	state.SensorID = sensorID
	state.Latched = true
	state.Since = time.Now()
	return l.blobs.Save("gas-latch", sensorID, state)
}

func (l *Latch) Clear(sensorID string) error {
	state := AlarmState{SensorID: sensorID}
	return l.blobs.Save("gas-latch", sensorID, state)
}

func (l *Latch) MarkConfirmed(sensorID string) error {
	state, err := l.State(sensorID)
	if err != nil {
		return err
	}
	state.Confirmed = true
	return l.blobs.Save("gas-latch", sensorID, state)
}
