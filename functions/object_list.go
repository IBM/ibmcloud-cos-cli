package functions

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	. "github.com/IBM/ibmcloud-cos-cli/i18n"

	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// ObjectsList lists all the objects in a certain bucket.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectsList(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize ListObjectsInput
	input := new(s3.ListObjectsInput)

	// Required parameter for ListObjects
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Optional parameters for ListObjects
	options := map[string]string{
		fields.Delimiter:    flags.Delimiter,
		fields.EncodingType: flags.EncodingType,
		fields.Prefix:       flags.Prefix,
		fields.Marker:       flags.Marker,
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
	input.MaxKeys = nextPageSize

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Listing objects..."))

	// Initialize counter for Common Prefixes
	prefixesCounter := 0

	// Sets template for Common Prefixes
	commonPrefixesTable := ui.Table([]string{
		"Common Prefixes:",
	})

	// Create a table object that will hold all of the objects in the table.
	objectsTable := ui.Table([]string{
		T("Name"),
		T("Last Modified"),
		T("Object Size"),
	})
	// var to hold next marker
	nextMarker := ""

	// ListObjectPages API - iterates over the pages of a ListObjects operation
	err = client.ListObjectsPages(input, objectPagesIter(
		&prefixesCounter,
		&commonPrefixesTable,
		&objectsTable,
		paginationHelper,
		&nextMarker,
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

	// Check if CommonPrefixes are greater than 0
	if prefixesCounter > 0 {
		commonPrefixesTable.Print()
		ui.Say("")
	}

	// Simple formatting to be outputted, break down by
	// 0 vs 1 vs multiple objects in the bucket.
	var numItemsString string
	if paginationHelper.GetTotal() == 0 {
		numItemsString = T("no objects in bucket '") + utils.EntityNameColor(*input.Bucket) + "'.\n"
	} else if paginationHelper.GetTotal() == 1 {
		numItemsString = T("1 object in bucket '") + utils.EntityNameColor(*input.Bucket) + "':\n"
	} else {
		numItemsString = strconv.FormatInt(paginationHelper.GetTotal(), 10) + T(" objects in bucket '") +
			utils.EntityNameColor(*input.Bucket) + "':\n"
	}

	// Output the # of objects
	ui.Say(T("Found ") + numItemsString)
	if paginationHelper.GetTotal() != 0 {
		objectsTable.Print()
		ui.Say("")
	}

	// Output the next marker if users are paging sets of objects
	if nextMarker != "" {
		ui.Say(T("To retrieve the next set of objects use this Key as your --marker for the next command: "))
		ui.Say(utils.EntityNameColor(nextMarker))
		ui.Say("")
	}

	// Return
	return nil
}

// objectPagesIter - Iterate through List Object Pages
func objectPagesIter(prefixesCounter *int, commonPrefixesTable *terminal.Table, objectsTable *terminal.Table,
	paginationHelper *PaginationHelper, nextMarker *string) func(*s3.ListObjectsOutput, bool) bool {
	return func(page *s3.ListObjectsOutput, last bool) bool {
		// add common prefixes in the page to prefixes table
		*prefixesCounter += len(page.CommonPrefixes)
		for _, commonPrefix := range page.CommonPrefixes {
			(*commonPrefixesTable).Add(utils.SaneEntityForColorTerm(aws.StringValue(commonPrefix.Prefix)))
		}

		// add all contents in the current page to the objects table
		for _, key := range page.Contents {
			(*objectsTable).Add(
				utils.SaneEntityForColorTerm(aws.StringValue(key.Key)),
				aws.TimeValue(key.LastModified).Format("Jan 02, 2006 at 15:04:05"),
				FormatFileSize(*key.Size),
			)
		}

		// update the PaginationHelper with the number of retrieved objects
		// and check if more pages are needed to reach max ( when set )
		// will also re-adjust the next page value of the ListObjectsInput
		shouldContinue := paginationHelper.UpdateTotal(len(page.Contents))
		// checks if page iteration will stop and there is more pages to be retrieved
		if !shouldContinue && !last {
			// populate next marker
			*nextMarker = aws.StringValue(page.NextMarker)
		}

		// Return true or false (via shouldContinue)
		return shouldContinue
	}
}
