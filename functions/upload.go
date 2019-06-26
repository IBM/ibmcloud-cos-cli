package functions

import (
	"io"

	"github.com/IBM/ibmcloud-cos-cli/errors"

	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

//
// Upload - utilizes S3 Manager Uploader API to upload objects from S3 concurrently.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func Upload(c *cli.Context) (err error) {

	// Check the number of arguments
	// and error out if so
	if c.NArg() > 0 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		return
	}

	// Build S3 Manager UploadInput
	input := new(s3manager.UploadInput)

	// Required parameters for S3 Manager UploadInput
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
		fields.Body:   flags.File,
	}

	//
	// Optional parameters for S3 Manager UploadInput
	options := map[string]string{
		fields.CacheControl:       flags.CacheControl,
		fields.ContentDisposition: flags.ContentDisposition,
		fields.ContentEncoding:    flags.ContentEncoding,
		fields.ContentLanguage:    flags.ContentLanguage,
		fields.ContentMD5:         flags.ContentMD5,
		fields.ContentType:        flags.ContentType,
		fields.Metadata:           flags.Metadata,
	}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// Defer closing on body of the file to be uploaded
	if closeAble, ok := input.Body.(io.Closer); ok {
		defer closeAble.Close()
	}

	//
	// S3 Manager Uploader Options
	uploadOptions := map[string]string{
		fields.PartSize:           flags.PartSize,
		fields.Concurrency:        flags.Concurrency,
		fields.LeavePartsOnErrors: flags.LeavePartsOnErrors,
		fields.MaxUploadParts:     flags.MaxUploadParts,
	}

	// Load COS Context
	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}

	// Error holding
	errorHolder := new(error)

	// Build a s3manager upload
	var uploader utils.Uploader
	if uploader, err = cosContext.GetUploader(c.String(flags.Region),
		applyConfigUploader(c, uploadOptions, errorHolder)); err != nil {
		return
	}

	// Uploader Error Checking
	if *errorHolder != nil {
		return *errorHolder
	}

	// Upload Op
	var output *s3manager.UploadOutput
	if output, err = uploader.Upload(input); err != nil {
		return
	}
	// Success
	// Output the successful message
	err = cosContext.GetDisplay(c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}

// Build Uploader from the Config
func applyConfigUploader(c *cli.Context, options map[string]string, err *error) func(u *s3manager.Uploader) {
	return func(u *s3manager.Uploader) {
		mandatory := map[string]string{}
		*err = MapToSDKInput(c, u, mandatory, options)
	}
}
