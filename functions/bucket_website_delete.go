package functions

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM/ibmcloud-cos-cli/errors"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/render"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// BucketWebsiteDelete removes website configuration from a bucket.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketWebsiteDelete(c *cli.Context) (err error) {
	// check the number of arguments
	if c.NArg() > 0 {
		// Build Command Error struct
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

	// Initialize input structure
	input := new(s3.DeleteBucketWebsiteInput)

	// Required parameters
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Optional parameters
	options := map[string]string{}

	// Validate User Inputs
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// No force on deleting bucket website, alert users
	if !c.Bool(flags.Force) {
		confirmed := false
		// Warn the user that they're about to delete a bucket website.
		cosContext.UI.Warn(render.WarningDeleteBucketWebsite(input))
		cosContext.UI.Prompt(render.MessageConfirmationContinue(), &terminal.PromptOptions{}).Resolve(&confirmed)
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

	// DeleteBucketWebsite Op
	var output *s3.DeleteBucketWebsiteOutput
	if output, err = client.DeleteBucketWebsite(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
