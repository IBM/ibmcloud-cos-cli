//+build unit

package functions_test

import (
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

func TestBucketReplicationFailuresText(t *testing.T) {
	defer providers.MocksRESET()

	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	firstSyncBefore := "2025-01-01T00:00:00Z"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	firstSyncAttempted, _ := time.Parse(time.RFC3339, "2024-12-01T10:00:00Z")
	lastSyncAttempted, _ := time.Parse(time.RFC3339, "2024-12-02T11:00:00Z")

	var capturedInput = new(s3.ListBucketReplicationFailuresInput)
	providers.MockS3API.
		On("ListBucketReplicationFailures", mock.MatchedBy(
			func(input *s3.ListBucketReplicationFailuresInput) bool {
				capturedInput = input
				return aws.StringValue(input.Bucket) == targetBucket &&
					aws.StringValue(input.FirstSyncAttemptedBefore) == firstSyncBefore &&
					aws.Int64Value(input.MaxKeys) == 10 &&
					aws.StringValue(input.EncodingType) == "url" &&
					aws.StringValue(input.ContinuationToken) == "TOKEN"
			})).
		Return(new(s3.ListBucketReplicationFailuresOutput).
			SetListReplicationFailureResult(new(s3.ListReplicationFailureResult).
				SetName(targetBucket).
				SetIsTruncated(true).
				SetKeyCount(1).
				SetMaxKeys(10).
				SetNextContinuationToken("NEXTTOKEN").
				SetContents([]*s3.ObjectSyncAttempted{
					new(s3.ObjectSyncAttempted).
						SetKey("obj-1").
						SetVersionId("v1").
						SetSyncType("metadata").
						SetSyncFailureCause("network error").
						SetFirstSyncAttempted(firstSyncAttempted).
						SetLastSyncAttempted(lastSyncAttempted),
				})), nil).
		Once()

	os.Args = []string{"-", commands.BucketReplicationFailures,
		"--bucket", targetBucket,
		"--first-sync-attempted-before", firstSyncBefore,
		"--max-keys", "10",
		"--encoding-type", "url",
		"--starting-token", "TOKEN",
		"--region", "REG",
		"--output", "text"}
	plugin.Start(new(cos.Plugin))

	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketReplicationFailures", 1)
	assert.Equal(t, (*int)(nil), exitCode)
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "replication failure")
	assert.Contains(t, output, "obj-1")
	assert.Contains(t, output, "NEXTTOKEN")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketReplicationFailuresJSON(t *testing.T) {
	defer providers.MocksRESET()

	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	firstSyncAttempted, _ := time.Parse(time.RFC3339, "2024-12-01T10:00:00Z")
	lastSyncAttempted, _ := time.Parse(time.RFC3339, "2024-12-02T11:00:00Z")

	providers.MockS3API.
		On("ListBucketReplicationFailures", mock.MatchedBy(
			func(input *s3.ListBucketReplicationFailuresInput) bool {
				return aws.StringValue(input.Bucket) == targetBucket
			})).
		Return(new(s3.ListBucketReplicationFailuresOutput).
			SetListReplicationFailureResult(new(s3.ListReplicationFailureResult).
				SetName(targetBucket).
				SetIsTruncated(false).
				SetKeyCount(1).
				SetContents([]*s3.ObjectSyncAttempted{
					new(s3.ObjectSyncAttempted).
						SetKey("obj-1").
						SetVersionId("v1").
						SetSyncType("metadata").
						SetSyncFailureCause("network error").
						SetFirstSyncAttempted(firstSyncAttempted).
						SetLastSyncAttempted(lastSyncAttempted),
				})), nil).
		Once()

	os.Args = []string{"-", commands.BucketReplicationFailures,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "json"}
	plugin.Start(new(cos.Plugin))

	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketReplicationFailures", 1)
	assert.Equal(t, (*int)(nil), exitCode)
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "ListReplicationFailureResult")
	assert.Contains(t, output, "\"Key\": \"obj-1\"")
	assert.Contains(t, output, "\"VersionId\": \"v1\"")
	assert.NotContains(t, errors, "FAIL")
}

