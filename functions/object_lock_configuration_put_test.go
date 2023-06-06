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
	objectLockConfigurationJSONStr = `{
		"ObjectLockEnabled": "Enabled",
		"Rule": {
				"DefaultRetention": {"Mode": "COMPLIANCE", "Years": 1}
		}
	}`
		// replaced '==' instead of ':' - invalid json format
	objectLockConfigurationMalformedJSONStr = `{
		"ObjectLockEnabled" == "Enabled",
		"Rule": {
				"DefaultRetention" == {"Mode" == "COMPLIANCE", "Years"== 1}
		}
	}`

	objectLockConfigurationSimpleJSONStr = `
	ObjectLockEnabled=Enabled,
	Rule={DefaultRetention={Mode=COMPLIANCE, Years=1}}
	`
	defaultRetention = new(s3.DefaultRetention).SetMode("Enabled").SetYears(1)

	objectLockConfigurationObject = new(s3.ObjectLockConfiguration).
						SetObjectLockEnabled("Enabled").
						SetRule(new(s3.ObjectLockRule).
							SetDefaultRetention(new(s3.DefaultRetention).
								SetMode("COMPLIANCE").
								SetYears(1),
							),
						)
)

func TestObjectLockPutValidJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedObjectLockConfiguration *s3.ObjectLockConfiguration

	providers.MockS3API.
		On("PutObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.PutObjectLockConfigurationInput) bool {
				capturedObjectLockConfiguration = input.ObjectLockConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLockConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockPut, "--bucket", targetBucket, "--region", "REG",
		"--object-lock-configuration", objectLockConfigurationJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLockConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, objectLockConfigurationObject, capturedObjectLockConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLockPutValidJSONFile(t *testing.T) {
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

	var capturedObjectLockConfiguration *s3.ObjectLockConfiguration

	providers.MockS3API.
		On("PutObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.PutObjectLockConfigurationInput) bool {
				capturedObjectLockConfiguration = input.ObjectLockConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLockConfigurationOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(objectLockConfigurationJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockPut, "--bucket", targetBucket, "--region", "REG",
		"--object-lock-configuration", "file://" + targetFileName}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLockConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, objectLockConfigurationObject, capturedObjectLockConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLockPutValidSimplifiedJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedObjectLockConfiguration *s3.ObjectLockConfiguration

	providers.MockS3API.
		On("PutObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.PutObjectLockConfigurationInput) bool {
				capturedObjectLockConfiguration = input.ObjectLockConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLockConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockPut, "--bucket", targetBucket, "--region", "REG",
		"--object-lock-configuration", objectLockConfigurationSimpleJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLockConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, objectLockConfigurationObject, capturedObjectLockConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLockPutWithoutBucket(t *testing.T) {
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
		On("PutObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.PutObjectLockConfigurationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLockConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockPut, "--region", "REG",
		"--object-lock-configuration", objectLockConfigurationJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLockConfiguration", 0)
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

func TestObjectLockPutWithoutObjectLockConfiguration(t *testing.T) {
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
		On("PutObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.PutObjectLockConfigurationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLockConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockPut, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLockConfiguration", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "Mandatory Flag '--object-lock-configuration' is missing")
}

func TestObjectLockPutWithMalformedJsonObjectLockConfiguration(t *testing.T) {
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
		On("PutObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.PutObjectLockConfigurationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLockConfigurationOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockPut, "--bucket", targetBucket, "--region", "REG",
		"--object-lock-configuration", objectLockConfigurationMalformedJSONStr} // replaced '==' instead of ':' - invalid json format
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLockConfiguration", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "The value in flag '--object-lock-configuration' is invalid")
}
