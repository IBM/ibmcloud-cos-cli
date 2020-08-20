package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/errors"

	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// PartUploadCopy - uploads an individual part of a multiple upload (file).
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func PartUploadCopy(c *cli.Context) (err error) {
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

	// Initialize UploadPartCopyInput
	input := new(s3.UploadPartCopyInput)

	// Required parameters for UploadPartCopy
	mandatory := map[string]string{
		fields.Bucket:     flags.Bucket,
		fields.Key:        flags.Key,
		fields.UploadID:   flags.UploadID,
		fields.PartNumber: flags.PartNumber,
		fields.CopySource: flags.CopySource,
	}

	//
	// Optional parameter for UploadPartCopy
	options := map[string]string{
		fields.CopySourceIfMatch:           flags.CopySourceIfMatch,
		fields.CopySourceIfModifiedSince:   flags.CopySourceIfModifiedSince,
		fields.CopySourceIfNoneMatch:       flags.CopySourceIfNoneMatch,
		fields.CopySourceIfUnmodifiedSince: flags.CopySourceIfUnmodifiedSince,
		fields.CopySourceRange:             flags.CopySourceRange,
	}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// validate source
	if path := strings.Split(*input.CopySource, "/"); len(path[0]) == 0 {
		return awserr.New("copy.source.bucket.missing", "no source bucket", nil)
	}

	// Concatenate forward slash in the beginning of the copy source
	*input.CopySource = "/" + *input.CopySource

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// UploadPartCopy Op
	var output *s3.UploadPartCopyOutput
	if output, err = client.UploadPartCopy(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
