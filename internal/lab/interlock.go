package lab

import (
	"fmt"

	"labvent/internal/hood"
)

type Interlock struct {
	baseline *hood.BaselineService
}

func NewInterlock(baseline *hood.BaselineService) *Interlock {
	return &Interlock{baseline: baseline}
}

func (i *Interlock) Verdict(roomID string, measured float64) bool {
	base, err := i.baseline.Current(roomID)
	if err != nil {
		return false
	}
	return measured >= base
}

func (i *Interlock) Check(roomID string, measured float64) error {
	if !i.Verdict(roomID, measured) {
		return fmt.Errorf("buffer %s interlocked: pressure %.1f below baseline", roomID, measured)
	}
	return nil
}
