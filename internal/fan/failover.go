package fan

import "fmt"

// Failover switches the row's exhaust load from the primary fan to its standby.
//
// The order is make-before-break: the standby is started and confirmed running
// BEFORE the primary is stopped, so cabinet negative pressure is held by at
// least one fan throughout the switchover and never drops to zero. Stopping the
// primary first (break-before-make) leaves a window where no fan runs and
// volatile gases can escape. If the standby fails to reach the running state,
// the primary is left running untouched so the existing airflow is preserved.
func (c *Controller) Failover(primaryID string) error {
	primary, err := c.Get(primaryID)
	if err != nil {
		return err
	}
	standby, err := c.StandbyOf(primary.RowID, primaryID)
	if err != nil {
		return err
	}
	// Make before break: bring the standby up first so at least one fan keeps
	// negative pressure while the primary is withdrawn.
	if err := c.Start(standby.ID); err != nil {
		if c.audit != nil {
			_ = c.audit.Record("fan", primaryID, "failover-aborted", "standby start failed, primary left running")
		}
		return fmt.Errorf("failover aborted: standby %s failed to start, primary %s left running: %w", standby.ID, primary.ID, err)
	}
	// Safety gate: only withdraw the primary once the standby is confirmed
	// running. Never stop the existing fan on the assumption the new one is up.
	standbyNow, err := c.Get(standby.ID)
	if err != nil {
		return err
	}
	if !standbyNow.Running() {
		if c.audit != nil {
			_ = c.audit.Record("fan", primaryID, "failover-aborted", "standby not running, primary left running")
		}
		return fmt.Errorf("failover aborted: standby %s not running (state %s), primary %s left running", standby.ID, standbyNow.State, primary.ID)
	}
	if err := c.Stop(primary.ID); err != nil {
		return err
	}
	if c.audit != nil {
		_ = c.audit.Record("fan", primaryID, "failover", standby.ID)
	}
	return nil
}
