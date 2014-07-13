package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

// ShowCommand is a Command implementation that reads and outputs the
// contents of a Terraform plan or state file.
type ShowCommand struct {
	Meta
}

func (c *ShowCommand) Run(args []string) int {
	args = c.Meta.process(args)

	cmdFlags := flag.NewFlagSet("show", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()
	if len(args) != 1 {
		c.Ui.Error(
			"The show command expects exactly one argument with the path\n" +
				"to a Terraform state or plan file.\n")
		cmdFlags.Usage()
		return 1
	}
	path := args[0]

	var plan *terraform.Plan
	var state *terraform.State

	f, err := os.Open(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error loading file: %s", err))
		return 1
	}

	var planErr, stateErr error
	plan, err = terraform.ReadPlan(f)
	if err != nil {
		if _, err := f.Seek(0, 0); err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading file: %s", err))
			return 1
		}

		plan = nil
		planErr = err
	}
	if plan == nil {
		state, err = terraform.ReadState(f)
		if err != nil {
			stateErr = err
		}
	}
	if plan == nil && state == nil {
		c.Ui.Error(fmt.Sprintf(
			"Terraform couldn't read the given file as a state or plan file.\n"+
				"The errors while attempting to read the file as each format are\n"+
				"shown below.\n\n"+
				"State read error: %s\n\nPlan read error: %s",
			stateErr,
			planErr))
		return 1
	}

	if plan != nil {
		c.Ui.Output(FormatPlan(plan, c.Colorize()))
		return 0
	}

	c.Ui.Output(FormatState(state, c.Colorize()))
	return 0
}

func (c *ShowCommand) Help() string {
	helpText := `
Usage: terraform show [options] path

  Reads and outputs a Terraform state or plan file in a human-readable
  form.

Options:

  -no-color     If specified, output won't contain any color.

`
	return strings.TrimSpace(helpText)
}

func (c *ShowCommand) Synopsis() string {
	return "Inspect Terraform state or plan"
}
