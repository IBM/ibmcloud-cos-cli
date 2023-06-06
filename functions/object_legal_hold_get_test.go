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

func TestObjectLegalHoldGetConfigurationText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	var objectLegalHoldObject = new(s3.ObjectLockLegalHold).SetStatus("ON")

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.GetObjectLegalHoldInput)
	providers.MockS3API.
		On("GetObjectLegalHold", mock.MatchedBy(
			func(input *s3.GetObjectLegalHoldInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectLegalHoldOutput).
			SetLegalHold(objectLegalHoldObject),
				 nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLegalHoldGet,
		"--bucket", targetBucket,
		"--key", "foo",
		"--region", "REG",
		"--output", "text"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectLegalHold", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Legal Hold Status: ON")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLegalHoldGetJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	objectLegalHoldObject := new(s3.ObjectLockLegalHold).SetStatus("ON")

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedInput = new(s3.GetObjectLegalHoldInput)
	providers.MockS3API.
		On("GetObjectLegalHold", mock.MatchedBy(
			func(input *s3.GetObjectLegalHoldInput) bool {
				capturedInput = input
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectLegalHoldOutput).
			SetLegalHold(objectLegalHoldObject),
				 nil).
		Once()
	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLegalHoldGet,
		"--bucket", targetBucket,
		"--key", "foo",
		"--region", "REG",
		"--output", "json"}
	//call plugin
	plugin.Start(new(cos.Plugin))


	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectLegalHold", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "LegalHold")
	assert.Contains(t, output, "\"Status\": \"ON\"")
	assert.Equal(t, aws.StringValue(capturedInput.Bucket), targetBucket)
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLegalHoldGetWithoutBucket(t *testing.T) {
	defer providers.MocksRESET()
	objectLegalHoldObject := new(s3.ObjectLockLegalHold).SetStatus("ON")

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("GetObjectLegalHold", mock.MatchedBy(
			func(input *s3.GetObjectLegalHoldInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetObjectLegalHoldOutput).
			SetLegalHold(objectLegalHoldObject),
				 nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLegalHoldGet, "--key", "foo",
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectLegalHold", 0)
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

func TestObjectLegalHoldGetNoConfiguration(t *testing.T) {
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
		On("GetObjectLegalHold", mock.MatchedBy(
			func(input *s3.GetObjectLegalHoldInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(nil, errors.New("NoSuchObjectLegalHold")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLegalHoldGet,
		"--bucket", targetBucket,
		"--key", "foo",
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetObjectLegalHold", 1)
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
