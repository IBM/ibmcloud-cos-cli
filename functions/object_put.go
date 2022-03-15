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

// ObjectPut - puts an object to a bucket without using multipart upload
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectPut(c *cli.Context) (err error) {
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

	// Build PutObjectInput
	input := new(s3.PutObjectInput)

	// Required parameters for PutObject
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	// Optional parameters for PutObject
	options := map[string]string{
		fields.CacheControl:            flags.CacheControl,
		fields.ContentDisposition:      flags.ContentDisposition,
		fields.ContentEncoding:         flags.ContentEncoding,
		fields.ContentLanguage:         flags.ContentLanguage,
		fields.ContentLength:           flags.ContentLength,
		fields.ContentMD5:              flags.ContentMD5,
		fields.ContentType:             flags.ContentType,
		fields.Metadata:                flags.Metadata,
		fields.Body:                    flags.Body,
		fields.Tagging:                 flags.Tagging,
		fields.WebsiteRedirectLocation: flags.WebsiteRedirectLocation,
	}

	// Validate User Inputs
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// Deferring close of Body
	if closeAble, ok := input.Body.(io.Closer); ok {
		defer closeAble.Close()
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// PutObject Op
	var output *s3.PutObjectOutput
	if output, err = client.PutObject(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
