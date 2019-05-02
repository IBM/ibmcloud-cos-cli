package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// ObjectCopy copies a file from one bucket to another. The source and destination
// buckets must be in the same region.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectCopy(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize Copy Object Input
	input := new(s3.CopyObjectInput)

	// Required parameters for CopyObject
	mandatory := map[string]string{
		fields.Bucket:     flags.Bucket,
		fields.Key:        flags.Key,
		fields.CopySource: flags.CopySource,
	}

	//
	// Optional parameters for CopyObject
	options := map[string]string{
		fields.CacheControl:                flags.CacheControl,
		fields.ContentDisposition:          flags.ContentDisposition,
		fields.ContentEncoding:             flags.ContentEncoding,
		fields.ContentLanguage:             flags.ContentLanguage,
		fields.ContentType:                 flags.ContentType,
		fields.CopySourceIfMatch:           flags.CopySourceIfMatch,
		fields.CopySourceIfModifiedSince:   flags.CopySourceIfModifiedSince,
		fields.CopySourceIfNoneMatch:       flags.CopySourceIfNoneMatch,
		fields.CopySourceIfUnmodifiedSince: flags.CopySourceIfUnmodifiedSince,
		fields.Metadata:                    flags.Metadata,
		fields.MetadataDirective:           flags.MetadataDirective,
	}

	// Validate User Inputs and Retrieve Region
	region, err := ValidateUserInputsAndSetRegion(c, input, mandatory, options, conf)
	if err != nil {
		ui.Failed(err.Error())
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(err.Error(), 1)
	}

	// Split the path to obtain the key name
	var path = strings.Split(*input.CopySource, "/")

	// Check if there is a slash in the copy source path
	if len(path[0]) == 0 {
		ui.Failed(T("Empty source bucket."))
		cli.ShowCommandHelp(c, "copy-object")
		return cli.NewExitError("", 1)
	}
	// Build copy source
	var sourceBucket = path[0]

	// Concatentate forward slash in the beginning of the copy source
	*input.CopySource = "/" + *input.CopySource

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Copying object..."))

	// Call CopyObject API
	_, err = client.CopyObject(input)
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

	//
	// Output the successful message
	ui.Say(T("Successfully copied '{{.object}}' from bucket '{{.bucket1}}' to bucket '{{.bucket2}}'.",
		map[string]interface{}{"object": utils.EntityNameColor(*input.Key),
			"bucket1": utils.EntityNameColor(sourceBucket),
			"bucket2": utils.EntityNameColor(*input.Bucket)}))

	//
	// Return
	return nil
}
