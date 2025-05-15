//go:build unit
// +build unit

package functions_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

var (
	retentionDateGet         = "2025-01-01T00:00:00.000Z"
	retentionDateParseGet, _ = time.Parse(time.RFC3339, retentionDateGet)
	objectRetentionObject    = new(s3.ObjectLockRetention).SetMode("COMPLIANCE").SetRetainUntilDate(retentionDateParseGet)
)

func TestObjectRetentionGetConfigurationText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.GetObjectRetentionInput)
	providers.MockS3API.
		On("GetObjectRetention", mock.MatchedBy(
			func(input *s3.GetObjectRetentionInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectRetentionOutput).
			SetRetention(objectRetentionObject),
			nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectRetentionGet,
		"--bucket", targetBucket,
		"--key", "targetKey",
		"--region", "REG",
		"--output", "text"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectRetention", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Retention")
	assert.Contains(t, output, "Mode: COMPLIANCE")
	assert.Contains(t, output, "Retain Until Date (UTC): Jan 01, 2025 at 00:00:00")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectRetentionGetJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.GetObjectRetentionInput)
	providers.MockS3API.
		On("GetObjectRetention", mock.MatchedBy(
			func(input *s3.GetObjectRetentionInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectRetentionOutput).
			SetRetention(objectRetentionObject),
			nil).
		Once()
	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectRetentionGet,
		"--bucket", targetBucket,
		"--key", "targetKey",
		"--region", "REG",
		"--output", "json"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectRetention", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"Retention\"")
	assert.Contains(t, output, "\"Mode\": \"COMPLIANCE\"")
	assert.Contains(t, output, "\"RetainUntilDate\": \"2025-01-01T00:00:00Z\"")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectRetentionGetWithoutBucket(t *testing.T) {
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
		On("GetObjectRetention", mock.MatchedBy(
			func(input *s3.GetObjectRetentionInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectRetentionOutput).
			SetRetention(objectRetentionObject),
			nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectRetentionGet, "--key", "targetKey",
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectRetention", 0)
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

func TestObjectRetentionGetNoConfiguration(t *testing.T) {
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
		On("GetObjectRetention", mock.MatchedBy(
			func(input *s3.GetObjectRetentionInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(nil, errors.New("NoSuchObjectRetention")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectRetentionGet,
		"--bucket", targetBucket,
		"--key", "targetKey",
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectRetention", 1)
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
