package audit

import (
	"time"

	"github.com/google/uuid"
)

func (s *Service) Record(source string, subject string, action string, detail string) error {
	entry := Entry{
		ID:      uuid.NewString(),
		At:      time.Now(),
		Source:  source,
		Subject: subject,
		Action:  action,
		Detail:  detail,
	}
	return s.blobs.Save("audit", entry.ID, entry)
}
