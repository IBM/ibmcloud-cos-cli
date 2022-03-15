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
	"github.com/IBM/ibmcloud-cos-cli/utils"
)

var (
	taggingJSONStr = `{
  "TagSet": [
    {
      "Key": "category",
      "Value": "storage"
    }
  ]
}`

	taggingSimpleJSONStr = `
  TagSet=[{
    Key=category,
    Value=storage
  }]
`

	taggingObject = new(s3.Tagging).SetTagSet([]*s3.Tag{
		new(s3.Tag).
			SetKey("category").
			SetValue("storage"),
	})
)

func TestObjectTaggingPutRegularBucketText(t *testing.T) {
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

	var capturedInput = new(s3.PutObjectTaggingInput)
	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.PutObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--tagging", taggingJSONStr,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 1)
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
	assert.Equal(t, taggingObject, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingPutRegularBucketJson(t *testing.T) {
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

	var capturedInput = new(s3.PutObjectTaggingInput)
	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.PutObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--tagging", taggingJSONStr,
		"--region", "REG",
		"--output", "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"VersionId\": null")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.Equal(t, aws.StringValue(capturedInput.Key), targetKey)
	assert.Equal(t, taggingObject, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingPutRegularBucketValidJsonString(t *testing.T) {
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

	var capturedInput = new(s3.PutObjectTaggingInput)
	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.PutObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--tagging", taggingJSONStr,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 1)
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
	assert.Equal(t, taggingObject, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingPutRegularBucketValidJsonFile(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetFileName := "fileMock"
	isClosed := false

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	var capturedInput = new(s3.PutObjectTaggingInput)
	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.PutObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectTaggingOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(taggingJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--tagging", taggingJSONStr,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 1)
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
	assert.Equal(t, taggingObject, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingPutRegularBucketValidSimplifiedJsonString(t *testing.T) {
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

	var capturedInput = new(s3.PutObjectTaggingInput)
	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.PutObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--tagging", taggingJSONStr,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 1)
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
	assert.Equal(t, taggingObject, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingPutVersionedBucketWithVersionIdText(t *testing.T) {
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

	var capturedInput = new(s3.PutObjectTaggingInput)
	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.PutObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectTaggingOutput).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--version-id", targetVersionId,
		"--tagging", taggingJSONStr,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 1)
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
	assert.Equal(t, taggingObject, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingPutVersionedBucketWithVersionIdJson(t *testing.T) {
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

	var capturedInput = new(s3.PutObjectTaggingInput)
	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.PutObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectTaggingOutput).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--version-id", targetVersionId,
		"--tagging", taggingJSONStr,
		"--region", "REG",
		"--output", "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 1)
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
	assert.Equal(t, taggingObject, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingPutVersionedBucketWithoutVersionIdText(t *testing.T) {
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

	var capturedInput = new(s3.PutObjectTaggingInput)
	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.PutObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectTaggingOutput).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--tagging", taggingJSONStr,
		"--region", "REG",
		"--output", "text"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 1)
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
	assert.Equal(t, taggingObject, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingPutVersionedBucketWithoutVersionIdJson(t *testing.T) {
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

	var capturedInput = new(s3.PutObjectTaggingInput)
	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.PutObjectTaggingInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectTaggingOutput).
			SetVersionId(targetVersionId), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--tagging", taggingJSONStr,
		"--region", "REG",
		"--output", "json"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 1)
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
	assert.Equal(t, taggingObject, capturedInput.Tagging)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectTaggingPutWithoutBucket(t *testing.T) {
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
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--key", targetKey,
		"--tagging", taggingJSONStr,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 0)
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

func TestObjectTaggingPutWithoutKey(t *testing.T) {
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
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--tagging", taggingJSONStr,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 0)
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

func TestObjectTaggingPutWithoutTagging(t *testing.T) {
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
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 0)
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

func TestObjectTaggingPutMalformedJsonTagging(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	malformedJsonTagging := "{\"TagSet\": [{\"Key\": \"category\", \"Value\": \"storage\",]}" // Trailing comma

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilObjectNotExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("PutObjectTagging", mock.MatchedBy(
			func(input *s3.DeleteObjectTaggingInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectTaggingOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectTaggingPut,
		"--bucket", targetBucket,
		"--key", targetKey,
		"--tagging", malformedJsonTagging,
		"--region", "REG"}
	// call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectTagging", 0)
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
