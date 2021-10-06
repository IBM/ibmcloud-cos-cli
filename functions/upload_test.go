//+build unit

package functions_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/IBM/ibm-cos-sdk-go/aws"

	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

func TestUploadSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetFile := "TargetFile"

	targetCacheControl := "TargetCacheControl"
	targetContentDisposition := "TargetContentDisposition"
	targetContentEncoding := "TargetContentEncoding"
	targetContentLanguage := "TargetContentLanguage"
	//targetContentLength := "TargetContentLength"
	targetContentMD5 := "TargetContentMD5"
	targetContentType := "TargetContentType"

	targetPartSize := int64(9999)
	targetConcurrency := 99
	targetLeavePartsOnErrors := true
	targetMaxUploadParts := 999

	isClosed := false

	referenceUploader := new(s3manager.Uploader)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFile
		})).
		Return(utils.WrapString("", &isClosed), nil).
		Once()

	var inputCapture *s3manager.UploadInput
	providers.MockUploaderAPI.
		On("Upload", mock.MatchedBy(func(input *s3manager.UploadInput) bool {
			inputCapture = input
			return true
		})).
		Return(nil, nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.Upload, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.File, targetFile,
		"--" + flags.Region, "REG",
		"--" + flags.CacheControl, targetCacheControl,
		"--" + flags.ContentDisposition, targetContentDisposition,
		"--" + flags.ContentEncoding, targetContentEncoding,
		"--" + flags.ContentLanguage, targetContentLanguage,
		"--" + flags.ContentMD5, targetContentMD5,
		"--" + flags.ContentType, targetContentType,

		"--" + flags.PartSize, strconv.FormatInt(targetPartSize, 10),
		"--" + flags.Concurrency, strconv.Itoa(targetConcurrency),
		"--" + flags.LeavePartsOnErrors,
		"--" + flags.MaxUploadParts, strconv.Itoa(targetMaxUploadParts),
	}

	// call  plugin
	providers.ReferenceUploader = referenceUploader
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	output := providers.FakeUI.Outputs()
	fmt.Println(output)

	// assert all command line options where set
	assert.NotNil(t, inputCapture) // assert file name matched
	assert.Equal(t, targetBucket, aws.StringValue(inputCapture.Bucket))
	assert.Equal(t, targetKey, aws.StringValue(inputCapture.Key))

	assert.Equal(t, targetCacheControl, aws.StringValue(inputCapture.CacheControl))
	assert.Equal(t, targetContentDisposition, aws.StringValue(inputCapture.ContentDisposition))
	assert.Equal(t, targetContentEncoding, aws.StringValue(inputCapture.ContentEncoding))
	assert.Equal(t, targetContentLanguage, aws.StringValue(inputCapture.ContentLanguage))
	assert.Equal(t, targetContentMD5, aws.StringValue(inputCapture.ContentMD5))
	assert.Equal(t, targetContentType, aws.StringValue(inputCapture.ContentType))

	assert.Equal(t, targetPartSize, referenceUploader.PartSize)
	assert.Equal(t, targetConcurrency, referenceUploader.Concurrency)
	assert.Equal(t, targetLeavePartsOnErrors, referenceUploader.LeavePartsOnError)
	assert.Equal(t, targetMaxUploadParts, referenceUploader.MaxUploadParts)

	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	//output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, errors, "FAIL")

}
