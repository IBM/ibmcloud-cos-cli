//+build unit

package functions_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws"
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

func TestObjectDeleteSunnyPath(t *testing.T) {
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

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("DeleteObject", mock.MatchedBy(
			func(input *s3.DeleteObjectInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectOutput), nil).
		Once()

	providers.FakeUI.Inputs("Y")

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectDelete, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.Region, "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, fmt.Sprintf("Delete '%s' from bucket '%s' ran successfully.", targetKey, targetBucket))
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectDeleteSunnyPathForce(t *testing.T) {
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

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("DeleteObject", mock.MatchedBy(
			func(input *s3.DeleteObjectInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectDelete, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.Region, "REG",
		"--" + flags.Force}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, fmt.Sprintf("Delete '%s' from bucket '%s' ran successfully.", targetKey, targetBucket))
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectDeleteRainyPath(t *testing.T) {
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
		On("DeleteObject", mock.MatchedBy(
			func(input *s3.DeleteObjectInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("NoSuchKey")).
		Once()

	providers.FakeUI.Inputs("Y")

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectDelete, "--bucket", targetBucket,
		"--" + flags.Key, badKey,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObject", 1)
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

func TestObjectDeleteWithoutKey(t *testing.T) {
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
		On("DeleteObject", mock.MatchedBy(
			func(input *s3.DeleteObjectInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("NoSuchKey")).
		Once()

	providers.FakeUI.Inputs("Y")

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectDelete, "--bucket", targetBucket,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObject", 0)
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

func TestObjectDeleteVersionIdDeleteMarker(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetVersionId := "TargetVersionId"

	var capturedDeleteObjectInput *s3.DeleteObjectInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("DeleteObject", mock.MatchedBy(
			func(input *s3.DeleteObjectInput) bool {
				capturedDeleteObjectInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectOutput).
			SetVersionId(targetVersionId).
			SetDeleteMarker(true), nil).
		Once()

	providers.FakeUI.Inputs("Y")

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectDelete, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.Region, "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, targetVersionId, aws.StringValue(capturedDeleteObjectInput.VersionId))
	assert.Contains(t, output, "Delete marker created for")
	assert.Contains(t, output, "Delete marker version ID:")
	assert.NotContains(t, output, "with version ID")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectDeleteVersionIdTargetedDelete(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetVersionId := "TargetVersionId"

	var capturedDeleteObjectInput *s3.DeleteObjectInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("DeleteObject", mock.MatchedBy(
			func(input *s3.DeleteObjectInput) bool {
				capturedDeleteObjectInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectOutput).
			SetVersionId("TargetVersionId"), nil).
		Once()

	providers.FakeUI.Inputs("Y")

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectDelete, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--version-id", targetVersionId,
		"--" + flags.Region, "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, targetVersionId, aws.StringValue(capturedDeleteObjectInput.VersionId))
	assert.Contains(t, output, "with version ID")
	assert.NotContains(t, output, "Delete marker created for")
	assert.NotContains(t, output, "Delete marker version ID:")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}
