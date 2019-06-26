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

// PartsList prints out a list of all the parts currently uploaded in a multipart upload session.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func PartsList(c *cli.Context) (err error) {

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

	// Initialize ListPartsInput
	pageIterInput := new(s3.ListPartsInput)

	//
	// Required parameters for ListParts
	mandatory := map[string]string{
		fields.Bucket:   flags.Bucket,
		fields.Key:      flags.Key,
		fields.UploadID: flags.UploadID,
	}

	// Optional parameters for ListParts
	options := map[string]string{
		fields.PartNumberMarker: flags.PartNumberMarker,
	}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, pageIterInput, mandatory, options); err != nil {
		return
	}

	// Initialize Parts Input
	input := new(s3.ListPartsInput)
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
	pageIterInput.MaxParts = nextPageSize

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
			input.MaxParts = &value
		}
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// List Parts Op
	output := new(s3.ListPartsOutput)
	if err = client.ListPartsPages(pageIterInput, ListPartsItx(paginationHelper, output)); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return

}

// ListPartsItx - Initialize List Parts Output from the first page
func ListPartsItx(paginationHelper *PaginationHelper, output *s3.ListPartsOutput) func(*s3.ListPartsOutput, bool) bool {
	// Set first page to true
	firstPage := true
	return func(page *s3.ListPartsOutput, _ bool) bool {
		if firstPage {
			// Check if first page, initialize the output
			initListPartsOutputFromPage(output, page)
			firstPage = false
		}
		// Merge subsequent pages into output
		mergeListPartsOutputPage(output, page)
		// Return
		return paginationHelper.UpdateTotal(len(page.Parts))
	}
}

//
// initListPartsOutputFromPage - Initialize List Parts Output from the first page
func initListPartsOutputFromPage(output, page *s3.ListPartsOutput) {
	output.AbortDate = page.AbortDate
	output.AbortRuleId = page.AbortRuleId
	output.Bucket = page.Bucket
	output.Initiator = page.Initiator
	output.Key = page.Key
	output.Owner = page.Owner
	output.PartNumberMarker = page.PartNumberMarker
	output.StorageClass = page.StorageClass
	output.UploadId = page.UploadId
}

// mergeListPartsOutputPage - Merge List Parts Output into previous page print-outs
func mergeListPartsOutputPage(output, page *s3.ListPartsOutput) {
	output.IsTruncated = page.IsTruncated
	output.NextPartNumberMarker = page.NextPartNumberMarker
	output.Parts = append(output.Parts, page.Parts...)
}
