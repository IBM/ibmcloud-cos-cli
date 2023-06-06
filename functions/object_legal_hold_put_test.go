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
	objectLegalHoldJSONStr = `{"Status": "ON"}`
		// replaced '==' instead of ':' - invalid json format
	objectLegalHoldMalformedJSONStr = `{"Status" == "ON"}`

	objectLegalHoldSimpleJSONStr = `Status=ON`
	objectLegalHoldObject = new(s3.ObjectLockLegalHold).SetStatus("ON")

)
func TestObjectLegalHoldValidJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedObjectLegalHold *s3.ObjectLockLegalHold

	providers.MockS3API.
		On("PutObjectLegalHold", mock.MatchedBy(
			func(input *s3.PutObjectLegalHoldInput) bool {
				capturedObjectLegalHold = input.LegalHold
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLegalHoldOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectLegalHoldPut,
		"--bucket", targetBucket,
		"--key", "targetKey",
		"--region", "REG",
		"--legal-hold", objectLegalHoldJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLegalHold", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, objectLegalHoldObject, capturedObjectLegalHold)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLegalHoldPutValidJSONFile(t *testing.T) {
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

	var capturedObjectLegalHold *s3.ObjectLockLegalHold

	providers.MockS3API.
		On("PutObjectLegalHold", mock.MatchedBy(
			func(input *s3.PutObjectLegalHoldInput) bool {
				capturedObjectLegalHold = input.LegalHold
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLegalHoldOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(objectLegalHoldJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLegalHoldPut, "--bucket", targetBucket, "--key", "targetKey", "--region", "REG",
		"--legal-hold", "file://" + targetFileName}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLegalHold", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, objectLegalHoldObject, capturedObjectLegalHold)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLegalHoldPutValidSimplifiedJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedObjectLegalHold *s3.ObjectLockLegalHold

	providers.MockS3API.
		On("PutObjectLegalHold", mock.MatchedBy(
			func(input *s3.PutObjectLegalHoldInput) bool {
				capturedObjectLegalHold = input.LegalHold
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLegalHoldOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLegalHoldPut, "--bucket", targetBucket, "--key", "targetKey", "--region", "REG",
		"--legal-hold", objectLegalHoldSimpleJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLegalHold", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, objectLegalHoldObject, capturedObjectLegalHold)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestObjectLegalHoldPutWithoutBucket(t *testing.T) {
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
		On("PutObjectLegalHold", mock.MatchedBy(
			func(input *s3.PutObjectLegalHoldInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLegalHoldOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLegalHoldPut, "--key", "targetKey", "--region", "REG",
		"--legal-hold", objectLegalHoldJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLegalHold", 0)
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

func TestObjectLockPutWithoutObjectLegalHold(t *testing.T) {
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
		On("PutObjectLegalHold", mock.MatchedBy(
			func(input *s3.PutObjectLegalHoldInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLegalHoldOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLegalHoldPut, "--bucket", targetBucket,"--key", "targetKey", "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLegalHold", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "Mandatory Flag '--legal-hold' is missing")
}

func TestObjectLockPutWithMalformedJsonObjectLegalHold(t *testing.T) {
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
		On("PutObjectLegalHold", mock.MatchedBy(
			func(input *s3.PutObjectLegalHoldInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutObjectLegalHoldOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.ObjectLegalHoldPut, "--bucket", targetBucket, "--key", "targetKey", "--region", "REG",
		"--legal-hold", objectLegalHoldMalformedJSONStr} // replaced '==' instead of ':' - invalid json format
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutObjectLegalHold", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "The value in flag '--legal-hold' is invalid")
}
