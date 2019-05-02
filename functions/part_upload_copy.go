package functions

import (
	"strconv"
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// PartUploadCopy - uploads an individual part of a multiple upload (file).
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func PartUploadCopy(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize UploadPartCopyInput
	input := new(s3.UploadPartCopyInput)

	// Required parameters for UploadPartCopy
	mandatory := map[string]string{
		fields.Bucket:     flags.Bucket,
		fields.Key:        flags.Key,
		fields.UploadID:   flags.UploadID,
		fields.PartNumber: flags.PartNumber,
		fields.CopySource: flags.CopySource,
	}

	//
	// Optional parameter for UploadPartCopy
	options := map[string]string{
		fields.CopySourceIfMatch:           flags.CopySourceIfMatch,
		fields.CopySourceIfModifiedSince:   flags.CopySourceIfModifiedSince,
		fields.CopySourceIfNoneMatch:       flags.CopySourceIfNoneMatch,
		fields.CopySourceIfUnmodifiedSince: flags.CopySourceIfUnmodifiedSince,
		fields.CopySourceRange:             flags.CopySourceRange,
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
	ui.Say(T("Uploading part copy..."))

	// UploadPart API
	_, err = client.UploadPartCopy(input)
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

	// Output the successful upload party copy operation with partnumber of
	// the multipart upload
	partNum := strconv.FormatInt(*input.PartNumber, 10)
	ui.Say(T("Uploaded part copy '{{.part}}' of object '{{.object}}'.",
		map[string]interface{}{"part": utils.EntityNameColor(partNum),
			"object": utils.EntityNameColor(*input.Key)}))

	// Return
	return nil
}
