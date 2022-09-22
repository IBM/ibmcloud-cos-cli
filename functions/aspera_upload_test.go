//go:build unit
// +build unit

package functions_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibmcloud-cos-cli/aspera"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
)

func TestAsperaUpload(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	srcFile := "/path/to/source"

	var inputCapture *s3manager.UploadInput

	os.Setenv("IBMCLOUD_API_KEY", "apikey")
	os.Setenv("ASPERA_SDK_PATH", "/mock/aspera/sdk")
	transferdBin := "asperatransferd"
	if runtime.GOOS == "window" {
		transferdBin = "asperatransferd.exe"
	}
	transferBinPath := filepath.Join("/mock/aspera/sdk", "bin", transferdBin)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == transferBinPath })).
		Return(os.Stat(os.TempDir()))
	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == srcFile })).
		Return(os.Stat(os.TempDir()))
	providers.MockFileOperations.
		On("GetTotalBytes", mock.MatchedBy(func(path string) bool { return path == srcFile })).
		Return(int64(1024), nil)

	providers.MockAsperaTransfer.
		On("Upload", mock.Anything, mock.MatchedBy(func(input *aspera.COSInput) bool {
			inputCapture = &s3manager.UploadInput{Bucket: &input.Bucket, Key: &input.Key}
			return true
		})).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{
		"-", commands.AsperaUpload,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--region", "REG",
		srcFile,
	}

	// call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	output := providers.FakeUI.Outputs()

	// assert all command line options where set
	assert.NotNil(t, inputCapture) // assert file name matched
	assert.Equal(t, targetBucket, aws.StringValue(inputCapture.Bucket))
	assert.Equal(t, targetKey, aws.StringValue(inputCapture.Key))

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

func TestAsperaUploadOutputText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	srcFile := "/path/to/source"

	os.Setenv("IBMCLOUD_API_KEY", "apikey")
	os.Setenv("ASPERA_SDK_PATH", "/mock/aspera/sdk")
	transferdBin := "asperatransferd"
	if runtime.GOOS == "window" {
		transferdBin = "asperatransferd.exe"
	}
	transferBinPath := filepath.Join("/mock/aspera/sdk", "bin", transferdBin)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == transferBinPath })).
		Return(os.Stat(os.TempDir()))
	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == srcFile })).
		Return(os.Stat(os.TempDir()))
	providers.MockFileOperations.
		On("GetTotalBytes", mock.MatchedBy(func(path string) bool { return path == srcFile })).
		Return(int64(1024), nil)

	providers.MockAsperaTransfer.
		On("Upload", mock.Anything, mock.MatchedBy(func(input *aspera.COSInput) bool {
			return true
		})).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{
		"-", commands.AsperaUpload,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--region", "REG",
		"--output", "text",
		srcFile,
	}

	// call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	output := providers.FakeUI.Outputs()

	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	errors := providers.FakeUI.Errors()
	//assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Successfully uploaded object 'TargetKey' to bucket 'TargetBucket'.")
	//assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestAsperaUploadOutputJSON(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	srcFile := "/path/to/source"

	os.Setenv("IBMCLOUD_API_KEY", "apikey")
	os.Setenv("ASPERA_SDK_PATH", "/mock/aspera/sdk")
	transferdBin := "asperatransferd"
	if runtime.GOOS == "window" {
		transferdBin = "asperatransferd.exe"
	}
	transferBinPath := filepath.Join("/mock/aspera/sdk", "bin", transferdBin)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == transferBinPath })).
		Return(os.Stat(os.TempDir()))
	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == srcFile })).
		Return(os.Stat(os.TempDir()))
	providers.MockFileOperations.
		On("GetTotalBytes", mock.MatchedBy(func(path string) bool { return path == srcFile })).
		Return(int64(1024), nil)

	providers.MockAsperaTransfer.
		On("Upload", mock.Anything, mock.MatchedBy(func(input *aspera.COSInput) bool {
			return true
		})).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{
		"-", commands.AsperaUpload,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--region", "REG",
		"--output", "json",
		srcFile,
	}

	// call  plugin
	plugin.Start(new(cos.Plugin))

	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	//assert TotalBytes in JSON output
	assert.Contains(t, output, `"TotalBytes": 1024`)
	//assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}
