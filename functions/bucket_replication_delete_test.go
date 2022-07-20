//+build unit

package functions_test

import (
	"os"
	"testing"

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

func TestBucketReplicationDeleteConfigurationText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockS3API.On("WaitUntilBucketNotExists", mock.Anything).Return(nil).Once()

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.DeleteBucketReplicationInput)
	providers.MockS3API.
		On("DeleteBucketReplication", mock.MatchedBy(
			func(input *s3.DeleteBucketReplicationInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteBucketReplicationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationDelete,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "text"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteBucketReplication", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketReplicationDeleteConfigurationJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockS3API.On("WaitUntilBucketNotExists", mock.Anything).Return(nil).Once()

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.DeleteBucketReplicationInput)
	providers.MockS3API.
		On("DeleteBucketReplication", mock.MatchedBy(
			func(input *s3.DeleteBucketReplicationInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteBucketReplicationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationDelete,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "json"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteBucketReplication", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "{}")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketReplicationDeleteWithoutBucket(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilBucketNotExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("DeleteBucketReplication", mock.MatchedBy(
			func(input *s3.DeleteBucketReplicationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteBucketOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationDelete,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteBucketReplication", 0)
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
