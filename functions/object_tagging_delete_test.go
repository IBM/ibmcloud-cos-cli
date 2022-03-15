//+build unit

package functions_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
)

func TestObjectTaggingDeleteRegularBucketText(t *testing.T) {
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

	var capturedInput = new(s3.DeleteObjectTaggingInput)
	providers.MockS3API.
		On("DeleteObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingDelete,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "Version ID")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingDeleteRegularBucketJson(t *testing.T) {
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

	var capturedInput = new(s3.DeleteObjectTaggingInput)
	providers.MockS3API.
		On("DeleteObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingDelete,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG",
		"--output", "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"VersionId\": null")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingDeleteVersionedBucketWithVersionIdText(t *testing.T) {
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

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.DeleteObjectTaggingInput)
	providers.MockS3API.
		On("DeleteObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingDelete,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--version-id", targetVersionId,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Version ID: "+targetVersionId)
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	assert.Equal(t, aws.StringValue(capturedInput.VersionId), targetVersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingDeleteVersionedBucketWithVersionIdJson(t *testing.T) {
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

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.DeleteObjectTaggingInput)
	providers.MockS3API.
		On("DeleteObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingDelete,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--version-id", targetVersionId,
		"--region", "REG",
		"--output", "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"VersionId\": \""+targetVersionId+"\"")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	assert.Equal(t, aws.StringValue(capturedInput.VersionId), targetVersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingDeleteVersionedBucketWithoutVersionIdText(t *testing.T) {
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

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.DeleteObjectTaggingInput)
	providers.MockS3API.
		On("DeleteObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingDelete,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Version ID: "+targetVersionId)
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	assert.Nil(t, capturedInput.VersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingDeleteVersionedBucketWithoutVersionIdJson(t *testing.T) {
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

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.DeleteObjectTaggingInput)
	providers.MockS3API.
		On("DeleteObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingDelete,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG",
		"--output", "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"VersionId\": \""+targetVersionId+"\"")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	assert.Nil(t, capturedInput.VersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingDeleteWithoutBucket(t *testing.T) {
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
		On("DeleteObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingDelete,
		"--key", targetKey,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjectTagging", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
}

func TestObjectTaggingDeleteWithoutKey(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("DeleteObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingDelete,
		"--bucket", targetBucket,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjectTagging", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
}
