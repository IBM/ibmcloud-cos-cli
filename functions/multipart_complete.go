package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"

	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/urfave/cli"
)

// MultipartComplete completes a multipart upload instance, calling the appropriate Amazon AWS function to do so.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func MultipartComplete(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize CompleteMultipartUploadInput
	input := new(s3.CompleteMultipartUploadInput)

	// Required parameters for CompleteMultipartUpload
	mandatory := map[string]string{
		fields.Bucket:   flags.Bucket,
		fields.Key:      flags.Key,
		fields.UploadID: flags.UploadID,
	}

	// Optional parameter for CompleteMultipartUpload
	options := map[string]string{
		fields.MultipartUpload: flags.MultipartUpload,
	}

	// Validate User Inputs and Retrieve Region
	region, err := ValidateUserInputsAndSetRegion(c, input, mandatory, options, conf)
	if err != nil {
		ui.Failed(err.Error())
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(err.Error(), 1)
	}

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Completing multipart upload..."))

	// CompleteMultipartUpload API
	_, err = client.CompleteMultipartUpload(input)
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

	// Output the successful message
	ui.Say(T("Successfully uploaded '{{.Key}}' to bucket '{{.Bucket}}'.",
		map[string]interface{}{"Key": utils.EntityNameColor(*input.Key),
			"Bucket": utils.EntityNameColor(*input.Bucket)}))

	// Return
	return nil
}
