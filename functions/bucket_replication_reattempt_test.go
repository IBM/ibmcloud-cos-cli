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

func TestBucketReplicationReattempt(t *testing.T) {
	defer providers.MocksRESET()

	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.PutBucketReplicationReattemptInput)
	providers.MockS3API.
		On("PutBucketReplicationReattempt", mock.MatchedBy(
			func(input *s3.PutBucketReplicationReattemptInput) bool {
				capturedInput = input
				return aws.StringValue(input.Bucket) == targetBucket
			})).
		Return(new(s3.PutBucketReplicationReattemptOutput), nil).
		Once()

	os.Args = []string{"-", commands.BucketReplicationReattempt,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "text"}
	plugin.Start(new(cos.Plugin))

	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketReplicationReattempt", 1)
	assert.Equal(t, (*int)(nil), exitCode)
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.NotContains(t, errors, "FAIL")
}

