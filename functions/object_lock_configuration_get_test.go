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

func TestObjectLockGetConfigurationText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	defaultRetention := new(s3.DefaultRetention).SetMode("COMPLIANCE").SetYears(1)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.GetObjectLockConfigurationInput)
	providers.MockS3API.
		On("GetObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.GetObjectLockConfigurationInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectLockConfigurationOutput).
			SetObjectLockConfiguration(new(s3.ObjectLockConfiguration).
				SetObjectLockEnabled("Enabled").
				SetRule(new(s3.ObjectLockRule).
				SetDefaultRetention(defaultRetention),
				)), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockGet,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "text"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectLockConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Object Lock Status: Enabled")
	assert.Contains(t, output, "Retention Mode: COMPLIANCE")
	assert.Contains(t, output, "Retention Period: 1 Years")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLockGetConfigurationJson(t *testing.T) {
	defer providers.MocksRESET()
	defaultRetention := new(s3.DefaultRetention).SetMode("COMPLIANCE").SetDays(1)

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.GetObjectLockConfigurationInput)
	providers.MockS3API.
		On("GetObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.GetObjectLockConfigurationInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectLockConfigurationOutput).
			SetObjectLockConfiguration(new(s3.ObjectLockConfiguration).
				SetObjectLockEnabled("Enabled").
				SetRule(new(s3.ObjectLockRule).
				SetDefaultRetention(defaultRetention),
				)), nil).
		Once()
	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockGet,
		"--bucket", targetBucket,
		"--region", "REG",
		"--output", "json"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectLockConfiguration", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "ObjectLockConfiguration")
	assert.Contains(t, output, "Rule")
	assert.Contains(t, output, "DefaultRetention")
	assert.Contains(t, output, "\"ObjectLockEnabled\": \"Enabled\"")
	assert.Contains(t, output, "\"Mode\": \"COMPLIANCE\"")
	assert.Contains(t, output, "\"Years\":")
	assert.Contains(t, output, "\"Days\": 1")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLockGetConfigurationWithoutBucket(t *testing.T) {
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
		On("GetObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.GetObjectLockConfigurationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectLockConfigurationOutput).
			SetObjectLockConfiguration(new(s3.ObjectLockConfiguration)), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockGet,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectLockConfiguration", 0)
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

func TestObjectLockGetConfigurationNoConfiguration(t *testing.T) {
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
		On("GetObjectLockConfiguration", mock.MatchedBy(
			func(input *s3.GetObjectLockConfigurationInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(nil, errors.New("NoSuchObjectLockConfiguration")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLockGet,
		"--bucket", targetBucket,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectLockConfiguration", 1)
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
