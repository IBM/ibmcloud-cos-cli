package functions

import (
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

func BucketReplicationPut(c *cli.Context) (err error) {
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

	// Initialize input structure
	input := new(s3.PutBucketReplicationInput)

	// Required parameters
	mandatory := map[string]string{
		fields.Bucket:                   flags.Bucket,
		fields.ReplicationConfiguration: flags.ReplicationConfiguration,
	}

	// Optional parameters
	options := map[string]string{}

	// Validate User Inputs
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// Setting client to do the call
	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	// PutPutBucketReplication Op
	var output *s3.PutBucketReplicationOutput
	if output, err = client.PutBucketReplication(input); err != nil {
		return
	}

	// Display either in JSON or text
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
	return
}
