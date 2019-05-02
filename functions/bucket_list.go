package functions

import (
	"strconv"
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// BucketsList prints out a list of all the buckets in an IBM Cloud account.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketsList(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and list of known regions
	ui := cosContext.UI
	listKnownRegions := cosContext.ListKnownRegions

	// Set ListBucketsInput
	input := new(s3.ListBucketsInput)

	// Required parameter for CreateBucket
	mandatory := map[string]string{}

	// Optional parameters for CreateBucket
	options := map[string]string{
		fields.IBMServiceInstanceID: flags.IbmServiceInstanceID,
	}

	// Check through user inputs for validation
	err := MapToSDKInput(c, input, mandatory, options)
	if err != nil {
		ui.Failed(T("Incorrect Usage."))
		cli.ShowCommandHelp(c, c.Command.Name)
		ui.Say("")
		return cli.NewExitError("", 1)
	}

	// any valid end point will be valid, using first from known list
	regions, err := listKnownRegions.GetAllRegions()
	if err != nil || len(regions) == 0 {
		ui.Failed(T("Unable to load default Region."))
		return cli.NewExitError("", 1)
	}

	// Retrieve the first listed region
	region := regions[0]

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Loading buckets..."))

	// The ListBuckets API
	result, err := client.ListBuckets(input)
	// Error Handling
	if err != nil {
		if strings.Contains(err.Error(), "EmptyStaticCreds") {
			ui.Failed(err.Error() + "\n" + T("Try logging in using 'ibmcloud login'."))
		} else {
			ui.Failed(err.Error())
		}
		return cli.NewExitError("", 1)
	}
	// Success
	ui.Ok()

	// Prepare the output
	numBuckets := len(result.Buckets)

	// Formats the number of buckets.
	if numBuckets == 0 {
		ui.Say(T("No buckets found in your account."))
		return nil // cli.NewExitError("", 0)
	} else if numBuckets == 1 {
		ui.Say(T("1 bucket found in your account:\n"))
	} else {
		ui.Say(strconv.Itoa(numBuckets) + T(" buckets found in your account:\n"))
	}

	// Create a table object to display each bucket in an organized fashion.
	table := ui.Table([]string{T("Name"), T("Date Created")})

	for _, b := range result.Buckets {
		// Add each bucket's name and date created to the table.
		t := aws.TimeValue(b.CreationDate)
		// Format the "Date Created" in a certain way.
		table.Add(utils.EntityNameColor(aws.StringValue(b.Name)), t.Format("Jan 02, 2006 at 15:04:05"))
	}
	table.Print()
	ui.Say("")

	// Return
	return nil
}
