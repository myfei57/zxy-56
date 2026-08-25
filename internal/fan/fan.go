package fan

type Fan struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	RowID  string `json:"row_id"`
	HoodID string `json:"hood_id"`
	ZoneID string `json:"zone_id"`
	Role   string `json:"role"`
	State  string `json:"state"`
	Speed  int    `json:"speed"`
}

const (
	StateStopped  = "stopped"
	StateStarting = "starting"
	StateRunning  = "running"
)

func (f Fan) Running() bool {
	return f.State == StateRunning
}
