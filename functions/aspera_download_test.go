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
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/IBM/ibmcloud-cos-cli/functions"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

func TestAsperaDownload(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "aspera-bucket"
	targetKey := "aspera-key"

	// mock writing destination
	targetPath := "/mock/path"
	targetFileName := targetPath + "/MockFileName"

	// random string to mock file content
	objectContent := "followTheWhiteRabbit"
	// checks if file removed in the end
	isRemoved := false

	var inputCapture *s3.GetObjectInput

	downloadPath := "/mock/downloads"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	providers.MockPluginConfig.
		On(
			"GetStringWithDefault",
			mock.MatchedBy(func(key string) bool { return key == config.DownloadLocation }),
			mock.AnythingOfType("string"),
		).
		Return(downloadPath, nil)
	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us-east", nil)

	os.Setenv("IBMCLOUD_API_KEY", "apikey")
	os.Setenv("ASPERA_SDK_PATH", "/mock/aspera/sdk")
	transferdBin := "asperatransferd"
	if runtime.GOOS == "window" {
		transferdBin = "asperatransferd.exe"
	}
	transferBinPath := filepath.Join("/mock/aspera/sdk", "bin", transferdBin)

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == transferBinPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetFileName })).
		Return(os.Stat(targetFileName)) // some random name that must not exist

	providers.MockFileOperations.
		On("Remove", mock.MatchedBy(func(path string) bool { return path == targetFileName })).
		Run(func(args mock.Arguments) { isRemoved = true }).
		Return(nil)

	providers.MockS3API.
		On("GetObject", mock.MatchedBy(
			func(input *s3.GetObjectInput) bool {
				inputCapture = input
				return true
			})).
		Return(new(s3.GetObjectOutput).
			SetContentLength(int64(len(objectContent))), nil).
		Once()

	// Download is handled by asperatransferd deamon and ascp command,
	// so we can't mock the aspera file operation
	providers.MockAsperaTransfer.
		On("Download", mock.Anything, mock.Anything).
		Return(nil)

	// --- Act ----
	// set os args
	// all args for later verification
	os.Args = []string{
		"-",
		commands.AsperaDownload,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--region", "REG",
		targetFileName,
	}
	//call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all wroteContent //
	wroteContent := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()

	//assert OK
	assert.Contains(t, wroteContent, "OK")
	//assert Not Fail
	assert.NotContains(t, errors, "FAIL")

	// verify all parameters fwd
	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetKey, *inputCapture.Key)

	// assert file not removed
	assert.False(t, isRemoved, "Is Removed")
}

func TestAsperaDownloadOutputText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "aspera-bucket"
	targetKey := "aspera-key"

	// mock writing destination
	targetPath := "/mock/path"
	targetFileName := targetPath + "/MockFileName"

	// random string to mock file content
	objectContent := "followTheWhiteRabbit"

	downloadPath := "/mock/downloads"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	providers.MockPluginConfig.
		On(
			"GetStringWithDefault",
			mock.MatchedBy(func(key string) bool { return key == config.DownloadLocation }),
			mock.AnythingOfType("string"),
		).
		Return(downloadPath, nil)
	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us-east", nil)

	os.Setenv("IBMCLOUD_API_KEY", "apikey")
	os.Setenv("ASPERA_SDK_PATH", "/mock/aspera/sdk")
	transferdBin := "asperatransferd"
	if runtime.GOOS == "window" {
		transferdBin = "asperatransferd.exe"
	}
	transferBinPath := filepath.Join("/mock/aspera/sdk", "bin", transferdBin)

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == transferBinPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetFileName })).
		Return(os.Stat(targetFileName)) // some random name that must not exist

	providers.MockS3API.
		On("GetObject", mock.MatchedBy(
			func(input *s3.GetObjectInput) bool {
				return true
			})).
		Return(new(s3.GetObjectOutput).
			SetContentLength(int64(len(objectContent))), nil).
		Once()

	// Download is handled by asperatransferd deamon and ascp command,
	// so we can't mock the aspera file operation
	providers.MockAsperaTransfer.
		On("Download", mock.Anything, mock.Anything).
		Return(nil)

	// --- Act ----
	// set os args
	// all args for later verification
	os.Args = []string{
		"-",
		commands.AsperaDownload,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--region", "REG",
		targetFileName,
	}
	//call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all wroteContent //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()

	//assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Successfully downloaded 'aspera-key' from bucket 'aspera-bucket'")
	//assert Totalbytes in text output
	assert.Contains(t, output, "20 B downloaded.")
	//assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestAsperaDownloadOutputJSON(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "aspera-bucket"
	targetKey := "aspera-key"

	// mock writing destination
	targetPath := "/mock/path"
	targetFileName := targetPath + "/MockFileName"

	// random string to mock file content
	objectContent := "followTheWhiteRabbit"

	downloadPath := "/mock/downloads"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	providers.MockPluginConfig.
		On(
			"GetStringWithDefault",
			mock.MatchedBy(func(key string) bool { return key == config.DownloadLocation }),
			mock.AnythingOfType("string"),
		).
		Return(downloadPath, nil)
	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us-east", nil)

	os.Setenv("IBMCLOUD_API_KEY", "apikey")
	os.Setenv("ASPERA_SDK_PATH", "/mock/aspera/sdk")
	transferdBin := "asperatransferd"
	if runtime.GOOS == "window" {
		transferdBin = "asperatransferd.exe"
	}
	transferBinPath := filepath.Join("/mock/aspera/sdk", "bin", transferdBin)

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == transferBinPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetFileName })).
		Return(os.Stat(targetFileName)) // some random name that must not exist

	providers.MockS3API.
		On("GetObject", mock.MatchedBy(
			func(input *s3.GetObjectInput) bool {
				return true
			})).
		Return(new(s3.GetObjectOutput).
			SetContentLength(int64(len(objectContent))), nil).
		Once()

	// Download is handled by asperatransferd deamon and ascp command,
	// so we can't mock the aspera file operation
	providers.MockAsperaTransfer.
		On("Download", mock.Anything, mock.Anything).
		Return(nil)

	// --- Act ----
	// set os args
	// all args for later verification
	os.Args = []string{
		"-",
		commands.AsperaDownload,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--region", "REG",
		"--output", "json",
		targetFileName,
	}
	//call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all wroteContent //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()

	//assert TotalBytes in JSON output
	assert.Contains(t, output, `"TotalBytes": 20`)
	//assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestAsperaDownloadMissingKey(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "aspera-bucket"

	// mock writing destination
	targetPath := "/mock/path"
	targetFileName := targetPath + "/MockFileName"

	// random string to mock file content
	objectContent := "followTheWhiteRabbit"

	downloadPath := "/mock/downloads"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	providers.MockPluginConfig.
		On(
			"GetStringWithDefault",
			mock.MatchedBy(func(key string) bool { return key == config.DownloadLocation }),
			mock.AnythingOfType("string"),
		).
		Return(downloadPath, nil)
	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us-east", nil)

	os.Setenv("IBMCLOUD_API_KEY", "apikey")
	os.Setenv("ASPERA_SDK_PATH", "/mock/aspera/sdk")
	transferdBin := "asperatransferd"
	if runtime.GOOS == "window" {
		transferdBin = "asperatransferd.exe"
	}
	transferBinPath := filepath.Join("/mock/aspera/sdk", "bin", transferdBin)

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == transferBinPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetFileName })).
		Return(os.Stat(targetFileName)) // some random name that must not exist

	providers.MockS3API.
		On("GetObject", mock.MatchedBy(
			func(input *s3.GetObjectInput) bool {
				return true
			})).
		Return(new(s3.GetObjectOutput).
			SetContentLength(int64(len(objectContent))), nil).
		Once()

	// Download is handled by asperatransferd deamon and ascp command,
	// so we can't mock the aspera file operation
	providers.MockAsperaTransfer.
		On("Download", mock.Anything, mock.Anything).
		Return(nil)

	// --- Act ----
	// set os args
	// all args for later verification
	os.Args = []string{
		"-",
		commands.AsperaDownload,
		"--" + flags.Bucket, targetBucket,
		"--region", "REG",
		targetFileName,
	}
	//call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	//assert exit code is 1
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all wroteContent //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()

	//Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "'--key' is missing")
}

func TestGetObjectSize(t *testing.T) {
	defer providers.MocksRESET()
	key := "a-key"
	providers.MockS3API.
		On("GetObject", mock.MatchedBy(
			func(input *s3.GetObjectInput) bool {
				return *input.Key == key
			})).
		Return(new(s3.GetObjectOutput).
			SetContentLength(int64(1024)), nil)

	size, err := functions.GetObjectSize(providers.MockS3API, &s3.GetObjectInput{Key: &key})
	assert.Nil(t, err)
	assert.Equal(t, int64(1024), size)
}
