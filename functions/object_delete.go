package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/urfave/cli"
)

// ObjectDelete deletes an object from a bucket.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectDelete(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Set DeleteObjectInput
	input := new(s3.DeleteObjectInput)

	// Required parameters and no optional parameter for DeleteObject
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	// Validate User Inputs and Retrieve Region
	region, err := ValidateUserInputsAndSetRegion(c, input, mandatory, map[string]string{}, conf)
	if err != nil {
		ui.Failed(err.Error())
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(err.Error(), 1)
	}

	// Check if user provides force parameter
	force := false
	if c.IsSet("force") {
		force = c.Bool("force")
	}

	// No force on deleting object, alert users
	if !force {
		resolve := false

		// Warn the user that they're about to delete the file (prevent accidental deletions)
		ui.Warn(T("WARNING: This will permanently delete the object '{{.Key}}' from the bucket '{{.Bucket}}'.",
			input))
		ui.Prompt(T("Are you sure you would like to continue?"), &terminal.PromptOptions{}).Resolve(&resolve)

		// If users cancel prior, cancel operation
		if !resolve {
			ui.Say(T("Operation canceled."))
			return cli.NewExitError("", 0)
		}
	}

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Deleting object..."))

	// DeleteObject API
	_, err = client.DeleteObject(input)
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
	ui.Say(T("Delete '{{.Key}}' from bucket '{{.Bucket}}' ran successfully.", input))

	// Return
	return nil
}
