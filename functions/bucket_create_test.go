//+build unit

package functions_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
)

func TestBucketCreateSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	targetRegion := "us"
	targetClass := "standard"
	targetIbmServiceInstanceID := "IbmServiceInstanceID"

	var capturedInput *s3.CreateBucketInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilBucketExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("CreateBucket", mock.MatchedBy(
			func(input *s3.CreateBucketInput) bool {
				capturedInput = input
				return true
			})).
		Return(new(s3.CreateBucketOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketCreate,
		"--bucket", targetBucket,
		"--region", targetRegion,
		"--class", targetClass,
		"--ibm-service-instance-id", targetIbmServiceInstanceID,
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CreateBucket", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	// assert request match cli parameters
	assert.Equal(t, *capturedInput.Bucket, targetBucket)
	assert.Equal(t, *capturedInput.IBMServiceInstanceId, targetIbmServiceInstanceID)
	assert.Equal(t, *capturedInput.CreateBucketConfiguration.LocationConstraint, targetRegion+"-"+targetClass)

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")

}

func TestBucketCreateRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	targetRegion := "us"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilBucketExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("CreateBucket", mock.MatchedBy(
			func(input *s3.CreateBucketInput) bool {
				return true
			})).
		Return(nil, errors.New("Internal Server Error")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketCreate,
		"--bucket", targetBucket,
		"--region", targetRegion,
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CreateBucket", 1)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}

func TestBucketCreateWithoutBucket(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetRegion := "us"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilBucketExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("CreateBucket", mock.MatchedBy(
			func(input *s3.CreateBucketInput) bool {
				return true
			})).
		Return(nil, errors.New("Internal Server Error")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketCreate,
		"--region", targetRegion,
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CreateBucket", 0)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}
