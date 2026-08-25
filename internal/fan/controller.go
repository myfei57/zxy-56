package fan

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"labvent/internal/audit"
	"labvent/internal/hood"
	"labvent/internal/store"
)

type RunRecord struct {
	FanID   string    `json:"fan_id"`
	Running bool      `json:"running"`
	At      time.Time `json:"at"`
}

type Controller struct {
	blobs    store.Blob
	pressure *hood.PressureService
	doors    *hood.DoorService
	audit    *audit.Service
	reserve  func(labID string, hours float64) error
}

func NewController(
	blobs store.Blob,
	pressure *hood.PressureService,
	doors *hood.DoorService,
	auditService *audit.Service,
	reserve func(labID string, hours float64) error,
) *Controller {
	return &Controller{blobs: blobs, pressure: pressure, doors: doors, audit: auditService, reserve: reserve}
}

func (c *Controller) Register(name string, rowID string, hoodID string, zoneID string, role string) (Fan, error) {
	if name == "" || rowID == "" || hoodID == "" {
		return Fan{}, fmt.Errorf("fan name, row and hood are required")
	}
	item := Fan{ID: uuid.NewString(), Name: name, RowID: rowID, HoodID: hoodID, ZoneID: zoneID, Role: role, State: StateStopped, Speed: 0}
	if err := c.blobs.Save("fan", item.ID, item); err != nil {
		return Fan{}, err
	}
	return item, nil
}

func (c *Controller) Get(id string) (Fan, error) {
	var item Fan
	if err := c.blobs.Load("fan", id, &item); err != nil {
		return Fan{}, err
	}
	return item, nil
}

func (c *Controller) List(rowID string) ([]Fan, error) {
	ids, err := c.blobs.List("fan")
	if err != nil {
		return nil, err
	}
	items := []Fan{}
	for _, id := range ids {
		item, err := c.Get(id)
		if err != nil {
			return nil, err
		}
		if item.RowID == rowID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (c *Controller) StateOf(id string) (string, error) {
	item, err := c.Get(id)
	if err != nil {
		return "", err
	}
	return item.State, nil
}

func (c *Controller) IssueStart(id string) error {
	item, err := c.Get(id)
	if err != nil {
		return err
	}
	item.State = StateStarting
	return c.blobs.Save("fan", item.ID, item)
}

func (c *Controller) ConfirmRunning(id string) error {
	item, err := c.Get(id)
	if err != nil {
		return err
	}
	item.State = StateRunning
	if err := c.blobs.Save("fan", item.ID, item); err != nil {
		return err
	}
	return c.blobs.Save("fan-run", item.ID, RunRecord{FanID: item.ID, Running: true, At: time.Now()})
}

func (c *Controller) Start(id string) error {
	if err := c.IssueStart(id); err != nil {
		return err
	}
	return c.ConfirmRunning(id)
}

func (c *Controller) Stop(id string) error {
	item, err := c.Get(id)
	if err != nil {
		return err
	}
	item.State = StateStopped
	item.Speed = 0
	return c.blobs.Save("fan", item.ID, item)
}

func (c *Controller) StandbyOf(rowID string, primaryID string) (Fan, error) {
	items, err := c.List(rowID)
	if err != nil {
		return Fan{}, err
	}
	for _, item := range items {
		if item.Role == "standby" && item.ID != primaryID {
			return item, nil
		}
	}
	return Fan{}, fmt.Errorf("no standby fan found in row %s", rowID)
}
