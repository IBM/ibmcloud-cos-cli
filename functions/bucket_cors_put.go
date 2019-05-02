package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"

	"github.com/IBM/ibmcloud-cos-cli/config/flags"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// BucketCorsPut creates a CORS configuration for the bucket.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketCorsPut(c *cli.Context) error {

	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Set PutBucketCorsInput
	input := new(s3.PutBucketCorsInput)

	// Required parameter for PutBucketCors
	mandatory := map[string]string{
		fields.Bucket:            flags.Bucket,
		fields.CORSConfiguration: flags.CorsConfiguration,
	}

	// No optional parameter for PutBucketCors
	options := map[string]string{}

	// Check through user inputs for validation
	err := MapToSDKInput(c, input, mandatory, options)
	if err != nil {
		ui.Failed(T("Incorrect Usage."))
		cli.ShowCommandHelp(c, c.Command.Name)
		ui.Say("")
		return cli.NewExitError("", 1)
	}

	// Set region from either user's input or default region in conf
	region, err := GetRegion(c, conf)
	if err != nil || region == "" {
		ui.Failed(T("Unable to get Region."))
		return cli.NewExitError("", 1)
	}

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Setting bucket CORS..."))

	// Put the bucket CORS by calling PutBucketCors
	_, err = client.PutBucketCors(input)
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
	ui.Say(T("Successfully set CORS configuration on bucket: {{.Bucket}}",
		map[string]interface{}{"Bucket": utils.EntityNameColor(*input.Bucket)}))

	// Return
	return nil
}
