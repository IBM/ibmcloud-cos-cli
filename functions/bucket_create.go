package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"

	"github.com/IBM/ibmcloud-cos-cli/config/flags"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// BucketCreate creates a new bucket, allowing the user to set the name, the class, and the region.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketCreate(c *cli.Context) error {

	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Set CreateBucketInput
	input := new(s3.CreateBucketInput)

	// Required parameter for CreateBucket
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Optional parameters for CreateBucket
	options := map[string]string{
		fields.IBMServiceInstanceID: flags.IbmServiceInstanceID,
	}

	// Check through user inputs for validation
	err := MapToSDKInput(c, input, mandatory, options)
	if err != nil {
		ui.Failed(T("Incorrect Usage."))
		cli.ShowCommandHelp(c, c.Command.Name)
		ui.Say("")
		return cli.NewExitError("", 1)
	}

	// To create a bucket, AWS accepts a LocationConstraint, which holds the region and the bucket class within a
	// single string, such as "us-south-flex" for a Flex bucket in region us-south. We cannot pass, for example,
	// "us-south" into the Region, because then AWS will not know what bucket type to be created. So, we need to
	// generate a LocationConstraint on the fly, which is generally just the location, a hyphen, and the class.
	// However, if the user creates a us-geo (or eu-geo/ap-geo) bucket, we have to remove "geo" from the string
	// because of AWS conventions. For example, AWS doesn't have "us-geo" region - it's just "us" for a cross
	// region United States bucket. So instead of a "us-geo-standard" LocationConstraint, we would have
	// "us-standard". This code takes care of these issues.
	bucketConfig := new(s3.CreateBucketConfiguration)

	// Set region from either user's input or default region in conf
	region, err := GetRegion(c, conf)
	if err != nil || region == "" {
		ui.Failed(T("Unable to get Region."))
		return cli.NewExitError("", 1)
	}

	// Initialize storage class and store class name
	class := ""
	if c.String("class") == "" {
		class = "standard"
	} else {
		class = strings.ToLower(c.String("class"))
	}
	// Build Location constraint with region and class
	bucketConfig.SetLocationConstraint(BuildLocationConstraint(region, class))

	// Put the location constraint into createbucketconfiguration of CreateBucketInput
	input.SetCreateBucketConfiguration(bucketConfig)

	// Append "Vault" if the storage class of the bucket is cold
	if strings.ToLower(class) == "cold" {
		class = "Cold Vault"
	} else {
		class = strings.Title(class)
	}

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Creating bucket..."))

	// Create the bucket calling CreateBucket
	_, err = client.CreateBucket(input)
	// Error Handling
	if err != nil {
		if strings.Contains(err.Error(), "EmptyStaticCreds") {
			ui.Failed(err.Error() + "\n" + T("Try logging in using 'ibmcloud login'."))
		} else {
			ui.Failed(err.Error())
		}
		return cli.NewExitError("", 1)
	}
	// Wait until bucket exists
	err = client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(*input.Bucket),
	})
	// Error handling when checking bucket exists fails
	if err != nil {
		ui.Failed(err.Error())
		return cli.NewExitError("", 1)
	}

	// Success
	ui.Ok()
	ui.Say(T("Successfully created bucket '{{.bucket}}'.",
		map[string]interface{}{"bucket": utils.EntityNameColor(*input.Bucket)}))
	ui.Say(T("Location: %s"), utils.EntityNameColor(region))
	ui.Say(T("Class: {{.class}}",
		map[string]interface{}{"class": utils.EntityNameColor(class)}))

	// Return
	return nil
}

// BuildLocationConstraint builds location constraint based on region and class
// More details:
// https://console.bluemix.net/docs/services/cloud-object-storage/basics/classes.html#use-storage-classes
func BuildLocationConstraint(region, class string) string {
	var result string
	if "" != region {
		if strings.HasSuffix(region, "-geo") {
			result = region[:len(region)-3]
		} else {
			result = region + "-"
		}
	}
	result += class
	return result
}
