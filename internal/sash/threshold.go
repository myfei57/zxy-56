package sash

func ThresholdFor(fullOpen int, height int) float64 {
	return 0.5
}

func HeightForTarget(threshold float64, fullOpen int) int {
	if threshold >= 0.9 {
		return 1
	}
	if threshold <= 0.5 {
		return fullOpen
	}
	ratio := 1 - (threshold-0.5)/0.4
	return int(float64(fullOpen) * ratio)
}
