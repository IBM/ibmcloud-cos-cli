package functions

import (
	"strings"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// Local Constants
const (
	incorrectUsage = "Incorrect Usage."
	noRegion       = "Unable to get region."
)

// MultipartCreate creates a new multipart upload instance, according to the Amazon AWS multipart specification
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func MultipartCreate(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize CreateMultipartUploadInput
	input := new(s3.CreateMultipartUploadInput)

	// Required parameters for CreateMultipartUpload
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	//
	// Optional parameters for CreateMultipartUpload
	options := map[string]string{
		fields.CacheControl:       flags.CacheControl,
		fields.ContentDisposition: flags.ContentDisposition,
		fields.ContentEncoding:    flags.ContentEncoding,
		fields.ContentLanguage:    flags.ContentLanguage,
		fields.ContentType:        flags.ContentType,
		fields.Metadata:           flags.Metadata,
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
	ui.Say(T("Creating new multipart upload instance..."))

	// Call the AWS Go SDK's multipart function, passing in the input object. This is where the multipart upload
	// instance is created.
	multiPartInstance, err := client.CreateMultipartUpload(input)
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

	// Output the successful message including
	// Bucket, Key and Upload ID
	ui.Say(T("Details about your multipart upload instance:"))
	ui.Say("Bucket: %s", utils.EntityNameColor(*multiPartInstance.Bucket))
	ui.Say("Key: %s", utils.EntityNameColor(*multiPartInstance.Key))
	ui.Say("Upload ID: %s", utils.EntityNameColor(*multiPartInstance.UploadId))

	// Return
	return nil
}
