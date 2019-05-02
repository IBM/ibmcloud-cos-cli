package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"

	"github.com/IBM/ibmcloud-cos-cli/config/flags"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
)

// BucketDelete deletes a bucket from a user's account. Requires the user to provide the region of the bucket, if it's
// not provided then get DefaultRegion from credentials.json
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketDelete(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Set DeleteBucketInput
	input := new(s3.DeleteBucketInput)

	// Required parameter and no optional parameter for DeleteBucket
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Validate User Inputs and Retrieve Region
	region, err := ValidateUserInputsAndSetRegion(c, input, mandatory, map[string]string{}, conf)
	if err != nil {
		ui.Failed(err.Error())
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(err.Error(), 1)
	}

	// Check if user opts to delete bucket by force
	force := false
	if c.IsSet(flags.Force) {
		force = c.Bool(flags.Force)
	}

	// If this flag is activated, then bypass the warning that the user is about to delete the bucket.
	if !force {
		confirmed := false

		// Warn the user that they're about to delete a bucket.
		ui.Warn(T("WARNING: This will permanently delete the bucket '{{.Bucket}}' from your account.", input))
		ui.Prompt(T("Are you sure you would like to continue?"), &terminal.PromptOptions{}).Resolve(&confirmed)

		if !confirmed {
			ui.Say(T("Operation canceled."))
			return cli.NewExitError("", 0)
		}
	}

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Deleting bucket..."))

	// Delete the bucket by calling DeleteBucket and creating a DeleteBucketInput object that holds the name of
	// the bucket to be deleted.
	_, err = client.DeleteBucket(input)
	// Error handling
	if err != nil {
		if strings.Contains(err.Error(), "EmptyStaticCreds") {
			ui.Failed(err.Error() + "\n" + T("Try logging in using 'ibmcloud login'."))
		} else {
			ui.Failed(err.Error())
		}
		return cli.NewExitError("", 1)
	}

	// Stall until the bucket is deleted.
	err = client.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: input.Bucket,
	})
	if err != nil {
		ui.Failed(err.Error())
		return cli.NewExitError("", 1)
	}
	// Success
	ui.Ok()

	// Output the successful message
	ui.Say(T("Successfully deleted bucket '{{.Bucket}}'. The bucket '{{.Bucket}}' will be available for reuse after 15 minutes.",
		map[string]interface{}{"Bucket": utils.EntityNameColor(*input.Bucket)}))

	// Return
	return nil
}
