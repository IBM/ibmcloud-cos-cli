package functions

import (
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/render"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// Download - utilizes S3 Manager Downloader API to download objects from S3 concurrently.
// Parameter:
//
//	CLI Context Application
//
// Returns:
//
//	Error = zero or non-zero
func Download(c *cli.Context) (err error) {

	// check the number of arguments
	if c.NArg() > 1 {
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

	// Monitor the file
	keepFile := false

	// Download location
	var dstPath string

	// In case of error removes incomplete downloads
	defer func() {
		if !keepFile && dstPath != "" {
			cosContext.Remove(dstPath)
		}
	}()

	// Build GetObjectInput
	input := new(s3.GetObjectInput)

	// Required parameters for GetObjectInput
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	//
	// Optional parameters for GetObjectInput
	options := map[string]string{
		fields.IfMatch:                    flags.IfMatch,
		fields.IfModifiedSince:            flags.IfModifiedSince,
		fields.IfNoneMatch:                flags.IfNoneMatch,
		fields.IfUnmodifiedSince:          flags.IfUnmodifiedSince,
		fields.Range:                      flags.Range,
		fields.ResponseCacheControl:       flags.ResponseCacheControl,
		fields.ResponseContentDisposition: flags.ResponseContentDisposition,
		fields.ResponseContentEncoding:    flags.ResponseContentEncoding,
		fields.ResponseContentLanguage:    flags.ResponseContentLanguage,
		fields.ResponseContentType:        flags.ResponseContentType,
		fields.ResponseExpires:            flags.ResponseExpires,
	}

	//
	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	//
	// Validate Download Location
	var file utils.WriteCloser
	if dstPath, file, err = getAndValidateDownloadPath(cosContext, c.Args().First(),
		aws.StringValue(input.Key), c.IsSet(flags.Force)); err != nil || file == nil {
		return
	}

	// Defer closing file
	defer file.Close()

	// Options for Downloader
	downloadOptions := map[string]string{
		fields.PartSize:    flags.PartSize,
		fields.Concurrency: flags.Concurrency,
	}
	// Error holder
	errorHolder := new(error)

	// Build a s3manager download
	var downloader utils.Downloader
	if downloader, err = cosContext.GetDownloader(c.String(flags.Region),
		applyConfigDownloader(c, downloadOptions, errorHolder)); err != nil {
		return
	}
	// Downloader Error Checking
	if *errorHolder != nil {
		return *errorHolder
	}

	// Download Op
	var totalBytes int64
	if totalBytes, err = downloader.Download(file, input); err != nil {
		return
	}

	// Maintain to keep file
	keepFile = true

	// Render DownloadOutput
	output := &render.DownloadOutput{
		TotalBytes: totalBytes,
	}
	// Output the successful message
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}

// Build Download from the Config
func applyConfigDownloader(c *cli.Context, options map[string]string, err *error) func(u *s3manager.Downloader) {
	return func(u *s3manager.Downloader) {
		mandatory := map[string]string{}
		*err = MapToSDKInput(c, u, mandatory, options)
	}
}
