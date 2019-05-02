package functions

import (
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

// BucketHead determines if a bucket exists in a user's account.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketHead(c *cli.Context) error {

	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Builds HeadBucketInput
	input := new(s3.HeadBucketInput)

	// Required parameter and no optional parameter for HeadBucket
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

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// HeadBucket
	_, err = client.HeadBucket(input)
	// Error handling
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
	ui.Say(T("Bucket '{{.bucket}}' in region {{.region}} found in your IBM Cloud Object Storage account.",
		map[string]interface{}{"bucket": utils.EntityNameColor(aws.StringValue(input.Bucket)),
			"region": utils.EntityNameColor(region)}))

	// Return
	return nil
}
