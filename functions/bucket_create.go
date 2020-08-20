package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/errors"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// BucketCreate creates a new bucket, allowing the user to set the name, the class, and the region.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketCreate(c *cli.Context) (err error) {
	// check the number of arguments
	if c.NArg() > 0 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		return
	}

	// Load COS Context
	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}
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
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	var region string
	if region, err = cosContext.GetCurrentRegion(c.String(flags.Region)); err != nil {
		return
	}
	setLocationConstraint(input, region, c.String(flags.Class))

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// Put the bucket CORS by calling PutBucketCors
	var output *s3.CreateBucketOutput
	if output, err = client.CreateBucket(input); err != nil {
		return
	}

	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	return
}

// To create a bucket, AWS accepts a LocationConstraint, which holds the region and the bucket class within a
// single string, such as "us-south-flex" for a Flex bucket in region us-south. We cannot pass, for example,
// "us-south" into the Region, because then AWS will not know what bucket type to be created. So, we need to
// generate a LocationConstraint on the fly, which is generally just the location, a hyphen, and the class.
// However, if the user creates a us-geo (or eu-geo/ap-geo) bucket, we have to remove "geo" from the string
// because of AWS conventions. For example, AWS doesn't have "us-geo" region - it's just "us" for a cross
// region United States bucket. So instead of a "us-geo-standard" LocationConstraint, we would have
// "us-standard". This code takes care of these issues.
// BuildLocationConstraint builds location constraint based on region and class
// More details:
// https://console.bluemix.net/docs/services/cloud-object-storage/basics/classes.html#use-storage-classes
func setLocationConstraint(input *s3.CreateBucketInput, region, class string) {
	var locationConstraint string
	if strings.HasSuffix(strings.ToLower(region), "-geo") {
		locationConstraint = region[:len(region)-3]
	} else {
		locationConstraint = region + "-"
	}
	if "" == class {
		class = "standard"
	}
	locationConstraint += class
	config := new(s3.CreateBucketConfiguration).SetLocationConstraint(strings.ToLower(locationConstraint))
	input.CreateBucketConfiguration = config
}
