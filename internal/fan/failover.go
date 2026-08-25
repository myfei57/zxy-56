package fan

func (c *Controller) Failover(primaryID string) error {
	primary, err := c.Get(primaryID)
	if err != nil {
		return err
	}
	standby, err := c.StandbyOf(primary.RowID, primaryID)
	if err != nil {
		return err
	}
	if err := c.Stop(primary.ID); err != nil {
		return err
	}
	if err := c.Start(standby.ID); err != nil {
		return err
	}
	if c.audit != nil {
		_ = c.audit.Record("fan", primaryID, "failover", standby.ID)
	}
	return nil
}
