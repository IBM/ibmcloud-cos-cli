//go:build unit
// +build unit

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

func TestBucketLifecycleGetConfigurationText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	lifecycleFilterObject = new(s3.LifecycleRuleFilter)
	expirationObject = new(s3.LifecycleExpiration).SetDays(30)
	abortMultipartUploadObject = new(s3.AbortIncompleteMultipartUpload).SetDaysAfterInitiation(30)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.GetBucketLifecycleConfigurationInput)
	providers.MockS3API.On("GetBucketLifecycleConfiguration", mock.MatchedBy(
		func(input *s3.GetBucketLifecycleConfigurationInput) bool {
			capturedInput = input
			return *input.Bucket == targetBucket
		})).
		Return(new(s3.GetBucketLifecycleConfigurationOutput).
			SetRules([]*s3.LifecycleRule{
				new(s3.LifecycleRule).
					SetStatus("Enabled").
					SetID("rule1").
					SetFilter(lifecycleFilterObject).
					SetExpiration(expirationObject).
					SetAbortIncompleteMultipartUpload(abortMultipartUploadObject),
			}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationGet,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "text"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketLifecycleConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Status: Enabled")
	assert.Contains(t, output, "Filter: Empty")
	assert.Contains(t, output, "Expiration in days: 30")
	assert.Contains(t, output, "Abort incomplete multipart upload initiated after, in days: 30")
	assert.Contains(t, output, "Found no noncurrentVersionTransitions in bucket lifecycle configuration")
	assert.Contains(t, output, "Found no transitions in bucket lifecycle configuration")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketLifecycleGetConfigurationJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	lifecycleFilterObject = new(s3.LifecycleRuleFilter).SetObjectSizeGreaterThan(1)
	expirationObject = new(s3.LifecycleExpiration).SetDays(30)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.GetBucketLifecycleConfigurationInput)
	providers.MockS3API.
		On("GetBucketLifecycleConfiguration", mock.MatchedBy(
			func(input *s3.GetBucketLifecycleConfigurationInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketLifecycleConfigurationOutput).
			SetRules([]*s3.LifecycleRule{
				new(s3.LifecycleRule).
					SetStatus("Enabled").
					SetID("rule1").
					SetFilter(lifecycleFilterObject).
					SetExpiration(expirationObject),
			}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationGet,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "json"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketLifecycleConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert
	assert.Contains(t, output, "Rules")
	assert.Contains(t, output, "Expiration")
	assert.Contains(t, output, "\"Days\": 30")
	assert.Contains(t, output, "Filter")
	assert.Contains(t, output, "\"ObjectSizeGreaterThan\": 1")
	assert.Contains(t, output, "\"ID\": \"rule1\"")
	assert.Contains(t, output, "\"Status\": \"Enabled\"")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketLifecycleGetConfigurationWithoutBucket(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	lifecycleFilterObject = new(s3.LifecycleRuleFilter).SetObjectSizeGreaterThan(1)
	expirationObject = new(s3.LifecycleExpiration).SetDays(30)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("GetBucketLifecycleConfiguration", mock.MatchedBy(
		func(input *s3.GetBucketLifecycleConfigurationInput) bool {
			return *input.Bucket == targetBucket
		})).
		Return(new(s3.GetBucketLifecycleConfigurationOutput).
			SetRules([]*s3.LifecycleRule{
				new(s3.LifecycleRule).
					SetStatus("Enabled").
					SetID("rule1").
					SetFilter(lifecycleFilterObject).
					SetExpiration(expirationObject),
			}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationGet,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketLifecycleConfiguration", 0)
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

func TestBucketLifecycleGetConfigurationNoConfiguration(t *testing.T) {
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
		On("GetBucketLifecycleConfiguration", mock.MatchedBy(
			func(input *s3.GetBucketLifecycleConfigurationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(nil, errors.New("InvalidParameter")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationGet,
		"--bucket", targetBucket,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketLifecycleConfiguration", 1)
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
