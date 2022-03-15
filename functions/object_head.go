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

// ObjectHead is used to get the head of an object, i.e., determine if it exists in a bucket.
// It also gets the file size and the last time it was modified.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectHead(c *cli.Context) (err error) {
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
		fields.VersionId:         flags.VersionId,
	}

	// Validate User Inputs
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// HeadObject API
	var output *s3.HeadObjectOutput
	if output, err = client.HeadObject(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
