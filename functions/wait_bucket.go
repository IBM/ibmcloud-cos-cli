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

// bucketWaiter is a function based off of HeadBucketInput to
// monitor whether the object exists or not
type bucketWaiter func(*s3.HeadBucketInput) error

// bucketWaiterBuild type
type bucketWaiterBuild func(s3iface.S3API) bucketWaiter

// WaitBucketExists function
func WaitBucketExists(c *cli.Context) error {
	bwb := func(client s3iface.S3API) bucketWaiter { return client.WaitUntilBucketExists }
	return doBucketWait(c, bwb)
}

// WaitBucketNotExists function
func WaitBucketNotExists(c *cli.Context) error {
	bwb := func(client s3iface.S3API) bucketWaiter { return client.WaitUntilBucketNotExists }
	return doBucketWait(c, bwb)
}

// doBucketWait function waits for when object exists or not
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func doBucketWait(c *cli.Context, bwb bucketWaiterBuild) (err error) {

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

	// Build PutObjectInput
	input := new(s3.HeadBucketInput)

	// Required parameters and no optional parameter for Headbucket
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

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

	// Wait until bucket condition
	if err = bwb(client)(input); err != nil {
		cosContext.UI.Failed(err.Error())
		return cli.NewExitError(err, 255)
	}

	// Return
	return
}
