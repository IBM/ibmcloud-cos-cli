//+build unit

package functions_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

func TestObjectHeadSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("HeadObject", mock.MatchedBy(
			func(input *s3.HeadObjectInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.HeadObjectOutput).
			SetContentLength(int64(1)).
			SetLastModified(time.Now()), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectHead, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.Region, "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "HeadObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectHeadRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	badKey := "NoSuchKey"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("HeadObject", mock.MatchedBy(
			func(input *s3.HeadObjectInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("NoSuchKey")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectHead, "--bucket", targetBucket,
		"--" + flags.Key, badKey,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "HeadObject", 1)
	// assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
}

func TestObjectHeadWithoutKey(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("HeadObject", mock.MatchedBy(
			func(input *s3.HeadObjectInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("InvalidUploadID")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectHead, "--bucket", targetBucket,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "HeadObject", 0)
	// assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	assert.Contains(t, errors, "--key' is missing")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
}

func TestObjectHeadWebsiteRedirectLocation(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	expectedWebsiteRedirectLocation := "https://cloud.ibm.com"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("HeadObject", mock.MatchedBy(
			func(input *s3.HeadObjectInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.HeadObjectOutput).
			SetContentLength(int64(1)).
			SetLastModified(time.Now()).
			SetWebsiteRedirectLocation(expectedWebsiteRedirectLocation), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectHead, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.Region, "REG",
		"--" + flags.Output, "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "HeadObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert expectedWebsiteRedirectLocation in JSON output
	assert.Contains(t, output, expectedWebsiteRedirectLocation)
	assert.Contains(t, output, "WebsiteRedirectLocation")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectHeadVersionIdText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetVersionId := "88024d62-d55f-4332-be26-9b65c4d73bc0"

	var inputCapture *s3.HeadObjectInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("HeadObject", mock.MatchedBy(
			func(input *s3.HeadObjectInput) bool {
				inputCapture = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.HeadObjectOutput).
			SetContentLength(int64(1)).
			SetLastModified(time.Now()).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectHead, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.VersionId, targetVersionId,
		"--" + flags.Region, "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	// verify all parameters fwd
	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetKey, *inputCapture.Key)
	assert.Equal(t, targetVersionId, *inputCapture.VersionId)

	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "HeadObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Version ID: "+targetVersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectHeadVersionIdJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetVersionId := "88024d62-d55f-4332-be26-9b65c4d73bc0"

	var inputCapture *s3.HeadObjectInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("HeadObject", mock.MatchedBy(
			func(input *s3.HeadObjectInput) bool {
				inputCapture = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.HeadObjectOutput).
			SetContentLength(int64(1)).
			SetLastModified(time.Now()).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectHead, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.VersionId, targetVersionId,
		"--" + flags.Output, "json",
		"--" + flags.Region, "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----

	// verify all parameters fwd
	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetKey, *inputCapture.Key)
	assert.Equal(t, targetVersionId, *inputCapture.VersionId)

	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "HeadObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"VersionId\": \""+targetVersionId+"\"")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}
