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
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
)

func TestWaitBucketNotExistsHappy(t *testing.T) {

	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetWaitBucket"

	var captureInput *s3.HeadBucketInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On(
			"WaitUntilBucketNotExists",
			mock.MatchedBy(func(input *s3.HeadBucketInput) bool { captureInput = input; return true })).
		Return(nil)

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Wait,
		commands.BucketNotExists,
		"--" + flags.Bucket, targetBucket,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "WaitUntilBucketNotExists", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *captureInput.Bucket)

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")
}

func TestWaitBucketExistsHappy(t *testing.T) {

	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetWaitBucket"

	var captureInput *s3.HeadBucketInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On(
			"WaitUntilBucketExists",
			mock.MatchedBy(func(input *s3.HeadBucketInput) bool { captureInput = input; return true })).
		Return(nil)

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Wait,
		commands.BucketExists,
		"--" + flags.Bucket, targetBucket,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "WaitUntilBucketExists", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *captureInput.Bucket)

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")
}
