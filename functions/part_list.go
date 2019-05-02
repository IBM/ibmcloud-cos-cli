package functions

import (
	"strconv"
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// PartsList prints out a list of all the parts currently uploaded in a multipart upload session.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func PartsList(c *cli.Context) error {

	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize ListPartsInput
	input := new(s3.ListPartsInput)

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
	input.MaxParts = nextPageSize

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Listing parts..."))

	// create a table to display all the parts
	partsTable := ui.Table([]string{
		T("Part Number"),
		T("Last Modified"),
		T("ETag"),
		T("Size"),
	})
	// var to hold next part number marker
	nextPartNumberMarker := ""

	// ListPartPages API - iterates over the pages of a ListParts operation
	err = client.ListPartsPages(input, func(page *s3.ListPartsOutput, last bool) bool {
		// add all parts in the current page to the uploads table
		for _, part := range page.Parts {
			partsTable.Add(
				strconv.FormatInt(*part.PartNumber, 10),
				aws.TimeValue(part.LastModified).Format("Jan 02, 2006 at 15:04:05"),
				aws.StringValue(part.ETag),
				FormatFileSize(*part.Size),
			)
		}

		// update the PaginationHelper with the number of retrieved parts
		// and check if more pages are needed to reach max ( when set )
		// will also re-adjust the next page value of the ListPartsInput
		shouldContinue := paginationHelper.UpdateTotal(len(page.Parts))
		// checks if page iteration will stop and there is more pages to be retrieved
		if !shouldContinue && !last {
			// populate next part number marker
			nextPartNumberMarker = strconv.FormatInt(*page.NextPartNumberMarker, 10)
		}
		return shouldContinue
	})
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

	// Simple formatting to be outputted, break down by
	// 0 vs 1 vs multiple parts in the multiple upload.
	var numItemsString string
	if paginationHelper.GetTotal() == 0 {
		numItemsString = T("no parts in the multipart upload for '") + utils.EntityNameColor(*input.Key) + "'.\n"
	} else if paginationHelper.GetTotal() == 1 {
		numItemsString = T("1 part in the multipart upload for '") + utils.EntityNameColor(*input.Key) + "':\n"
	} else {
		numItemsString = strconv.FormatInt(paginationHelper.GetTotal(), 10) + T(" parts in the multipart upload for '") +
			utils.EntityNameColor(*input.Key) + "':\n"
	}

	// Output the # of parts
	ui.Say(T("Found ") + numItemsString)
	if paginationHelper.GetTotal() != 0 {
		partsTable.Print()
		ui.Say("")
	}

	// Output the next part number marker if users are pagination on the parts
	if nextPartNumberMarker != "" {
		ui.Say(T("To retrieve the next set of parts, use this marker as your part-number-marker for the next command: ") +
			utils.EntityNameColor(nextPartNumberMarker))
		ui.Say("")
	}

	// Return
	return nil
}
