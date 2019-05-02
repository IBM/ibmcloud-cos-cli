package functions

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/i18n"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// MultiPartList command is to display the list of the multipart uploads of
// respective bucket
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func MultiPartList(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize ListMultipartUploadsInput
	input := new(s3.ListMultipartUploadsInput)

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

	// Validate User Inputs and Retrieve Region
	region, err := ValidateUserInputsAndSetRegion(c, input, mandatory, options, conf)
	if err != nil {
		ui.Failed(err.Error())
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(err.Error(), 1)
	}

	// retrieves a PaginationHelper and a pointer to the next page size
	paginationHelper, nextPageSize := NewPaginationHelper(c, flags.MaxItems, flags.PageSize)
	// set next page size as MaxUploads
	input.MaxUploads = nextPageSize

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Listing Multipart Uploads..."))

	// Initialize Common Prefixes counter
	prefixesCounter := 0

	// Initialize Common Prefixes Table
	commonPrefixesTable := ui.Table([]string{
		"Common Prefixes:",
	})

	// create a table to display all the multi part uploads
	uploadsTable := ui.Table([]string{
		T("UploadId"),
		T("Key"),
		T("Initiated"),
	})

	// var to hold the next key marker
	var nextKeyMarker string
	// var to hold next upload id marker
	var nextUploadIDMarker string

	// call S3API page iterator for the ListMultipartUploadsInput created before
	err = client.ListMultipartUploadsPages(input, partsPagesIter(
		&prefixesCounter,
		&commonPrefixesTable,
		&uploadsTable,
		paginationHelper,
		&nextKeyMarker,
		&nextUploadIDMarker,
	))

	// Error Handling
	if err != nil {
		if strings.Contains(err.Error(), "EmptyStaticCreds") {
			ui.Failed(err.Error() + "\nTry logging in using 'ibmcloud login'.")
		} else {
			ui.Failed(err.Error())
		}
		return cli.NewExitError("", 1)
	}

	// Success
	ui.Ok()

	// Check if prefix counter is greater than 0, then print
	if prefixesCounter > 0 {
		commonPrefixesTable.Print()
		ui.Say("")
	}

	// Simple formatting to be outputted, break down by
	// 0 vs 1 vs multiple objects in the bucket.
	numItemsString := outputFormatHelper(paginationHelper, aws.StringValue(input.Bucket))

	// Output the # of objects
	ui.Say(T("Found ") + numItemsString)
	if paginationHelper.GetTotal() != 0 {
		uploadsTable.Print()
		ui.Say("")
	}

	// Output the next marker if users are paging sets of objects
	if nextKeyMarker != "" || nextUploadIDMarker != "" {
		ui.Say(T("To retrieve the next set of multipart uploads, use the following markers in the next command:"))
		ui.Say("--key-marker %s --upload-id-marker %s",
			utils.EntityNameColor(nextKeyMarker),
			utils.EntityNameColor(nextUploadIDMarker))
		ui.Say("")
	}

	// return
	return nil
}

// partsPagesIter iteratese every multipart upload
func partsPagesIter(prefixesCounter *int, commonPrefixesTable *terminal.Table, uploadsTable *terminal.Table,
	paginationHelper *PaginationHelper, nextKeyMarker *string, nextUploadIDMarker *string) func(
	*s3.ListMultipartUploadsOutput, bool) bool {
	return func(page *s3.ListMultipartUploadsOutput, last bool) bool {
		// add common prefixes in the page to prefixes table
		*prefixesCounter += len(page.CommonPrefixes)
		for _, commonPrefix := range page.CommonPrefixes {
			(*commonPrefixesTable).Add(utils.SaneEntityForColorTerm(aws.StringValue(commonPrefix.Prefix)))
		}

		// add all uploads in the current page to the uploads table
		for _, upload := range page.Uploads {
			(*uploadsTable).Add(
				aws.StringValue(upload.UploadId),
				utils.SaneEntityForColorTerm(aws.StringValue(upload.Key)),
				aws.TimeValue(upload.Initiated).Format("Jan 02, 2006 at 15:04:05"),
			)
		}

		// update the PaginationHelper with the number of retrieved uploads
		// and check if more pages are needed to reach max ( when set )
		// will also re-adjust the next page value of the ListMultipartUploadsInput
		shouldContinue := paginationHelper.UpdateTotal(len(page.Uploads))
		// checks if page iteration will stop and there is more pages to be retrieved
		if !shouldContinue && !last {
			// populate next key marker
			*nextKeyMarker = aws.StringValue(page.NextKeyMarker)
			// populate next upload id marker
			*nextUploadIDMarker = aws.StringValue(page.NextUploadIdMarker)
		}
		// Return
		return shouldContinue
	}
}

func outputFormatHelper(paginationHelper *PaginationHelper, bucket string) string {
	// Simple formatting to be outputted, break down by
	// 0 vs 1 vs multiple objects in the bucket.
	var numItemsString string
	if paginationHelper.GetTotal() == 0 {
		numItemsString = T("no multipart uploads in bucket '") +
			utils.EntityNameColor(bucket) + "'.\n"
	} else if paginationHelper.GetTotal() == 1 {
		numItemsString = T("1 multipart upload in bucket '") +
			utils.EntityNameColor(bucket) + "':\n"
	} else {
		numItemsString = strconv.FormatInt(paginationHelper.GetTotal(), 10) + T(" multipart uploads in bucket '") +
			utils.EntityNameColor(bucket) + "':\n"
	}
	return numItemsString
}
