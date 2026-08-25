package audit

func (s *Service) Confirm(source string, subject string, detail string) error {
	if s.onConfirm != nil {
		if err := s.onConfirm(subject); err != nil {
			return err
		}
	}
	return s.Record(source, subject, "confirm", detail)
}

func (s *Service) ConfirmedCount(subject string) (int, error) {
	entries, err := s.List(subject)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.Action == "confirm" {
			count++
		}
	}
	return count, nil
}
