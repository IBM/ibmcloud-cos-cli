package functions

import (
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// bucketWaiter is a function based off of HeadBucketInput to
// monitor whether the object exists or not
type bucketWaiter func(*s3.HeadBucketInput) error

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
func doBucketWait(c *cli.Context, bwb bucketWaiterBuild) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Build PutObjectInput
	input := new(s3.HeadBucketInput)

	// Required parameters and no optional parameter for Headbucket
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
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

	// Wait until bucket condition
	err = bwb(client)(input)
	// Error handling when checking bucket exists fails
	if err != nil {
		ui.Failed(err.Error())
		return cli.NewExitError("", 255)
	}
	return err
}
