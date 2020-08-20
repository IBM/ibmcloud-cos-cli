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

// ObjectCopy copies a file from one bucket to another. The source and destination
// buckets must be in the same region.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectCopy(c *cli.Context) (err error) {
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

	// Initialize Copy Object Input
	input := new(s3.CopyObjectInput)

	// Required parameters for CopyObject
	mandatory := map[string]string{
		fields.Bucket:     flags.Bucket,
		fields.Key:        flags.Key,
		fields.CopySource: flags.CopySource,
	}

	//
	// Optional parameters for CopyObject
	options := map[string]string{
		fields.CacheControl:                flags.CacheControl,
		fields.ContentDisposition:          flags.ContentDisposition,
		fields.ContentEncoding:             flags.ContentEncoding,
		fields.ContentLanguage:             flags.ContentLanguage,
		fields.ContentType:                 flags.ContentType,
		fields.CopySourceIfMatch:           flags.CopySourceIfMatch,
		fields.CopySourceIfModifiedSince:   flags.CopySourceIfModifiedSince,
		fields.CopySourceIfNoneMatch:       flags.CopySourceIfNoneMatch,
		fields.CopySourceIfUnmodifiedSince: flags.CopySourceIfUnmodifiedSince,
		fields.Metadata:                    flags.Metadata,
		fields.MetadataDirective:           flags.MetadataDirective,
	}

	// Validate User Inputs
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

	// CopyObject Op
	var output *s3.CopyObjectOutput
	if output, err = client.CopyObject(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
