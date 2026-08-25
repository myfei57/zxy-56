package console

import (
	"encoding/json"
	"net/http"

	"labvent/internal/audit"
	"labvent/internal/chem"
	"labvent/internal/fan"
	"labvent/internal/gas"
	"labvent/internal/hood"
	"labvent/internal/lab"
	"labvent/internal/ns"
	"labvent/internal/quota"
	"labvent/internal/sash"
	"labvent/internal/store"
	"labvent/internal/valve"
)

type Deps struct {
	Namespaces *ns.NamespaceService
	Zones      *ns.ZoneService
	Labs       *lab.LabService
	Rooms      *lab.RoomService
	Mappings   *lab.MappingService
	Partitions *lab.PartitionService
	Interlock  *lab.Interlock
	Hoods      *hood.HoodService
	Airflow    *hood.AirflowService
	Velocity   *hood.VelocityService
	Pressure   *hood.PressureService
	Baseline   *hood.BaselineService
	Doors      *hood.DoorService
	Fans       *fan.Controller
	Sashes     *sash.SashService
	Moves      *sash.MoveService
	Gas        *gas.GasService
	Samples    *gas.SampleService
	Alarms     *gas.AlarmService
	Router     *gas.ZoneRouter
	Chems      *chem.ChemService
	Expiry     *chem.ExpiryService
	Classify   *chem.ClassifyService
	Issue      *store.IssueGate
	Cabinets   *store.CabinetService
	Quota      *quota.QuotaService
	Audit      *audit.Service
	Valves     *valve.ValveService
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, errorResponse{Error: err.Error()})
}

func readJSON(r *http.Request, out any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}

func queryParam(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}
