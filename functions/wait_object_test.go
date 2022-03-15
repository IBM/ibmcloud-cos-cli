//+build unit

package functions_test

import (
	"os"
	"strconv"
	"testing"
	"time"

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

func TestWaitObjectNotExistsHappy(t *testing.T) {

	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "WaitObjBucket"
	targetKey := "WaitObjKey"
	targetIfMatch := "WaitObjIfMatch"
	targetIfModifiedSince, _ := time.Parse(time.RFC3339, "2001-01-01")
	targetIfNoneMatch := "WaitObjIfNoneMatch"
	targetIfUnmodifiedSince, _ := time.Parse(time.RFC3339, "2002-02-02")
	targetRange := "WaitObjRange"
	targetPartNumber := int64(42)

	var inputCapture *s3.HeadObjectInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On(
			"WaitUntilObjectNotExists",
			mock.MatchedBy(func(input *s3.HeadObjectInput) bool { inputCapture = input; return true })).
		Return(nil)

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Wait,
		commands.ObjectNotExists,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.IfMatch, targetIfMatch,
		"--" + flags.IfModifiedSince, targetIfModifiedSince.Format(time.RFC3339),
		"--" + flags.IfNoneMatch, targetIfNoneMatch,
		"--" + flags.IfUnmodifiedSince, targetIfUnmodifiedSince.Format(time.RFC3339),
		"--" + flags.Range, targetRange,
		"--" + flags.PartNumber, strconv.FormatInt(targetPartNumber, 10),
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "WaitUntilObjectNotExists", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetKey, *inputCapture.Key)
	assert.Equal(t, targetIfMatch, *inputCapture.IfMatch)
	assert.Equal(t, targetIfModifiedSince, *inputCapture.IfModifiedSince)
	assert.Equal(t, targetIfNoneMatch, *inputCapture.IfNoneMatch)
	assert.Equal(t, targetIfUnmodifiedSince, *inputCapture.IfUnmodifiedSince)
	assert.Equal(t, targetRange, *inputCapture.Range)
	assert.Equal(t, targetPartNumber, *inputCapture.PartNumber)

	// capture all output //
	errors := providers.FakeUI.Errors()
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestWaitObjectExistsHappy(t *testing.T) {

	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	// flags
	targetBucket := "WaitObjBucket"
	targetKey := "WaitObjKey"

	var inputCapture *s3.HeadObjectInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On(
			"WaitUntilObjectExists",
			mock.MatchedBy(func(input *s3.HeadObjectInput) bool { inputCapture = input; return true })).
		Return(nil)

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Wait,
		commands.ObjectExists,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "WaitUntilObjectExists", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetKey, *inputCapture.Key)

	// capture all output //
	errors := providers.FakeUI.Errors()
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}
