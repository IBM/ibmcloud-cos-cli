//+build unit

package functions_test

import (
	"os"
	"testing"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

var (
	objectRetentionJSONStr = `{
		"Mode": "COMPLIANCE",
		"RetainUntilDate": "2025-01-01T00:00:00.000Z"
	}`
		// replaced '==' instead of ':' - invalid json format
	objectRetentionMalformedJSONStr = `{"Mode"== "COMPLIANCE", "RetainUntilDate"=="2025-01-01T00:00:00.000Z"}`
	objectRetentionSimpleJSONStr = `
		Mode=COMPLIANCE,
		RetainUntilDate="2025-01-01T00:00:00.000Z"
	`
	retentionDate = "2025-01-01T00:00:00.000Z"
	retentionDateParse, _ = time.Parse(time.RFC3339, retentionDate)
	objectRetentionPut = new(s3.ObjectLockRetention).SetMode("COMPLIANCE").SetRetainUntilDate(retentionDateParse)
)


func TestObjectRetentionValidJSONString(t *testing.T) {
	defer providers.MocksRESET()
	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedObjectRetention *s3.ObjectLockRetention

	providers.MockS3API.
		On("PutObjectRetention", mock.MatchedBy(
			func(input *s3.PutObjectRetentionInput) bool {
				capturedObjectRetention = input.Retention
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectRetentionOutput), nil).
		Once()
	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectRetentionPut,
		"--bucket", targetBucket,
		"--key", "targetKey",
		"--region", "REG",
		"--retention", objectRetentionJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))
	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectRetention", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, objectRetentionPut, capturedObjectRetention)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectRetentionPutValidJSONFile(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	targetFileName := "fileMock"
	isClosed := false

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedObjectRetention *s3.ObjectLockRetention

	providers.MockS3API.
		On("PutObjectRetention", mock.MatchedBy(
			func(input *s3.PutObjectRetentionInput) bool {
				capturedObjectRetention = input.Retention
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectRetentionOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(objectRetentionJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectRetentionPut, "--bucket", targetBucket, "--key", "targetKey", "--region", "REG",
		"--retention", "file://" + targetFileName}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectRetention", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, objectRetentionPut, capturedObjectRetention)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectRetentionPutValidSimplifiedJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedObjectRetention *s3.ObjectLockRetention

	providers.MockS3API.
		On("PutObjectRetention", mock.MatchedBy(
			func(input *s3.PutObjectRetentionInput) bool {
				capturedObjectRetention = input.Retention
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectRetentionOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectRetentionPut, "--bucket", targetBucket, "--key", "targetKey", "--region", "REG",
		"--retention", objectRetentionSimpleJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectRetention", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, objectRetentionPut, capturedObjectRetention)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectRetentionPutWithoutBucket(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("PutObjectRetention", mock.MatchedBy(
			func(input *s3.PutObjectRetentionInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectRetentionOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectRetentionPut, "--key", "targetKey", "--region", "REG",
		"--retention", objectRetentionJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectRetention", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
}

func TestObjectRetentionPutWithoutObjectRetention(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("PutObjectRetention", mock.MatchedBy(
			func(input *s3.PutObjectRetentionInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectRetentionOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectRetentionPut, "--bucket", targetBucket,"--key", "targetKey", "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectRetention", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "Mandatory Flag '--retention' is missing")
}

func TestObjectRetentionPutWithMalformedJsonObjectRetention(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("PutObjectRetention", mock.MatchedBy(
			func(input *s3.PutObjectRetentionInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectRetentionOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectRetentionPut, "--bucket", targetBucket, "--key", "targetKey", "--region", "REG",
		"--retention", objectRetentionMalformedJSONStr} // replaced '==' instead of ':' - invalid json format
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectRetention", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "The value in flag '--retention' is invalid")
}
