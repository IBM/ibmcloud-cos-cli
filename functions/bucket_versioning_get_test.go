//+build unit

package functions_test

import (
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

func TestBucketVersioningGetTextValidConfiguration(t *testing.T) {
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
		On("GetBucketVersioning", mock.MatchedBy(
			func(input *s3.GetBucketVersioningInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketVersioningOutput).
			SetStatus(s3.BucketVersioningStatusEnabled), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketVersioningGet, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketVersioning", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Status: Enabled")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketVersioningGetJsonValidConfiguration(t *testing.T) {
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
		On("GetBucketVersioning", mock.MatchedBy(
			func(input *s3.GetBucketVersioningInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketVersioningOutput).
			SetStatus(s3.BucketVersioningStatusEnabled), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketVersioningGet, "--bucket", targetBucket, "--region", "REG", "--output", "json"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketVersioning", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"Status\": \"Enabled\"")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketVersioningGetWithoutBucket(t *testing.T) {
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
		On("GetBucketVersioning", mock.MatchedBy(
			func(input *s3.GetBucketVersioningInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketVersioningOutput).
			SetStatus(s3.BucketVersioningStatusEnabled), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketVersioningGet, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketVersioning", 0)
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

func TestBucketVersioningGetNoVersioningConfiguration(t *testing.T) {
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
		On("GetBucketVersioning", mock.MatchedBy(
			func(input *s3.GetBucketVersioningInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketVersioningOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketVersioningGet, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketVersioning", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Versioning Configuration")
	assert.Contains(t, output, "(empty response from server; versioning has never been configured for this bucket)")
	assert.NotContains(t, output, "Status")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}
