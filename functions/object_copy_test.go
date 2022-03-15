//+build unit

package functions_test

import (
	"errors"
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
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
)

func TestObjectCopySunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetCopySource := "SourceBucket/SourceKey"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.CopyObjectOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectCopy, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CopySource, targetCopySource,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 1)
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

func TestObjectCopyRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetCopySource := "SourceBucket/SourceKey"
	badKey := "NoSuchKey"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("NoSuchKey")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectCopy, "--bucket", targetBucket,
		"--" + flags.CopySource, targetCopySource,
		"--" + flags.Key, badKey,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 1)
	// assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	// assert Not Fail
	assert.Contains(t, errors, "FAIL")
}

func TestObjectCopyWithoutCopySource(t *testing.T) {
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
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("NoSuchKey")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectCopy, "--bucket", targetBucket,
		"--" + flags.Key, badKey,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 0)
	// assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	assert.Contains(t, errors, "'--copy-source' is missing")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
}

func TestObjectCopyWebsiteRedirectLocation(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetCopySource := "SourceBucket/SourceKey"
	targetWebsiteRedirectLocation := "https://cloud.ibm.com"

	var capturedInput *s3.CopyObjectInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.CopyObjectOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectCopy, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CopySource, targetCopySource,
		"--region", "REG",
		"--" + flags.WebsiteRedirectLocation, targetWebsiteRedirectLocation}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert request match cli parameters
	assert.Equal(t, targetBucket, *capturedInput.Bucket)
	assert.Equal(t, targetKey, *capturedInput.Key)
	assert.Equal(t, targetCopySource, *capturedInput.CopySource)
	assert.Equal(t, targetWebsiteRedirectLocation, *capturedInput.WebsiteRedirectLocation)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectCopyVersionIdText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetCopySourceVersionId := "srcId"
	targetCopySource := "SourceBucket/SourceKey?versionId=" + targetCopySourceVersionId
	targetVersionId := "destId"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.CopyObjectOutput).
			SetVersionId(targetVersionId).
			SetCopySourceVersionId(targetCopySourceVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectCopy, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CopySource, targetCopySource,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Copy Source Version ID: "+targetCopySourceVersionId)
	assert.Contains(t, output, "Version ID: "+targetVersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectCopyVersionIdJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetCopySourceVersionId := "srcId"
	targetCopySource := "SourceBucket/SourceKey?versionId=" + targetCopySourceVersionId
	targetVersionId := "destId"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.CopyObjectOutput).
			SetVersionId(targetVersionId).
			SetCopySourceVersionId(targetCopySourceVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectCopy, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CopySource, targetCopySource,
		"--output", "json",
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert json
	assert.Contains(t, output, "\"CopySourceVersionId\": \""+targetCopySourceVersionId+"\"")
	assert.Contains(t, output, "\"VersionId\": \""+targetVersionId+"\"")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectCopyTaggingDirectiveCopy(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetCopySource := "SourceBucket/SourceKey"
	targetTaggingDirective := "COPY"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.CopyObjectInput)
	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.CopyObjectOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectCopy, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CopySource, targetCopySource,
		"--tagging-directive", targetTaggingDirective,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Equal(t, aws.StringValue(capturedInput.TaggingDirective), targetTaggingDirective)
	assert.Nil(t, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectCopyTaggingDirectiveReplace(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetCopySource := "SourceBucket/SourceKey"
	targetTaggingDirective := "REPLACE"
	targetTagging := "key1=value1"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.CopyObjectInput)
	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.CopyObjectOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectCopy, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CopySource, targetCopySource,
		"--tagging-directive", targetTaggingDirective,
		"--tagging", targetTagging,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Equal(t, aws.StringValue(capturedInput.TaggingDirective), targetTaggingDirective)
	assert.Equal(t, aws.StringValue(capturedInput.Tagging), targetTagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}
