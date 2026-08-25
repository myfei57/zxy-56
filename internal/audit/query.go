package audit

import (
	"sort"
)

func (s *Service) List(subject string) ([]Entry, error) {
	ids, err := s.blobs.List("audit")
	if err != nil {
		return nil, err
	}
	entries := []Entry{}
	for _, id := range ids {
		var entry Entry
		if err := s.blobs.Load("audit", id, &entry); err != nil {
			return nil, err
		}
		if subject == "" || entry.Subject == subject {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i int, j int) bool {
		return entries[i].At.After(entries[j].At)
	})
	return entries, nil
}

func (s *Service) Recent(limit int) ([]Entry, error) {
	entries, err := s.List("")
	if err != nil {
		return nil, err
	}
	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

func (s *Service) BySource(source string) ([]Entry, error) {
	entries, err := s.List("")
	if err != nil {
		return nil, err
	}
	filtered := []Entry{}
	for _, entry := range entries {
		if entry.Source == source {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}
