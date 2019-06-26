//+build unit

package functions_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
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

func TestObjectGetSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "GetObjBucket"
	targetKey := "GetObjKey"
	targetIfMatch := "GetObjIfMatch"
	targetIfModifiedSince, _ := time.Parse(time.RFC3339, "2001-01-01")
	targetIfNoneMatch := "GetObjIfNoneMatch"
	targetIfUnmodifiedSince, _ := time.Parse(time.RFC3339, "2002-02-02")
	targetRange := "GetObjRange"
	targetResponseCacheControl := "GetObjResponseCacheControl"
	targetResponseContentDisposition := "GetObjResponseContentDisposition"
	targetResponseContentEncoding := "GetObjResponseContentEncoding"
	targetResponseContentLanguage := "GetObjResponseContentLanguage"
	targetResponseContentType := "GetObjResponseContentType"
	targetResponseExpires, _ := time.Parse(time.RFC3339, "2003-03-03")

	// mock writing destination
	targetPath := "/mock/path"
	targetFileName := targetPath + "/MockFileName"

	// string to collect the writes
	wroteFileContent := ""
	// checks if destination file was closed
	isClosed := false

	// random string to mock file content
	objectContent := "magicBytesUnicornsAndFairyDustPurpurin"
	// checks if file removed in the end
	isRemoved := false

	var inputCapture *s3.GetObjectInput

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == targetFileName })).
		Return(os.Stat(targetFileName)) // some random name that must not exist

	providers.MockFileOperations.
		On("WriteCloserOpen", mock.MatchedBy(func(path string) bool { return path == targetFileName })).
		Return(utils.WriteToString(&wroteFileContent, &isClosed), nil)

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
			SetContentLength(int64(len(objectContent))).
			SetBody(ioutil.NopCloser(strings.NewReader(objectContent))).
			SetLastModified(time.Now()), nil).
		Once()

	// --- Act ----
	// set os args
	// all args for later verification
	os.Args = []string{
		"-",
		commands.GetObject,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.IfMatch, targetIfMatch,
		"--" + flags.IfModifiedSince, targetIfModifiedSince.Format(time.RFC3339),
		"--" + flags.IfNoneMatch, targetIfNoneMatch,
		"--" + flags.IfUnmodifiedSince, targetIfUnmodifiedSince.Format(time.RFC3339),
		"--" + flags.Range, targetRange,
		"--" + flags.ResponseCacheControl, targetResponseCacheControl,
		"--" + flags.ResponseContentDisposition, targetResponseContentDisposition,
		"--" + flags.ResponseContentEncoding, targetResponseContentEncoding,
		"--" + flags.ResponseContentLanguage, targetResponseContentLanguage,
		"--" + flags.ResponseContentType, targetResponseContentType,
		"--" + flags.ResponseExpires, targetResponseExpires.Format(time.RFC3339),
		"--" + flags.Region, "REG",
		targetFileName,
	}
	//call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	// verify all parameters fwd
	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetKey, *inputCapture.Key)
	assert.Equal(t, targetIfMatch, *inputCapture.IfMatch)
	assert.Equal(t, targetIfModifiedSince, *inputCapture.IfModifiedSince)
	assert.Equal(t, targetIfNoneMatch, *inputCapture.IfNoneMatch)
	assert.Equal(t, targetIfUnmodifiedSince, *inputCapture.IfUnmodifiedSince)
	assert.Equal(t, targetRange, *inputCapture.Range)
	assert.Equal(t, targetResponseCacheControl, *inputCapture.ResponseCacheControl)
	assert.Equal(t, targetResponseContentDisposition, *inputCapture.ResponseContentDisposition)
	assert.Equal(t, targetResponseContentEncoding, *inputCapture.ResponseContentEncoding)
	assert.Equal(t, targetResponseContentLanguage, *inputCapture.ResponseContentLanguage)
	assert.Equal(t, targetResponseContentType, *inputCapture.ResponseContentType)
	assert.Equal(t, targetResponseExpires, *inputCapture.ResponseExpires)

	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "GetObject", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all wroteContent //
	wroteContent := providers.FakeUI.Outputs()
	//assert OK
	assert.Contains(t, wroteContent, "OK")
	//assert Not Fail
	assert.NotContains(t, wroteContent, "FAIL")

	// assert file closed
	assert.True(t, isClosed, "Is Closed")
	// assert file not removed
	assert.False(t, isRemoved, "Is Removed")

	assert.Equal(t, objectContent, wroteFileContent)
}

func TestObjectGetSunnyPath2(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "GetObjBucket"
	targetKey := "GetObjKey"

	// string to collect the writes
	wroteFileContent := ""
	// checks if destination file was closed
	isClosed := false

	// random string to mock file content
	objectContent := "magicBytesUnicornsAndFairyDustPurpurin"
	// checks if file removed in the end
	isRemoved := false

	var inputCapture *s3.GetObjectInput

	downloadPath := "/mock/downloads"
	fileName := downloadPath + "/" + targetKey

	providers.MockPluginConfig.
		On(
			"GetStringWithDefault",
			mock.MatchedBy(func(key string) bool { return key == config.DownloadLocation }),
			mock.AnythingOfType("string"),
		).
		Return(downloadPath, nil)

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == downloadPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == fileName })).
		Return(os.Stat(fileName)) // some random name that must not exist

	providers.MockFileOperations.
		On("WriteCloserOpen", mock.MatchedBy(func(path string) bool { return path == fileName })).
		Return(utils.WriteToString(&wroteFileContent, &isClosed), nil)

	providers.MockFileOperations.
		On("Remove", mock.MatchedBy(func(path string) bool { return path == fileName })).
		Run(func(args mock.Arguments) { isRemoved = true }).
		Return(nil)

	providers.MockS3API.
		On("GetObject", mock.MatchedBy(
			func(input *s3.GetObjectInput) bool {
				inputCapture = input
				return true
			})).
		Return(new(s3.GetObjectOutput).
			SetContentLength(int64(len(objectContent))).
			SetBody(ioutil.NopCloser(strings.NewReader(objectContent))).
			SetLastModified(time.Now()), nil).
		Once()

	// --- Act ----
	// set os args
	// all args for later verification
	os.Args = []string{
		"-",
		commands.GetObject,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.Region, "REG",
	}
	//call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	// verify all parameters fwd
	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetKey, *inputCapture.Key)

	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "GetObject", 1)

	providers.MockFileOperations.AssertNotCalled(t, "Remove", mock.Anything)

	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all wroteContent //
	wroteContent := providers.FakeUI.Outputs()
	//assert OK
	assert.Contains(t, wroteContent, "OK")
	//assert Not Fail
	assert.NotContains(t, wroteContent, "FAIL")

	// assert file closed
	assert.True(t, isClosed, "Is Closed")
	// assert file not removed
	assert.False(t, isRemoved, "Is Removed")

	assert.Equal(t, objectContent, wroteFileContent)

}

func TestObjectGetFileAlreadyExists(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "GetObjBucket"
	targetKey := "GetObjKey"

	downloadPath := "/mock/downloads"
	fileName := downloadPath + "/" + targetKey

	execName, _ := os.Executable()

	providers.FakeUI.Inputs("n")

	providers.MockPluginConfig.
		On(
			"GetStringWithDefault",
			mock.MatchedBy(func(key string) bool { return key == config.DownloadLocation }),
			mock.AnythingOfType("string"),
		).
		Return(downloadPath, nil)

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == downloadPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == fileName })).
		Return(os.Stat(execName))

	// --- Act ----
	// set os args
	// all args for later verification
	os.Args = []string{
		"-",
		commands.GetObject,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.Region, "REG",
	}
	//call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	providers.MockFileOperations.AssertNotCalled(t, "Remove", mock.Anything)

	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all wroteContent //
	wroteContent := providers.FakeUI.Outputs()
	//assert OK
	assert.NotContains(t, wroteContent, "OK")
	//assert Not Fail
	assert.NotContains(t, wroteContent, "FAIL")

}

func TestObjectGetDestinationIsDir(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "GetObjBucket"
	targetKey := "GetObjKey"

	downloadPath := "/mock/downloads"
	fileName := downloadPath + "/" + targetKey

	providers.MockPluginConfig.
		On(
			"GetStringWithDefault",
			mock.MatchedBy(func(key string) bool { return key == config.DownloadLocation }),
			mock.AnythingOfType("string"),
		).
		Return(downloadPath, nil)

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == downloadPath })).
		Return(os.Stat(os.TempDir()))

	providers.MockFileOperations.
		On("GetFileInfo", mock.MatchedBy(func(path string) bool { return path == fileName })).
		Return(os.Stat(os.TempDir()))

	// --- Act ----
	// set os args
	// all args for later verification
	os.Args = []string{
		"-",
		commands.GetObject,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.Region, "REG",
	}
	//call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	providers.MockFileOperations.AssertNotCalled(t, "Remove", mock.Anything)

	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all wroteContent //
	wroteContent := providers.FakeUI.Outputs()
	//assert OK
	assert.NotContains(t, wroteContent, "OK")
	//assert Not Fail
	assert.Contains(t, wroteContent, "FAIL")
}
