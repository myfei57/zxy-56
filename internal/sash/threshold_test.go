package sash

import (
	"math"
	"testing"
)

func TestThresholdForScalesWithHeight(t *testing.T) {
	const fullOpen = 500

	full := ThresholdFor(fullOpen, fullOpen)
	mid := ThresholdFor(fullOpen, fullOpen/2)
	closed := ThresholdFor(fullOpen, 0) // true zero opening → low-sash bound

	if math.Abs(full-fullOpenThreshold) > 1e-9 {
		t.Fatalf("full-open threshold: got %.4f, want %.4f", full, fullOpenThreshold)
	}
	if math.Abs(closed-lowSashThreshold) > 1e-9 {
		t.Fatalf("closed-sash threshold: got %.4f, want %.4f", closed, lowSashThreshold)
	}
	// Threshold must rise as the sash closes: full-open < mid < closed.
	if !(full < mid && mid < closed) {
		t.Fatalf("threshold must rise as sash closes: full=%.4f mid=%.4f closed=%.4f", full, mid, closed)
	}
	// Linear midpoint: half-open sits halfway between the bounds.
	wantMid := fullOpenThreshold + (lowSashThreshold-fullOpenThreshold)/2
	if math.Abs(mid-wantMid) > 1e-9 {
		t.Fatalf("mid threshold: got %.4f, want %.4f", mid, wantMid)
	}
	// A small opening must sit just under the low-sash bound (not above it,
	// and not equal to the full-open value the old bug returned).
	small := ThresholdFor(fullOpen, 1)
	if !(small > fullOpenThreshold && small < lowSashThreshold) {
		t.Fatalf("small-opening threshold must be within bounds: got %.4f", small)
	}
}

func TestThresholdForGuards(t *testing.T) {
	if got := ThresholdFor(0, 100); got != fullOpenThreshold {
		t.Fatalf("zero fullOpen should fall back to full-open threshold: got %.4f", got)
	}
	if got := ThresholdFor(500, 0); got != lowSashThreshold {
		t.Fatalf("zero height should clamp to low-sash threshold: got %.4f", got)
	}
	if got := ThresholdFor(500, 600); got != fullOpenThreshold {
		t.Fatalf("height above full-open should clamp to full-open threshold: got %.4f", got)
	}
}

func TestHeightForTargetIsInverse(t *testing.T) {
	const fullOpen = 500
	for _, h := range []int{1, 50, 125, 250, 375, 500} {
		threshold := ThresholdFor(fullOpen, h)
		got := HeightForTarget(threshold, fullOpen)
		// Allow a unit of rounding from the int truncation in HeightForTarget.
		if abs(got-h) > 1 {
			t.Fatalf("round-trip height %d: threshold=%.4f → height=%d (want ~%d)", h, threshold, got, h)
		}
	}
}

func TestHeightForTargetBounds(t *testing.T) {
	const fullOpen = 500
	if got := HeightForTarget(lowSashThreshold, fullOpen); got != 1 {
		t.Fatalf("threshold at low-sash bound: got height %d, want 1", got)
	}
	if got := HeightForTarget(fullOpenThreshold, fullOpen); got != fullOpen {
		t.Fatalf("threshold at full-open bound: got height %d, want %d", got, fullOpen)
	}
	// Out-of-range thresholds clamp to the extremes.
	if got := HeightForTarget(1.5, fullOpen); got != 1 {
		t.Fatalf("threshold above range: got height %d, want 1", got)
	}
	if got := HeightForTarget(0.1, fullOpen); got != fullOpen {
		t.Fatalf("threshold below range: got height %d, want %d", got, fullOpen)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
