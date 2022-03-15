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
	publicAccessBlockConfigurationJSONStr = `{
  "BlockPublicAcls": true,
  "IgnorePublicAcls": false
}`

	publicAccessBlockConfigurationSimpleJSONStr = `
  BlockPublicAcls=true,
  IgnorePublicAcls=false
`

	publicAccessBlockConfigurationObject = new(s3.PublicAccessBlockConfiguration).
						SetBlockPublicAcls(true).
						SetIgnorePublicAcls(false)
)

func TestPublicAccessBlockPutValidJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedPublicAccessBlockConfiguration *s3.PublicAccessBlockConfiguration

	providers.MockS3API.
		On("PutPublicAccessBlock", mock.MatchedBy(
			func(input *s3.PutPublicAccessBlockInput) bool {
				capturedPublicAccessBlockConfiguration = input.PublicAccessBlockConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutPublicAccessBlockOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.PublicAccessBlockPut, "--bucket", targetBucket, "--region", "REG",
		"--public-access-block-configuration", publicAccessBlockConfigurationJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutPublicAccessBlock", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, publicAccessBlockConfigurationObject, capturedPublicAccessBlockConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestPublicAccessBlockPutValidJSONFile(t *testing.T) {
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

	var capturedPublicAccessBlockConfiguration *s3.PublicAccessBlockConfiguration

	providers.MockS3API.
		On("PutPublicAccessBlock", mock.MatchedBy(
			func(input *s3.PutPublicAccessBlockInput) bool {
				capturedPublicAccessBlockConfiguration = input.PublicAccessBlockConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutPublicAccessBlockOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(publicAccessBlockConfigurationJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.PublicAccessBlockPut, "--bucket", targetBucket, "--region", "REG",
		"--public-access-block-configuration", "file://" + targetFileName}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutPublicAccessBlock", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, publicAccessBlockConfigurationObject, capturedPublicAccessBlockConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestPublicAccessBlockPutValidSimplifiedJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedPublicAccessBlockConfiguration *s3.PublicAccessBlockConfiguration

	providers.MockS3API.
		On("PutPublicAccessBlock", mock.MatchedBy(
			func(input *s3.PutPublicAccessBlockInput) bool {
				capturedPublicAccessBlockConfiguration = input.PublicAccessBlockConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutPublicAccessBlockOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.PublicAccessBlockPut, "--bucket", targetBucket, "--region", "REG",
		"--public-access-block-configuration", publicAccessBlockConfigurationSimpleJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutPublicAccessBlock", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, publicAccessBlockConfigurationObject, capturedPublicAccessBlockConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestPublicAccessBlockPutWithoutBucket(t *testing.T) {
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
		On("PutPublicAccessBlock", mock.MatchedBy(
			func(input *s3.PutPublicAccessBlockInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutPublicAccessBlockOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.PublicAccessBlockPut, "--region", "REG",
		"--public-access-block-configuration", publicAccessBlockConfigurationJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutPublicAccessBlock", 0)
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

func TestPublicAccessBlockPutWithoutPublicAccessBlockConfiguration(t *testing.T) {
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
		On("PutPublicAccessBlock", mock.MatchedBy(
			func(input *s3.PutPublicAccessBlockInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutPublicAccessBlockOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.PublicAccessBlockPut, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutPublicAccessBlock", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "Mandatory Flag '--public-access-block-configuration' is missing")
}

func TestPublicAccessBlockPutWithMalformedJsonPublicAccessBlockConfiguration(t *testing.T) {
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
		On("PutPublicAccessBlock", mock.MatchedBy(
			func(input *s3.PutPublicAccessBlockInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutPublicAccessBlockOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.PublicAccessBlockPut, "--bucket", targetBucket, "--region", "REG",
		"--public-access-block-configuration", "{\"BlockPublicAcls\": true, \"IgnorePublicAcls\": false,}"} // trailing comma invalid
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutPublicAccessBlock", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "The value in flag '--public-access-block-configuration' is invalid")
}
