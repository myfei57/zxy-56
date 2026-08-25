package sash

// Face-velocity alarm thresholds (m/s). The acceptable upper limit rises as the
// sash is lowered: a smaller opening yields a higher face velocity for the same
// exhaust flow and is more prone to turbulence, so the alarm limit relaxes from
// fullOpenThreshold (sash fully open) up to lowSashThreshold (sash near closed).
// ThresholdFor and HeightForTarget are inverses across this range.
const (
	fullOpenThreshold = 0.5 // sash fully open → lowest acceptable limit
	lowSashThreshold  = 0.9 // sash near closed → highest acceptable limit
)

// ThresholdFor returns the face-velocity alarm threshold for the given sash
// height. The threshold scales linearly with the open ratio so that it tracks
// the sash in real time instead of staying pinned to the full-open value.
func ThresholdFor(fullOpen int, height int) float64 {
	if fullOpen <= 0 {
		return fullOpenThreshold
	}
	ratio := float64(height) / float64(fullOpen)
	switch {
	case ratio >= 1:
		return fullOpenThreshold
	case ratio <= 0:
		return lowSashThreshold
	}
	return lowSashThreshold - (lowSashThreshold-fullOpenThreshold)*ratio
}

// HeightForTarget is the inverse of ThresholdFor: it returns the sash height
// that yields the requested face-velocity threshold.
func HeightForTarget(threshold float64, fullOpen int) int {
	if threshold >= lowSashThreshold {
		return 1
	}
	if threshold <= fullOpenThreshold {
		return fullOpen
	}
	ratio := 1 - (threshold-fullOpenThreshold)/(lowSashThreshold-fullOpenThreshold)
	return int(float64(fullOpen) * ratio)
}
