package gas

type Mapping interface {
	CurrentZone(sensorID string) (string, error)
}

type ZoneRouter struct {
	mapping Mapping
	cached  map[string]string
}

func NewZoneRouter(mapping Mapping) *ZoneRouter {
	return &ZoneRouter{mapping: mapping, cached: map[string]string{}}
}

func (r *ZoneRouter) Cache(sensorID string, zoneID string) {
	r.cached[sensorID] = zoneID
}

func (r *ZoneRouter) Route(sensorID string) (string, error) {
	if zone, ok := r.cached[sensorID]; ok {
		return zone, nil
	}
	return r.mapping.CurrentZone(sensorID)
}

func (r *ZoneRouter) CachedZone(sensorID string) string {
	return r.cached[sensorID]
}
