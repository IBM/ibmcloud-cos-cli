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

// ObjectsVersions lists all the object versions in a certain bucket.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectVersions(c *cli.Context) (err error) {
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

	// Initialize ListObjectVersionsInput
	pageIterInput := new(s3.ListObjectVersionsInput)

	// Required parameter for ListObjectVersions
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	//
	// Optional parameters for ListObjectVersions
	options := map[string]string{
		fields.Delimiter:       flags.Delimiter,
		fields.EncodingType:    flags.EncodingType,
		fields.KeyMarker:       flags.KeyMarker,
		fields.Prefix:          flags.Prefix,
		fields.VersionIdMarker: flags.VersionIdMarker,
	}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, pageIterInput, mandatory, options); err != nil {
		return
	}

	// Initialize ListObjectVersions Input
	input := new(s3.ListObjectVersionsInput)
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

	// ListObjectVersions Op
	output := new(s3.ListObjectVersionsOutput)
	if err = client.ListObjectVersionsPages(pageIterInput, ListObjectVersionsItx(paginationHelper, output)); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}

// ListObjectVersionsItx - Initialize List Object Versions Output from the first page
func ListObjectVersionsItx(paginationHelper *PaginationHelper, output *s3.ListObjectVersionsOutput) func(
	*s3.ListObjectVersionsOutput, bool) bool {
	// Set first page to true
	firstPage := true
	return func(page *s3.ListObjectVersionsOutput, _ bool) bool {
		if firstPage {
			// Check if first page, initialize the output
			initListObjectVersionsOutputFromPage(output, page)
			firstPage = false
		}
		// Merge subsequent pages into output
		mergeListObjectVersionsOutputPage(output, page)
		// Return
		return paginationHelper.UpdateTotal(len(page.Versions) + len(page.DeleteMarkers))
	}
}

// initListObjectVersionsOutputFromPage - Initialize List Object Versions Output from the first page
func initListObjectVersionsOutputFromPage(output, page *s3.ListObjectVersionsOutput) {
	output.Delimiter = page.Delimiter
	output.EncodingType = page.EncodingType
	output.KeyMarker = page.KeyMarker
	output.Name = page.Name
	output.Prefix = page.Prefix
	output.VersionIdMarker = page.VersionIdMarker
}

// mergeListObjectVersionsOutputPage - Merge List Object Versions Output into previous page print-outs
func mergeListObjectVersionsOutputPage(output, page *s3.ListObjectVersionsOutput) {
	output.CommonPrefixes = append(output.CommonPrefixes, page.CommonPrefixes...)
	output.Versions = append(output.Versions, page.Versions...)
	output.DeleteMarkers = append(output.DeleteMarkers, page.DeleteMarkers...)
	output.IsTruncated = page.IsTruncated
	output.NextKeyMarker = page.NextKeyMarker
	output.NextVersionIdMarker = page.NextVersionIdMarker
}
