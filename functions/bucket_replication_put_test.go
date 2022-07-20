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
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

var (
	replicationConfigurationJSONStr = `{
	"Rules": [{
				"Status": "Enabled",
				"Priority": 1,
				"Filter": {},
				"Destination": { "Bucket": "mockCRN:bucket:TARGETBUCKET" },
				"DeleteMarkerReplication": { "Status": "Enabled" }
		}]
	}`

	// replaced '==' instead of ':' - invalid json format
	replicationConfigurationMalformedJSONStr = `{
	"Rules": [{
				"Status" == "Enabled"
				"Priority" == 1
				"Filter" == {}
				"Destination" == { "Bucket": "mockCRN:bucket:TARGETBUCKET" }
				"DeleteMarkerReplication" == { "Status": "Enabled" }
		}]
	}`

	replicationConfigurationSimpleJSONStr = `
   Rules=[{
			Status=Enabled,
			Priority=1,
			Filter={},
			Destination={Bucket="mockCRN:bucket:TARGETBUCKET"},
			DeleteMarkerReplication={Status=Enabled}
	}]`

	filterObject                  = new(s3.ReplicationRuleFilter)
	destinationObject             = new(s3.Destination).SetBucket("mockCRN:bucket:TARGETBUCKET")
	deleteMarkerReplicationObject = new(s3.DeleteMarkerReplication).SetStatus("Enabled")

	replicationConfigurationObject = new(s3.ReplicationConfiguration).
					SetRules([]*s3.ReplicationRule{
			new(s3.ReplicationRule).
				SetStatus("Enabled").
				SetPriority(1).
				SetFilter(filterObject).
				SetDestination(destinationObject).
				SetDeleteMarkerReplication(deleteMarkerReplicationObject),
		})
)

func TestBucketReplicationPutValidJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedReplicationConfiguration *s3.ReplicationConfiguration

	providers.MockS3API.
		On("PutBucketReplication", mock.MatchedBy(
			func(input *s3.PutBucketReplicationInput) bool {
				capturedReplicationConfiguration = input.ReplicationConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketReplicationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationPut, "--bucket", targetBucket, "--region", "REG",
		"--replication-configuration", replicationConfigurationJSONStr}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketReplication", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, replicationConfigurationObject, capturedReplicationConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketReplicationPutValidJSONFile(t *testing.T) {
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

	var capturedReplicationConfiguration *s3.ReplicationConfiguration

	providers.MockS3API.
		On("PutBucketReplication", mock.MatchedBy(
			func(input *s3.PutBucketReplicationInput) bool {
				capturedReplicationConfiguration = input.ReplicationConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketReplicationOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(replicationConfigurationJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationPut, "--bucket", targetBucket, "--region", "REG",
		"--replication-configuration", "file://" + targetFileName}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketReplication", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, replicationConfigurationObject, capturedReplicationConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketReplicationPutValidSimplifiedJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedReplicationConfiguration *s3.ReplicationConfiguration

	providers.MockS3API.
		On("PutBucketReplication", mock.MatchedBy(
			func(input *s3.PutBucketReplicationInput) bool {
				capturedReplicationConfiguration = input.ReplicationConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketReplicationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationPut, "--bucket", targetBucket, "--region", "REG",
		"--replication-configuration", replicationConfigurationSimpleJSONStr}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketReplication", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, replicationConfigurationObject, capturedReplicationConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketReplicationPutWithoutBucket(t *testing.T) {
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
		On("PutBucketReplication", mock.MatchedBy(
			func(input *s3.PutBucketReplicationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketReplicationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationPut, "--region", "REG",
		"--replication-configuration", replicationConfigurationJSONStr}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketReplication", 0)
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

func TestBucketReplicationPutWithoutReplicationConfiguration(t *testing.T) {
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
		On("PutBucketReplication", mock.MatchedBy(
			func(input *s3.PutBucketReplicationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketReplicationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationPut, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "BucketReplication", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "Mandatory Flag '--replication-configuration' is missing")
}

func TestBucketReplicationPutWithMalformedJsonReplicationConfiguration(t *testing.T) {
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
		On("PutBucketReplication", mock.MatchedBy(
			func(input *s3.PutBucketReplicationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketReplicationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketReplicationPut, "--bucket", targetBucket, "--region", "REG",
		"--replication-configuration", replicationConfigurationMalformedJSONStr} // replaced '==' instead of ':' - invalid json format
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketReplication", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "The value in flag '--replication-configuration' is invalid")
}
