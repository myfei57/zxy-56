package lab

import (
	"fmt"

	"labvent/internal/hood"
)

type Interlock struct {
	baseline *hood.BaselineService
	cache    map[string]float64
}

func NewInterlock(baseline *hood.BaselineService) *Interlock {
	return &Interlock{baseline: baseline, cache: map[string]float64{}}
}

func (i *Interlock) Verdict(roomID string, measured float64) bool {
	base, ok := i.cache[roomID]
	if !ok {
		current, err := i.baseline.Current(roomID)
		if err != nil {
			return false
		}
		base = current
		i.cache[roomID] = base
	}
	return measured >= base
}

func (i *Interlock) Check(roomID string, measured float64) error {
	if !i.Verdict(roomID, measured) {
		return fmt.Errorf("buffer %s interlocked: pressure %.1f below baseline", roomID, measured)
	}
	return nil
}
