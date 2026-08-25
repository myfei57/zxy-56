package ns

import (
	"fmt"
)

type Summary struct {
	NamespaceID string  `json:"namespace_id"`
	ZoneCount   int     `json:"zone_count"`
	TotalArea   float64 `json:"total_area"`
	AverageArea float64 `json:"average_area"`
}

func (s *NamespaceService) Summary(namespaceID string) (Summary, error) {
	zones, err := s.blobs.List("zone")
	if err != nil {
		return Summary{}, err
	}
	count := 0
	total := 0.0
	for _, id := range zones {
		var zone Zone
		if err := s.blobs.Load("zone", id, &zone); err != nil {
			return Summary{}, err
		}
		if zone.NamespaceID != namespaceID {
			continue
		}
		count++
		total += zone.Area
	}
	average := 0.0
	if count > 0 {
		average = total / float64(count)
	}
	if _, err := s.Get(namespaceID); err != nil {
		return Summary{}, fmt.Errorf("namespace %s not found", namespaceID)
	}
	return Summary{NamespaceID: namespaceID, ZoneCount: count, TotalArea: total, AverageArea: average}, nil
}
