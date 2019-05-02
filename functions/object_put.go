package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// ObjectPut - puts an object to a bucket without using multipart upload
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectPut(c *cli.Context) error {

	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Build PutObjectInput
	input := new(s3.PutObjectInput)

	// Required parameters for PutObject
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	//
	// Optional parameters for PutObject
	options := map[string]string{
		fields.CacheControl:       flags.CacheControl,
		fields.ContentDisposition: flags.ContentDisposition,
		fields.ContentEncoding:    flags.ContentEncoding,
		fields.ContentLanguage:    flags.ContentLanguage,
		fields.ContentLength:      flags.ContentLength,
		fields.ContentMD5:         flags.ContentMD5,
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

	// Retrieve the user body path
	bodyFile := c.String("body")

	// If body is used in the command
	if bodyFile != "" {
		// Open the file in order to read the binary to store into cos
		// Otherwise, return an error letting users know object cannot be open
		file, err := cosContext.ReadSeekerCloserOpen(bodyFile)
		if err != nil {
			ui.Failed(T("Unable to open object '{{.Key}}' for upload.",
				map[string]interface{}{fields.Key: bodyFile}))
			return cli.NewExitError("", 2)
		}
		// Close the file afer the file operation is used
		defer file.Close()

		// Set body into the PutObjectInput
		input.SetBody(file)
	}

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Putting object..."))

	// Upload the object by calling PutObject
	_, err = client.PutObject(input)
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

	//
	// Output the successful message
	ui.Say(T("Successfully uploaded object '{{.Key}}' to bucket '{{.Bucket}}'.",
		map[string]interface{}{fields.Key: utils.EntityNameColor(*input.Key),
			fields.Bucket: utils.EntityNameColor(*input.Bucket)}))

	// Return
	return nil
}
