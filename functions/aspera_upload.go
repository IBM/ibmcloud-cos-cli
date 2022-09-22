package functions

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	sdk "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
	"github.com/IBM/ibmcloud-cos-cli/aspera"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/render"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

func AsperaUpload(c *cli.Context) (err error) {

	// check the number of arguments
	if c.NArg() > 1 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		return
	}

	// Check required Environment Var
	APIKeyEnv := sdk.EnvAPIKey.Get()
	if APIKeyEnv == "" {
		err = fmt.Errorf(T("missing Environment Variable: %s"), "IBMCLOUD_API_KEY")
		return
	}

	// Load COS Context
	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}

	// Download SDK if not present
	if err = DownloadSDK(cosContext); err != nil {
		return fmt.Errorf(T("unable to install Aspera SDK: %s"), err)
	}

	// Build UploadInput
	input := new(s3manager.UploadInput)

	// Required parameters for UploadInput
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	//
	// Optional parameters for UploadInput
	options := map[string]string{}

	//
	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	srcPath := c.Args().First()

	// Check if source path actually exists
	if _, err = cosContext.GetFileInfo(srcPath); err != nil {
		err = fmt.Errorf("%s: %s", err, srcPath)
		return
	}

	asp, err := cosContext.GetAsperaTransfer(APIKeyEnv, c.String(flags.Region))
	if err != nil {
		return
	}

	// Sync stop signal like CTRL+c
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	var bytes int64
	bytes, err = cosContext.GetTotalBytes(srcPath)
	transferInput := &aspera.COSInput{
		Bucket: aws.StringValue(input.Bucket),
		Key:    aws.StringValue(input.Key),
		Path:   srcPath,
		Sub:    aspera.NewProgressBarSubscriber(bytes, cosContext.UI.Writer()),
	}

	if err = asp.Upload(ctx, transferInput); err != nil {
		return
	}

	// displaying total uploaded bytes for json output
	output := &render.AsperaUploadOutput{TotalBytes: bytes}
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	return
}
