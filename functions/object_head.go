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

// ObjectHead is used to get the head of an object, i.e., determine if it exists in a bucket.
// It also gets the file size and the last time it was modified.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectHead(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Initialize HeadObjectInput
	input := new(s3.HeadObjectInput)

	// Required parameters for HeadObject
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	// Optional parameters for HeadObject
	options := map[string]string{
		fields.IfMatch:           flags.IfMatch,
		fields.IfModifiedSince:   flags.IfModifiedSince,
		fields.IfNoneMatch:       flags.IfNoneMatch,
		fields.IfUnmodifiedSince: flags.IfUnmodifiedSince,
		fields.Range:             flags.Range,
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

	// HeadObject API
	result, err := client.HeadObject(input)
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

	/* Formatting the output containing the following:
	 * ContentLength : Size of the body in bytes.
	 * Last Modified : date of which it has been modified last
	 * The header can also contain much more information
	 **/
	ui.Say(T("Object '{{.Key}}' was found in bucket '{{.Bucket}}'.",
		map[string]interface{}{fields.Key: utils.EntityNameColor(*input.Key),
			fields.Bucket: utils.EntityNameColor(*input.Bucket)}))
	ui.Say(T("Object Size: {{.objectsize}}", map[string]interface{}{"objectsize": FormatFileSize(*result.ContentLength)}))
	ui.Say(T("Last Modified: {{.lastmodified}}",
		map[string]interface{}{"lastmodified": result.LastModified.Format("Monday, January 02, 2006 at 15:04:05")}))

	// Return
	return nil
}
