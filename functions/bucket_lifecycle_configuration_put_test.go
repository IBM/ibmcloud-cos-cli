//go:build unit
// +build unit

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
	lifecycleConfigurationJSONStr = `{
	"Rules": [{
				"Status": "Enabled",
				"ID": "rule1",
				"Filter": {
					"And": null,
					"ObjectSizeGreaterThan": 1,
					"ObjectSizeLessThan": null,
					"Prefix": null,
					"Tag": null
				},
				"Expiration": {
					"Date": null,
					"Days": 30,
					"ExpiredObjectDeleteMarker": null
				}
		}]
	}`

	// replaced '==' instead of ':' - invalid json format
	lifecycleConfigurationMalformedJSONStr = `{
	"Rules": [{
				"Status" == "Enabled"
				"ID" == "rule1",
				"Filter" == {}
		     }]
	}`

	lifecycleConfigurationSimpleJSONStr = `
   Rules=[{
			Status=Enabled,
			ID=rule1,
			Filter={ObjectSizeGreaterThan=1},
			Expiration={Days=30}
		 }]`

	lifecycleFilterObject      = new(s3.LifecycleRuleFilter).SetObjectSizeGreaterThan(1)
	expirationObject           = new(s3.LifecycleExpiration).SetDays(30)
	abortMultipartUploadObject = new(s3.AbortIncompleteMultipartUpload).SetDaysAfterInitiation(30)

	lifecycleConfigurationObject = new(s3.LifecycleConfiguration).SetRules([]*s3.LifecycleRule{
		new(s3.LifecycleRule).
			SetStatus("Enabled").
			SetID("rule1").
			SetFilter(lifecycleFilterObject).
			SetExpiration(expirationObject),
	})
)

func TestBucketLifecycleConfigurationPutValidJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedLifecycleConfiguration *s3.LifecycleConfiguration

	providers.MockS3API.
		On("PutBucketLifecycleConfiguration", mock.MatchedBy(
			func(input *s3.PutBucketLifecycleConfigurationInput) bool {
				capturedLifecycleConfiguration = input.LifecycleConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketLifecycleConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationPut, "--bucket", targetBucket, "--region", "REG",
		"--lifecycle-configuration", lifecycleConfigurationJSONStr}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketLifecycleConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, lifecycleConfigurationObject, capturedLifecycleConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketLifecycleConfigurationPutValidJSONFile(t *testing.T) {
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

	var capturedLifecycleConfiguration *s3.LifecycleConfiguration

	providers.MockS3API.
		On("PutBucketLifecycleConfiguration", mock.MatchedBy(
			func(input *s3.PutBucketLifecycleConfigurationInput) bool {
				capturedLifecycleConfiguration = input.LifecycleConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketLifecycleConfigurationOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(lifecycleConfigurationJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationPut, "--bucket", targetBucket, "--region", "REG",
		"--lifecycle-configuration", "file://" + targetFileName}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketLifecycleConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, lifecycleConfigurationObject, capturedLifecycleConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketLifecycleConfigurationPutValidSimplifiedJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedLifecycleConfiguration *s3.LifecycleConfiguration

	providers.MockS3API.
		On("PutBucketLifecycleConfiguration", mock.MatchedBy(
			func(input *s3.PutBucketLifecycleConfigurationInput) bool {
				capturedLifecycleConfiguration = input.LifecycleConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketLifecycleConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationPut, "--bucket", targetBucket, "--region", "REG",
		"--lifecycle-configuration", lifecycleConfigurationSimpleJSONStr}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketLifecycleConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, lifecycleConfigurationObject, capturedLifecycleConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketLifecycleConfigurationPutWithoutBucket(t *testing.T) {
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
		On("PutBucketLifecycleConfiguration", mock.MatchedBy(
			func(input *s3.PutBucketLifecycleConfigurationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketLifecycleConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationPut, "--region", "REG",
		"--lifecycle-configuration", lifecycleConfigurationJSONStr}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketLifecycleConfiguration", 0)
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

func TestBucketLifecycleConfigurationPutWithoutLifecycleConfiguration(t *testing.T) {
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
		On("PutBucketLifecycleConfiguration", mock.MatchedBy(
			func(input *s3.PutBucketLifecycleConfigurationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketLifecycleConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationPut, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "BucketLifecycleConfiguration", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "Mandatory Flag '--lifecycle-configuration' is missing")
}

func TestBucketLifecycleConfigurationPutWithMalformedJson(t *testing.T) {
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
		On("PutBucketLifecycleConfiguration", mock.MatchedBy(
			func(input *s3.PutBucketLifecycleConfigurationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketLifecycleConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketLifeCycleConfigurationPut, "--bucket", targetBucket, "--region", "REG",
		"--lifecycle-configuration", lifecycleConfigurationMalformedJSONStr} // replaced '==' instead of ':' - invalid json format
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketLifecycleConfiguration", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "The value in flag '--lifecycle-configuration' is invalid")
}
