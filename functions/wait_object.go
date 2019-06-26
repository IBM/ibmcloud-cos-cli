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
func doObjectWait(c *cli.Context, owb objectWaiterBuild) (err error) {

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
	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// Wait until object condition
	if err = owb(client)(input); err != nil {
		cosContext.UI.Failed(err.Error())
		return cli.NewExitError(err, 255)
	}

	return

}
