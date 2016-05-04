package command

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// ImportCommand is a cli.Command implementation that imports resources
// into the Terraform state.
type ImportCommand struct {
	Meta
}

func (c *ImportCommand) Run(args []string) int {
	args = c.Meta.process(args, true)

	cmdFlags := c.Meta.flagSet("import")
	cmdFlags.StringVar(&c.Meta.statePath, "state", DefaultStateFilename, "path")
	cmdFlags.IntVar(&c.Meta.parallelism, "parallelism", 0, "parallelism")
	cmdFlags.StringVar(&c.Meta.stateOutPath, "state-out", "", "path")
	cmdFlags.StringVar(&c.Meta.backupPath, "backup", "", "path")
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()
	if len(args) != 2 {
		c.Ui.Error("The import command expects two arguments.")
		cmdFlags.Usage()
		return 1
	}

	// Build the context based on the arguments given
	ctx, _, err := c.Context(contextOpts{
		Path:        configPath,
		StatePath:   c.Meta.statePath,
		Parallelism: c.Meta.parallelism,
	})
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	newState, err := ctx.Refresh()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error refreshing state: %s", err))
		return 1
	}

	log.Printf("[INFO] Writing state output to: %s", c.Meta.StateOutPath())
	if err := c.Meta.PersistState(newState); err != nil {
		c.Ui.Error(fmt.Sprintf("Error writing state file: %s", err))
		return 1
	}

	return 0
}

func (c *ImportCommand) Help() string {
	helpText := `
Usage: terraform import [options] ADDR ID

  Import existing infrastructure into your Terraform state.

  This will find and import the specified resource into your Terraform
  state, allowing existing infrastructure to come under Terraform
  management without having to be initially created by Terraform.

  The ADDR specified is the address to import the resource to. Please
  see the documentation online for resource addresses. The ID is a
  resource-specific ID to identify that resource being imported. Please
  reference the documentation for the resource type you're importing to
  determine the ID syntax to use. It typically matches directly to the ID
  that the provider uses.

  In the current state of Terraform import, the resource is only imported
  into your state file. Once it is imported, you must manually write
  configuration for the new resource or Terraform will mark it for destruction.
  Future versions of Terraform will expand the functionality of Terraform
  import.

  This command will not modify your infrastructure, but it will make
  network requests to inspect parts of your infrastructure relevant to
  the resource being imported.

Options:

  -backup=path        Path to backup the existing state file before
                      modifying. Defaults to the "-state-out" path with
                      ".backup" extension. Set to "-" to disable backup.

  -input=true         Ask for input for variables if not directly set.

  -no-color           If specified, output won't contain any color.

  -state=path         Path to read and save state (unless state-out
                      is specified). Defaults to "terraform.tfstate".

  -state-out=path     Path to write updated state file. By default, the
                      "-state" path will be used.

`
	return strings.TrimSpace(helpText)
}

func (c *ImportCommand) Synopsis() string {
	return "Import existing infrastructure into Terraform"
}
