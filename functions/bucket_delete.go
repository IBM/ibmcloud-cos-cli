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

// BucketDelete deletes a bucket from a user's account. Requires the user to provide the region of the bucket, if it's
// not provided then get DefaultRegion from credentials.json
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketDelete(c *cli.Context) (err error) {
	// check the number of arguments
	if c.NArg() > 0 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		return
	}

	// Load COS Context
	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}

	// Set DeleteBucketInput
	input := new(s3.DeleteBucketInput)

	// Required parameter and no optional parameter for DeleteBucket
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Optional parameters
	options := map[string]string{}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// No force on deleting object, alert users
	if !c.Bool(flags.Force) {
		confirmed := false
		// Warn the user that they're about to delete a bucket.
		cosContext.UI.Warn(render.WarningDeleteBucket(input))
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

	// DeleteBucket Op
	var output *s3.DeleteBucketOutput
	if output, err = client.DeleteBucket(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
