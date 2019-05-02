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

// objectWaiter is a function based off of HeadObjectInput to
// monitor whether the object exists or not
type objectWaiter func(input *s3.HeadObjectInput) error

// objectWaiterBuild is a function based off of s3 interface API
// for the objectWaiter type.
type objectWaiterBuild func(api s3iface.S3API) objectWaiter

// WaitObjectExists function
func WaitObjectExists(c *cli.Context) error {
	owb := func(client s3iface.S3API) objectWaiter { return client.WaitUntilObjectExists }
	return doObjectWait(c, owb)
}

// WaitObjectNotExists function
func WaitObjectNotExists(c *cli.Context) error {
	owb := func(client s3iface.S3API) objectWaiter { return client.WaitUntilObjectNotExists }
	return doObjectWait(c, owb)
}

// doObjectWait function waits for when object exists or not
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func doObjectWait(c *cli.Context, owb objectWaiterBuild) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize HeadObjectInput
	input := new(s3.HeadObjectInput)

	// Required parameters for HeadObject
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	// Optional parameters for HeadObject
	options := map[string]string{
		fields.IfMatch:           flags.IfMatch,
		fields.IfModifiedSince:   flags.IfModifiedSince,
		fields.IfNoneMatch:       flags.IfNoneMatch,
		fields.IfUnmodifiedSince: flags.IfUnmodifiedSince,
		fields.Range:             flags.Range,
		fields.PartNumber:        flags.PartNumber,
	}

	// Validate User Inputs and Retrieve Region
	region, err := ValidateUserInputsAndSetRegion(c, input, mandatory, options, conf)
	if err != nil {
		ui.Failed(err.Error())
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(err.Error(), 1)
	}

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Wait until object condition
	err = owb(client)(input)
	// Error handling when checking object exists fails
	if err != nil {
		ui.Failed(err.Error())
		return cli.NewExitError("", 255)
	}
	// Return error
	return err

}
