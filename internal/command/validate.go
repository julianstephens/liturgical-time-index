package command

import (
	"github.com/julianstephens/go-utils/cliutil"
	"github.com/julianstephens/liturgical-time-index/internal/plan"
)

type ValidateCmd struct {
	Plan string `help:"Path to the plan YAML file." type:"existingfile"`
}

func (c *ValidateCmd) Run() error {
	if _, err := plan.LoadAndValidatePlan(c.Plan); err != nil {
		cliutil.PrintError("Plan validation failed")
		return err
	}

	cliutil.PrintSuccess("Plan valid.")

	return nil
}
