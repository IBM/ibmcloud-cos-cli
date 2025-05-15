package functions

import (
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// Version prints out the version of cos plugin
// Parameter:
//
//	CLI Context Application
//
// Returns:
//
//	Error = zero or non-zero
func Version(c *cli.Context) (err error) {
	if c.NArg() > 0 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		return
	}

	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}

	pluginName := "ibmcloud cos cli plugin"
	version := c.App.Version

	// Print version here using terminal.UI
	// If we use render.Display it will print OK before printing version
	cosContext.UI.Print("%s %s", pluginName, version)

	return
}
