//+build unit

package functions_test

import (
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

func TestBucketWebsiteDeleteSunnyPath(t *testing.T) {
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
		On("DeleteBucketWebsite", mock.MatchedBy(
			func(input *s3.DeleteBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteBucketWebsiteOutput), nil).
		Once()

	providers.FakeUI.Inputs("Y")

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsiteDelete, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteBucketWebsite", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, errors, "FAIL")
	// Confirm warning prompt is shown when --force not present
	assert.Contains(t, output, "Are you sure you would like to continue?")
	assert.Contains(t, errors, "WARNING")
}

func TestBucketWebsiteDeleteSunnyPathForce(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockS3API.On("WaitUntilBucketNotExists", mock.Anything).Return(nil).Once()

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("DeleteBucketWebsite", mock.MatchedBy(
			func(input *s3.DeleteBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteBucketWebsiteOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsiteDelete, "--bucket", targetBucket, "--region", "REG", "--force"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteBucketWebsite", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")
	// Confirm warning prompt is not shown when --force present
	assert.NotContains(t, output, "Are you sure you would like to continue?")
	assert.NotContains(t, errors, "WARNING")
}

func TestBucketDeleteWebsiteWithoutBucket(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.On("WaitUntilBucketNotExists", mock.Anything).Return(nil).Once()

	providers.MockS3API.
		On("DeleteBucketWebsite", mock.MatchedBy(
			func(input *s3.DeleteBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteBucketOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsiteDelete,
		"--region", "REG",
		"--force"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteBucketWebsite", 0)
	//assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, errors, "FAIL")
}
