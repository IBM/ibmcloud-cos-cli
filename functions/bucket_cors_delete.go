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

// BucketCorsDelete - Deletes CORS confuration from a bucket
// Parameter:
//     	CLI Context Application
// Returns:
//  	Error if triggered
func BucketCorsDelete(c *cli.Context) (err error) {
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

	// Set DeleteBucketCorsInput
	input := new(s3.DeleteBucketCorsInput)

	// Required parameter and no optional parameter for DeleteBucketCors
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
	// DELETE CORS
	var output *s3.DeleteBucketCorsOutput
	if output, err = client.DeleteBucketCors(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
