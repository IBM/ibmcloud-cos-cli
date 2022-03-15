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

func TestObjectTaggingGetRegularBucketText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetTagKey := "TargetTagA"   // The table has "Key" and "Value" in it, so don't double up here
	targetTagValue := "TargetTagB" // Otherwise if the header doesn't print, we may miss it

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.GetObjectTaggingInput)
	providers.MockS3API.
		On("GetObjectTagging", mock.MatchedBy(
			func(input *s3.GetObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectTaggingOutput).
			SetTagSet([]*s3.Tag{new(s3.Tag).
				SetKey(targetTagKey).
				SetValue(targetTagValue)}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingGet,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "Version ID")
	assert.Contains(t, output, "Key")
	assert.Contains(t, output, "Value")
	assert.Contains(t, output, targetTagKey)
	assert.Contains(t, output, targetTagValue)
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingGetRegularBucketJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetTagKey := "TargetTagA"   // The table has "Key" and "Value" in it, so don't double up here
	targetTagValue := "TargetTagB" // Otherwise if the header doesn't print, we may miss it

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.GetObjectTaggingInput)
	providers.MockS3API.
		On("GetObjectTagging", mock.MatchedBy(
			func(input *s3.GetObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectTaggingOutput).
			SetTagSet([]*s3.Tag{new(s3.Tag).
				SetKey(targetTagKey).
				SetValue(targetTagValue)}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingGet,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG",
		"--output", "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"TagSet\":")
	assert.Contains(t, output, "\"Key\": \""+targetTagKey+"\"")
	assert.Contains(t, output, "\"Value\": \""+targetTagValue+"\"")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingGetVersionedBucketWithVersionIdText(t *testing.T) {
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
	targetTagKey := "TargetTagA"   // The table has "Key" and "Value" in it, so don't double up here
	targetTagValue := "TargetTagB" // Otherwise if the header doesn't print, we may miss it

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.GetObjectTaggingInput)
	providers.MockS3API.
		On("GetObjectTagging", mock.MatchedBy(
			func(input *s3.GetObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectTaggingOutput).
			SetVersionId(targetVersionId).
			SetTagSet([]*s3.Tag{new(s3.Tag).
				SetKey(targetTagKey).
				SetValue(targetTagValue)}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingGet,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--version-id", targetVersionId,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Version ID: "+targetVersionId)
	assert.Contains(t, output, "Key")
	assert.Contains(t, output, "Value")
	assert.Contains(t, output, targetTagKey)
	assert.Contains(t, output, targetTagValue)
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	assert.Equal(t, aws.StringValue(capturedInput.VersionId), targetVersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingGetVersionedBucketWithVersionIdJson(t *testing.T) {
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
	targetTagKey := "TargetTagA"   // The table has "Key" and "Value" in it, so don't double up here
	targetTagValue := "TargetTagB" // Otherwise if the header doesn't print, we may miss it

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.GetObjectTaggingInput)
	providers.MockS3API.
		On("GetObjectTagging", mock.MatchedBy(
			func(input *s3.GetObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectTaggingOutput).
			SetVersionId(targetVersionId).
			SetTagSet([]*s3.Tag{new(s3.Tag).
				SetKey(targetTagKey).
				SetValue(targetTagValue)}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingGet,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--version-id", targetVersionId,
		"--region", "REG",
		"--output", "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"VersionId\": \""+targetVersionId+"\"")
	assert.Contains(t, output, "\"TagSet\":")
	assert.Contains(t, output, "\"Key\": \""+targetTagKey+"\"")
	assert.Contains(t, output, "\"Value\": \""+targetTagValue+"\"")
	assert.Contains(t, output, "\"Value\": \""+targetTagValue+"\"")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	assert.Equal(t, aws.StringValue(capturedInput.VersionId), targetVersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingGetVersionedBucketWithoutVersionIdText(t *testing.T) {
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
	targetTagKey := "TargetTagA"   // The table has "Key" and "Value" in it, so don't double up here
	targetTagValue := "TargetTagB" // Otherwise if the header doesn't print, we may miss it

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.GetObjectTaggingInput)
	providers.MockS3API.
		On("GetObjectTagging", mock.MatchedBy(
			func(input *s3.GetObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectTaggingOutput).
			SetVersionId(targetVersionId).
			SetTagSet([]*s3.Tag{new(s3.Tag).
				SetKey(targetTagKey).
				SetValue(targetTagValue)}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingGet,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Version ID: "+targetVersionId)
	assert.Contains(t, output, "Key")
	assert.Contains(t, output, "Value")
	assert.Contains(t, output, targetTagKey)
	assert.Contains(t, output, targetTagValue)
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	assert.Nil(t, capturedInput.VersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingGetVersionedBucketWithoutVersionIdJson(t *testing.T) {
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
	targetTagKey := "TargetTagA"   // The table has "Key" and "Value" in it, so don't double up here
	targetTagValue := "TargetTagB" // Otherwise if the header doesn't print, we may miss it

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.GetObjectTaggingInput)
	providers.MockS3API.
		On("GetObjectTagging", mock.MatchedBy(
			func(input *s3.GetObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectTaggingOutput).
			SetVersionId(targetVersionId).
			SetTagSet([]*s3.Tag{new(s3.Tag).
				SetKey(targetTagKey).
				SetValue(targetTagValue)}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingGet,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG",
		"--output", "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"VersionId\": \""+targetVersionId+"\"")
	assert.Contains(t, output, "\"TagSet\":")
	assert.Contains(t, output, "\"Key\": \""+targetTagKey+"\"")
	assert.Contains(t, output, "\"Value\": \""+targetTagValue+"\"")
	assert.Contains(t, output, "\"Value\": \""+targetTagValue+"\"")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	assert.Nil(t, capturedInput.VersionId)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingGetNoTagsetReturned(t *testing.T) {
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

	var capturedInput = new(s3.GetObjectTaggingInput)
	providers.MockS3API.
		On("GetObjectTagging", mock.MatchedBy(
			func(input *s3.GetObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingGet,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "Version ID")
	assert.NotContains(t, output, "Key")
	assert.NotContains(t, output, "Value")
	assert.Contains(t, output, "No tags returned")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingGetWithoutBucket(t *testing.T) {
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
		On("GetObjectTagging", mock.MatchedBy(
			func(input *s3.GetObjectTaggingInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingGet,
		"--key", targetKey,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectTagging", 0)
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

func TestObjectTaggingGetWithoutKey(t *testing.T) {
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
		On("GetObjectTagging", mock.MatchedBy(
			func(input *s3.GetObjectTaggingInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingGet,
		"--key", targetKey,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectTagging", 0)
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
