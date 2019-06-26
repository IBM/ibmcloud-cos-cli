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

// MultiPartList command is to display the list of the multipart uploads of
// respective bucket
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func MultiPartList(c *cli.Context) (err error) {
	// check the number of arguments
	if c.NArg() > 0 {
		// Build Command Error Struct
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

	// Initialize ListMultipartUploadsInput
	pageIterInput := new(s3.ListMultipartUploadsInput)

	// Required parameters for ListMultipartUploadsInput
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	//
	// Optional parameters for ListMultipartUploadsInput
	options := map[string]string{
		fields.Delimiter:      flags.Delimiter,
		fields.EncodingType:   flags.EncodingType,
		fields.Prefix:         flags.Prefix,
		fields.KeyMarker:      flags.KeyMarker,
		fields.UploadIDMarker: flags.UploadIDMarker,
	}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, pageIterInput, mandatory, options); err != nil {
		return
	}

	// Initialize List Multipart Uploads Input
	input := new(s3.ListMultipartUploadsInput)
	if err = DeepCopyIntoUsingJSON(input, pageIterInput); err != nil {
		return
	}

	// Build Pagination Helper
	var paginationHelper *PaginationHelper
	var nextPageSize *int64
	// retrieves a PaginationHelper and a pointer to the next page size
	if paginationHelper, nextPageSize, err = NewPaginationHelper(c, flags.MaxItems, flags.PageSize); err != nil {
		return
	}
	// set next page size as MaxUploads
	pageIterInput.MaxUploads = nextPageSize

	// Check if Max Items is set
	if c.IsSet(flags.MaxItems) {
		// Parse if the integer is correctly set
		if value, errInner := parseInt64(c.String(flags.MaxItems)); errInner != nil {
			commandError := new(errors.CommandError)
			commandError.CLIContext = c
			commandError.Cause = errors.InvalidValue
			commandError.Flag = flags.MaxItems
			commandError.IError = errInner
			err = commandError
			return
		} else {
			input.MaxUploads = &value
		}
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// List Multipart Uploads Op
	output := new(s3.ListMultipartUploadsOutput)
	if err = client.ListMultipartUploadsPages(pageIterInput,
		ListMultipartUploadsItx(paginationHelper, output)); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}

// ListMultipartUploadsItx - Initialize List Multipart Uploads Output from the first page
func ListMultipartUploadsItx(paginationHelper *PaginationHelper, output *s3.ListMultipartUploadsOutput) func(
	*s3.ListMultipartUploadsOutput, bool) bool {
	// Set first page to true
	firstPage := true
	return func(page *s3.ListMultipartUploadsOutput, _ bool) bool {
		// Check if first page, initialize the output
		if firstPage {
			initListMultipartUploadsOutputFromPage(output, page)
			firstPage = false
		}
		// Merge subsequent pages into output
		mergeListMultipartUploadsOutputPage(output, page)
		// Return
		return paginationHelper.UpdateTotal(len(page.Uploads))
	}
}

//
// initListMultipartUploadsOutputFromPage - Initialize List Multipart Uploads Output from the first page
func initListMultipartUploadsOutputFromPage(output, page *s3.ListMultipartUploadsOutput) {
	output.Bucket = page.Bucket
	output.Delimiter = page.Delimiter
	output.EncodingType = page.EncodingType
	output.KeyMarker = page.KeyMarker
	output.Prefix = page.Prefix
	output.UploadIdMarker = page.UploadIdMarker
}

// mergeListMultipartUploadsOutputPage - Merge List Multipart Uploads Output into previous page print-outs
func mergeListMultipartUploadsOutputPage(output, page *s3.ListMultipartUploadsOutput) {
	output.CommonPrefixes = append(output.CommonPrefixes, page.CommonPrefixes...)
	output.IsTruncated = page.IsTruncated
	output.NextKeyMarker = page.NextKeyMarker
	output.NextUploadIdMarker = page.NextUploadIdMarker
	output.Uploads = append(output.Uploads, page.Uploads...)
}
