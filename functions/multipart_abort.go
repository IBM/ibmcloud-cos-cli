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

// MultipartAbort will end a current active multipart upload instance
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func MultipartAbort(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize AbortMultipartUploadInput
	input := new(s3.AbortMultipartUploadInput)

	// Required parameters and no optional parameter for AbortMultipartUpload
	mandatory := map[string]string{
		fields.Bucket:   flags.Bucket,
		fields.Key:      flags.Key,
		fields.UploadID: flags.UploadID,
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

	// Alert User that we are performing the call
	ui.Say(T("Aborting multipart upload..."))

	// AbortMultipartUpload API
	_, err = client.AbortMultipartUpload(input)
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
	ui.Say(T("Successfully aborted a multipart upload instance with key '{{.Key}}' and bucket '{{.Bucket}}'.",
		map[string]interface{}{"Key": utils.EntityNameColor(*input.Key),
			"Bucket": utils.EntityNameColor(*input.Bucket)}))

	// Return
	return nil
}
