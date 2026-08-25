package fan

import (
	"fmt"
)

func (c *Controller) Failover(primaryID string) error {
	primary, err := c.Get(primaryID)
	if err != nil {
		return err
	}
	standby, err := c.StandbyOf(primary.RowID, primaryID)
	if err != nil {
		return err
	}
	if err := c.Start(standby.ID); err != nil {
		return err
	}
	healthy, err := c.pressure.Healthy(primary.HoodID)
	if err != nil {
		return err
	}
	if !healthy {
		return fmt.Errorf("hood %s pressure not healthy during failover", primary.HoodID)
	}
	if err := c.Stop(primary.ID); err != nil {
		return err
	}
	if c.audit != nil {
		_ = c.audit.Record("fan", primaryID, "failover", standby.ID)
	}
	return nil
}
