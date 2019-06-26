package functions

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/render"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// ObjectDelete deletes an object from a bucket.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectDelete(c *cli.Context) (err error) {

	// check the number of arguments
	if c.NArg() > 0 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		// Return error
		return
	}

	// Load COS Context
	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}

	// Set DeleteObjectInput
	input := new(s3.DeleteObjectInput)

	// Required parameters and no optional parameter for DeleteObject
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	options := map[string]string{}

	// Validate User Inputs
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// No force on deleting object, alert users
	if !c.Bool(flags.Force) {
		confirmed := false

		// Warn the user that they're about to delete the file (prevent accidental deletions)
		cosContext.UI.Warn(render.WarningDeleteObject(input))
		cosContext.UI.Prompt(render.MessageConfirmationContinue(), &terminal.PromptOptions{}).Resolve(&confirmed)

		// If users cancel prior, cancel operation
		if !confirmed {
			cosContext.UI.Say(render.MessageOperationCanceled())
			return
		}
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// DeleteObject Op
	var output *s3.DeleteObjectOutput
	if output, err = client.DeleteObject(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
