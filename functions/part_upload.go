package functions

import (
	"io"

	"github.com/IBM/ibmcloud-cos-cli/errors"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// PartUpload - uploads an individual part of a multiple upload (file).
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func PartUpload(c *cli.Context) (err error) {
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

	// Initialize UploadPartInput
	input := new(s3.UploadPartInput)

	// Required parameters for UploadPart
	mandatory := map[string]string{
		fields.Bucket:     flags.Bucket,
		fields.Key:        flags.Key,
		fields.UploadID:   flags.UploadID,
		fields.PartNumber: flags.PartNumber,
	}

	// Optional parameters for UploadPart
	options := map[string]string{
		fields.ContentLength: flags.ContentLength,
		fields.ContentMD5:    flags.ContentMD5,
		fields.Body:          flags.Body,
	}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// Defer closing body
	if closeAble, ok := input.Body.(io.Closer); ok {
		defer closeAble.Close()
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// UploadPart Op
	var output *s3.UploadPartOutput
	if output, err = client.UploadPart(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
