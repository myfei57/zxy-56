package audit

func (s *Service) Reset(source string, subject string, detail string) error {
	if s.onReset != nil {
		if err := s.onReset(subject); err != nil {
			return err
		}
	}
	return s.Record(source, subject, "reset", detail)
}

func (s *Service) ResetCount(subject string) (int, error) {
	entries, err := s.List(subject)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.Action == "reset" {
			count++
		}
	}
	return count, nil
}
