//+build unit

package functions_test

import (
	"errors"
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

func TestBucketReplicationGetConfigurationText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	filterObject = new(s3.ReplicationRuleFilter)
	destinationObject = new(s3.Destination).SetBucket("mockCRN:bucket:TARGETBUCKET")
	deleteMarkerReplicationObject = new(s3.DeleteMarkerReplication).SetStatus("Enabled")

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.GetBucketReplicationInput)
	providers.MockS3API.
		On("GetBucketReplication", mock.MatchedBy(
			func(input *s3.GetBucketReplicationInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketReplicationOutput).
			SetReplicationConfiguration(new(s3.ReplicationConfiguration).
				SetRules([]*s3.ReplicationRule{
					new(s3.ReplicationRule).
						SetStatus("Enabled").
						SetPriority(1).
						SetFilter(filterObject).
						SetDestination(destinationObject).
						SetDeleteMarkerReplication(deleteMarkerReplicationObject),
				})), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationGet,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "text"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketReplication", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Status: Enabled")
	assert.Contains(t, output, "Priority: 1")
	assert.Contains(t, output, "Filter")
	assert.Contains(t, output, "Destination")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketReplicationGetConfigurationJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	filterObject = new(s3.ReplicationRuleFilter)
	destinationObject = new(s3.Destination).SetBucket("mockCRN:bucket:TARGETBUCKET")
	deleteMarkerReplicationObject = new(s3.DeleteMarkerReplication).SetStatus("Enabled")

	var capturedInput = new(s3.GetBucketReplicationInput)
	providers.MockS3API.
		On("GetBucketReplication", mock.MatchedBy(
			func(input *s3.GetBucketReplicationInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketReplicationOutput).
			SetReplicationConfiguration(new(s3.ReplicationConfiguration).
				SetRules([]*s3.ReplicationRule{
					new(s3.ReplicationRule).
						SetStatus("Enabled").
						SetPriority(1).
						SetFilter(filterObject).
						SetDestination(destinationObject).
						SetDeleteMarkerReplication(deleteMarkerReplicationObject),
				})), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationGet,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "json"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketReplication", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert
	assert.Contains(t, output, "ReplicationConfiguration")
	assert.Contains(t, output, "Rules")
	assert.Contains(t, output, "\"Status\": \"Enabled\"")
	assert.Contains(t, output, "\"Priority\": 1")
	assert.Contains(t, output, "Filter")
	assert.Contains(t, output, "Destination")
	assert.Contains(t, output, "DeleteMarkerReplication")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketReplicationGetConfigurationWithoutBucket(t *testing.T) {
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
		On("GetBucketReplication", mock.MatchedBy(
			func(input *s3.GetBucketReplicationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketReplicationOutput).
			SetReplicationConfiguration(new(s3.ReplicationConfiguration).
				SetRules([]*s3.ReplicationRule{
					new(s3.ReplicationRule).
						SetStatus("Enabled").
						SetPriority(1).
						SetFilter(filterObject).
						SetDestination(destinationObject).
						SetDeleteMarkerReplication(deleteMarkerReplicationObject),
				})), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationGet,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketReplication", 0)
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

func TestBucketReplicationGetConfigurationNoConfiguration(t *testing.T) {
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
		On("GetBucketReplication", mock.MatchedBy(
			func(input *s3.GetBucketReplicationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(nil, errors.New("InvalidParameter")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationGet,
		"--bucket", targetBucket,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketReplication", 1)
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
