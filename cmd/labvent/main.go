package main

import (
	"log"
	"net/http"

	"labvent/internal/audit"
	"labvent/internal/chem"
	"labvent/internal/console"
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

func main() {
	cfg := loadConfig()
	blobs := store.New(cfg.data)
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
	deps := console.Deps{
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
	server := console.NewServer(deps)
	log.Printf("LabVent listening on %s with data at %s", cfg.addr, cfg.data)
	if err := http.ListenAndServe(cfg.addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
