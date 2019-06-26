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

// Local Constants
const (
	incorrectUsage = "Incorrect Usage."
	noRegion       = "Unable to get region."
)

// MultipartCreate creates a new multipart upload instance, according to the Amazon AWS multipart specification
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func MultipartCreate(c *cli.Context) (err error) {
	// check the number of arguments
	if c.NArg() > 0 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		// Return with error
		return
	}

	// Load COS Context
	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}

	// Initialize CreateMultipartUploadInput
	input := new(s3.CreateMultipartUploadInput)

	// Required parameters for CreateMultipartUpload
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	//
	// Optional parameters for CreateMultipartUpload
	options := map[string]string{
		fields.CacheControl:       flags.CacheControl,
		fields.ContentDisposition: flags.ContentDisposition,
		fields.ContentEncoding:    flags.ContentEncoding,
		fields.ContentLanguage:    flags.ContentLanguage,
		fields.ContentType:        flags.ContentType,
		fields.Metadata:           flags.Metadata,
	}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// CreateMultipartUpload Op
	var output *s3.CreateMultipartUploadOutput
	if output, err = client.CreateMultipartUpload(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
