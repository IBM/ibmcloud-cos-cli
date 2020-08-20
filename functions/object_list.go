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

// ObjectsList lists all the objects in a certain bucket.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectsList(c *cli.Context) (err error) {
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

	// Initialize ListObjectsInput
	pageIterInput := new(s3.ListObjectsInput)

	// Required parameter for ListObjects
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	//
	// Optional parameters for ListObjects
	options := map[string]string{
		fields.Delimiter:    flags.Delimiter,
		fields.EncodingType: flags.EncodingType,
		fields.Prefix:       flags.Prefix,
		fields.Marker:       flags.Marker,
	}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, pageIterInput, mandatory, options); err != nil {
		return
	}

	// Initialize ListObjects Input
	input := new(s3.ListObjectsInput)
	if err = DeepCopyIntoUsingJSON(input, pageIterInput); err != nil {
		return
	}

	// Pagination Helper
	var paginationHelper *PaginationHelper
	var nextPageSize *int64
	// retrieves a PaginationHelper and a pointer to the next page size
	if paginationHelper, nextPageSize, err = NewPaginationHelper(c, flags.MaxItems, flags.PageSize); err != nil {
		return
	}
	// set next page size as MaxUploads
	pageIterInput.MaxKeys = nextPageSize

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
			input.MaxKeys = &value
		}
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// List Objects Op
	output := new(s3.ListObjectsOutput)
	if err = client.ListObjectsPages(pageIterInput, ListObjectsItx(paginationHelper, output)); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}

// ListObjectsItx - Initialize List Objects Output from the first page
func ListObjectsItx(paginationHelper *PaginationHelper, output *s3.ListObjectsOutput) func(
	*s3.ListObjectsOutput, bool) bool {
	// Set first page to true
	firstPage := true
	return func(page *s3.ListObjectsOutput, _ bool) bool {
		if firstPage {
			// Check if first page, initialize the output
			initListObjectsOutputFromPage(output, page)
			firstPage = false
		}
		// Merge subsequent pages into output
		mergeListObjectsOutputPage(output, page)
		// Return
		return paginationHelper.UpdateTotal(len(page.Contents))
	}
}

//
// initListObjectsOutputFromPage - Initialize List Objects Output from the first page
func initListObjectsOutputFromPage(output, page *s3.ListObjectsOutput) {
	output.Delimiter = page.Delimiter
	output.EncodingType = page.EncodingType
	output.Marker = page.Marker
	output.Name = page.Name
	output.Prefix = page.Prefix
}

// mergeListObjectsOutputPage - Merge List Objects Output into previous page print-outs
func mergeListObjectsOutputPage(output, page *s3.ListObjectsOutput) {
	output.CommonPrefixes = append(output.CommonPrefixes, page.CommonPrefixes...)
	output.Contents = append(output.Contents, page.Contents...)
	output.IsTruncated = page.IsTruncated
	output.NextMarker = page.NextMarker
}
