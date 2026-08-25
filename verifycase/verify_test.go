package verifycase

import (
	"testing"

	"labvent/internal/gas"
	"labvent/internal/lab"
	"labvent/internal/store"
)

func TestLvSensorZoneMappingFresh(t *testing.T) {
	dir := t.TempDir()
	blobs := store.New(dir)
	mapper := lab.NewMappingService(blobs)
	router := gas.NewZoneRouter(mapper)
	if err := mapper.RegisterSensor("s1", "zA"); err != nil {
		t.Fatal(err)
	}
	router.Cache("s1", "zA")
	plan := lab.PartitionPlan{
		RoomID: "room1",
		Moves:  []lab.ZoneMove{{SensorID: "s1", ToZone: "zB"}},
	}
	if err := mapper.ApplyPartition(plan); err != nil {
		t.Fatal(err)
	}
	zoneID, err := router.Route("s1")
	if err != nil {
		t.Fatal(err)
	}
	if zoneID != "zB" {
		t.Fatalf("alarm route must follow the refreshed zone mapping, got %s", zoneID)
	}
}
