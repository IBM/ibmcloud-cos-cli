package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

const (
	errorParsingDelete = "Error parsing parameter '--delete'"
)

// ObjectDeletes deletes multiple objects from a bucket.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectDeletes(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Set DeleteObjectInput
	input := new(s3.DeleteObjectsInput)

	// Required parameters and no optional parameter for DeleteObjects
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Delete: flags.Delete,
	}

	// Validate User Inputs and Retrieve Region
	region, err := ValidateUserInputsAndSetRegion(c, input, mandatory, map[string]string{}, conf)
	if err != nil {
		ui.Failed(err.Error())
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(err.Error(), 1)
	}

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Deleting objects..."))

	// DeleteObject API
	_, err = client.DeleteObjects(input)
	// Error handling
	if err != nil {
		if strings.Contains(err.Error(), "EmptyStaticCreds") {
			ui.Failed(err.Error() + "\n" + T("Try logging in using 'ibmcloud login'."))
		} else {
			ui.Failed(err.Error())
		}
		return cli.NewExitError("", 1)
	}
	// Success
	ui.Ok()
	ui.Say(T("Delete multiple objects from bucket '{{.Bucket}}' ran successfully.", input))

	// Return
	return nil
}
