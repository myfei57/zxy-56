package console

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func buildDeps(t *testing.T) Deps {
	t.Helper()
	blobs := store.New(t.TempDir())
	auditService := audit.NewService(blobs)
	namespaceService := ns.NewNamespaceService(blobs)
	zoneService := ns.NewZoneService(blobs)
	labService := lab.NewLabService(blobs, zoneService)
	roomService := lab.NewRoomService(blobs)
	baselineService := hood.NewBaselineService(blobs)
	interlock := lab.NewInterlock(baselineService)
	mappingService := lab.NewMappingService(blobs)
	partitionService := lab.NewPartitionService(blobs, mappingService)
	hoodService := hood.NewHoodService(blobs)
	airflowService := hood.NewAirflowService(blobs)
	velocityService := hood.NewVelocityService(blobs)
	pressureService := hood.NewPressureService(blobs)
	doorService := hood.NewDoorService(blobs)
	quotaService := quota.NewQuotaService(blobs)
	fanController := fan.NewController(blobs, pressureService, doorService, auditService, quotaService.Reserve)
	sashService := sash.NewSashService(blobs)
	moveService := sash.NewMoveService(blobs, airflowService.Persist, velocityService.SetFaceThreshold)
	latch := gas.NewLatch(blobs)
	gasService := gas.NewGasService(blobs)
	sampleService := gas.NewSampleService(blobs)
	zoneRouter := gas.NewZoneRouter(mappingService)
	alarmService := gas.NewAlarmService(gasService, latch, auditService, zoneRouter)
	auditService.SetConfirmCallback(latch.MarkConfirmed)
	auditService.SetResetCallback(fanController.ClearVibration)
	expiryService := chem.NewExpiryService(blobs)
	classifyService := chem.NewClassifyService(blobs)
	chemService := chem.NewChemService(blobs, expiryService, classifyService)
	issueGate := store.NewIssueGate(expiryService.EffectiveDate)
	cabinetService := store.NewCabinetService(classifyService.ClassOf)
	valveService := valve.NewValveService(blobs)
	return Deps{
		Namespaces: namespaceService,
		Zones:      zoneService,
		Labs:       labService,
		Rooms:      roomService,
		Mappings:   mappingService,
		Partitions: partitionService,
		Interlock:  interlock,
		Hoods:      hoodService,
		Airflow:    airflowService,
		Velocity:   velocityService,
		Pressure:   pressureService,
		Baseline:   baselineService,
		Doors:      doorService,
		Fans:       fanController,
		Sashes:     sashService,
		Moves:      moveService,
		Gas:        gasService,
		Samples:    sampleService,
		Alarms:     alarmService,
		Router:     zoneRouter,
		Chems:      chemService,
		Expiry:     expiryService,
		Classify:   classifyService,
		Issue:      issueGate,
		Cabinets:   cabinetService,
		Quota:      quotaService,
		Audit:      auditService,
		Valves:     valveService,
	}
}

func TestHealthz(t *testing.T) {
	server := NewServer(buildDeps(t))
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
}

func TestNamespaceFlow(t *testing.T) {
	server := NewServer(buildDeps(t))
	body := bytes.NewBufferString(`{"name":"A 区","code":"A"}`)
	request := httptest.NewRequest(http.MethodPost, "/namespaces", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("unexpected status: %d body=%s", response.Code, response.Body.String())
	}
	var created map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created["name"] != "A 区" {
		t.Fatalf("unexpected created: %v", created)
	}
	request = httptest.NewRequest(http.MethodGet, "/namespaces/"+created["id"], nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
}

func TestFanEndpoints(t *testing.T) {
	server := NewServer(buildDeps(t))
	body := bytes.NewBufferString(`{"name":"末端风机","row_id":"row1","hood_id":"h1","zone_id":"z1","role":"end"}`)
	request := httptest.NewRequest(http.MethodPost, "/fans", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("unexpected status: %d body=%s", response.Code, response.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	request = httptest.NewRequest(http.MethodPost, "/fans/"+created.ID+"/start", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected start status: %d body=%s", response.Code, response.Body.String())
	}
}
