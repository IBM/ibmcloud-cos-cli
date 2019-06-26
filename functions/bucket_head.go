package functions

import (
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// BucketHead determines if a bucket exists in a user's account.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketHead(c *cli.Context) (err error) {
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

	// Builds HeadBucketInput
	input := new(s3.HeadBucketInput)

	// Required parameter and no optional parameter for HeadBucket
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Optional parameters
	options := map[string]string{}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// HeadBucket Op
	var output *s3.HeadBucketOutput
	if output, err = client.HeadBucket(input); err != nil {
		return
	}

	// Build region for additional parameters to display
	var region string
	if region, err = cosContext.GetCurrentRegion(c.String(flags.Region)); err != nil {
		return
	}
	additionalParameters := map[string]interface{}{"region": region}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.Bool(flags.JSON)).Display(input, output, additionalParameters)

	// Return
	return
}
