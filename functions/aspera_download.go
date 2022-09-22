package functions

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	sdk "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/downloader"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/file_helpers"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/aspera"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/render"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

func AsperaDownload(c *cli.Context) (err error) {

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

	// Monitor the file
	keepFile := false

	// Download location
	var dstPath string

	// In case of error removes incomplete downloads
	defer func() {
		if !keepFile && dstPath != "" {
			cosContext.Remove(dstPath)
		}
	}()

	// Build GetObjectInput
	input := new(s3.GetObjectInput)

	// Required parameters for GetObjectInput
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	//
	// Optional parameters for GetObjectInput
	options := map[string]string{}

	//
	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	dstPath = c.Args().First()

	if dstPath == "" {
		// For consistence with the behavior of download command,
		// maybe we can change it to the current working directory.
		var downloadLocation string
		if downloadLocation, err = cosContext.GetDownloadLocation(); err != nil {
			return
		}
		dstPath = filepath.Join(downloadLocation, filepath.Base(aws.StringValue(input.Key)))
	}

	asp, err := cosContext.GetAsperaTransfer(APIKeyEnv, c.String(flags.Region))
	if err != nil {
		return
	}

	// Sync stop signal like CTRL+c
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	var client s3iface.S3API
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		return
	}

	var bytes int64
	if bytes, err = GetTotalBytes(client, input); err != nil {
		return
	}
	transferInput := &aspera.COSInput{
		Bucket: aws.StringValue(input.Bucket),
		Key:    aws.StringValue(input.Key),
		Path:   dstPath,
		Sub:    aspera.NewProgressBarSubscriber(bytes, cosContext.UI.Writer()),
	}

	if err = asp.Download(ctx, transferInput); err != nil {
		return
	}

	keepFile = true

	output := &render.DownloadOutput{
		TotalBytes: bytes,
	}
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	return
}

func GetTotalBytes(s s3iface.S3API, input *s3.GetObjectInput) (bytes int64, err error) {
	if strings.HasSuffix(aws.StringValue(input.Key), "/") {
		return GetDirectoryObjectSize(s, input)
	}
	return GetObjectSize(s, input)
}

func GetDirectoryObjectSize(s s3iface.S3API, input *s3.GetObjectInput) (size int64, err error) {
	pageIterInput := &s3.ListObjectsInput{
		Bucket: input.Bucket,
		Prefix: input.Key,
	}

	var objectContents []*s3.Object
	err = s.ListObjectsPages(pageIterInput, func(p *s3.ListObjectsOutput, _ bool) bool {
		objectContents = append(objectContents, p.Contents...)
		return true
	})
	if err != nil {
		return
	}

	if len(objectContents) == 0 {
		return 0, fmt.Errorf(T("no such directory: %s"), aws.StringValue(input.Key))
	}

	for _, object := range objectContents {
		size += aws.Int64Value(object.Size)
	}

	return
}

func GetObjectSize(s s3iface.S3API, input *s3.GetObjectInput) (size int64, err error) {
	output, err := s.GetObject(input)
	if err != nil {
		return
	}
	size = aws.Int64Value(output.ContentLength)
	return
}

func DownloadSDK(c *utils.CosContext) (err error) {
	if _, err = c.GetFileInfo(aspera.TransferdBinPath()); err == nil {
		return
	}
	c.UI.Warn(render.MessageAsperaBinaryNotFound())
	downloadURL, platform, err := aspera.GetSDKDownloadURL()
	if err != nil {
		return
	}

	tempDir, err := ioutil.TempDir("", "AsperaSDKDownload")
	if err != nil {
		return
	}
	SDKDownloader := downloader.New(tempDir)
	SDKDownloader.ProxyReader = downloader.NewProgressBar(c.UI.Writer())
	defer SDKDownloader.RemoveDir()

	pkgPath, _, err := SDKDownloader.Download(downloadURL)
	if err != nil {
		return
	}

	extractCmd := exec.Command("tar", "-xf", pkgPath, "-C", tempDir)
	if runtime.GOOS == "windows" {
		// built-in command for Windows 10: https://ibm.biz/Bdf7e
		// TODO: write a unzip functin with stdlib
		exec.Command("Expand-Archive", "-Path", pkgPath, "-DestinationPath", tempDir)
	}

	if err = extractCmd.Run(); err != nil {
		return
	}

	baseFolder := filepath.Join(tempDir, platform)
	if err = file_helpers.CopyDir(baseFolder, aspera.SDKDir()); err != nil {
		return
	}
	return
}
