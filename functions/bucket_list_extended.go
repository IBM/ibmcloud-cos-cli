package functions

import (
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// BucketsListExtended prints out an extensive list of all the buckets with their respective
// location constraints in an IBM Cloud account.  This also supports pagination.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketsListExtended(c *cli.Context) (err error) {
	// check the number of arguments
	if c.NArg() > 0 {
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

	// Set ListBucketsExtendedInput
	pageIterInput := new(s3.ListBucketsExtendedInput)

	// Required parameter for ListBucketsExtended
	mandatory := map[string]string{}

	//
	// Optional parameters for ListBucketsExtended
	options := map[string]string{
		fields.IBMServiceInstanceID: flags.IbmServiceInstanceID,
		fields.Marker:               flags.Marker,
		fields.Prefix:               flags.Prefix,
	}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, pageIterInput, mandatory, options); err != nil {
		return
	}

	// Initialize List Buckets Extended
	input := new(s3.ListBucketsExtendedInput)
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

	// ListBucketExtends Op
	output := new(s3.ListBucketsExtendedOutput)
	if err = client.ListBucketsExtendedPages(pageIterInput,
		ListBucketsExtendedItx(paginationHelper, output)); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}

// ListBucketsExtendedItx - iterate through each page of List Buckets Extended
func ListBucketsExtendedItx(paginationHelper *PaginationHelper, output *s3.ListBucketsExtendedOutput) func(
	*s3.ListBucketsExtendedOutput, bool) bool {
	// Set first page to true
	firstPage := true
	return func(page *s3.ListBucketsExtendedOutput, _ bool) bool {
		// Check if first page, initialize the output
		if firstPage {
			initListBucketsExtendedOutputFromPage(output, page)
			firstPage = false
		}
		// Merge subsequent pages into output
		mergeListBucketsExtendedOutputPage(output, page)

		// Return
		return paginationHelper.UpdateTotal(len(page.Buckets))
	}
}

// Initialize List Bucket Extended Output from the first page
func initListBucketsExtendedOutputFromPage(output, page *s3.ListBucketsExtendedOutput) {
	output.Marker = page.Marker
	output.Owner = page.Owner
	output.Prefix = page.Prefix
}

// Merge List Bucket Extended Output into previous page print-outs
func mergeListBucketsExtendedOutputPage(output, page *s3.ListBucketsExtendedOutput) {
	output.Buckets = append(output.Buckets, page.Buckets...)
	output.IsTruncated = page.IsTruncated
}
