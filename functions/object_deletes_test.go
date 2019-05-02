//+build unit

package functions_test

import (
	"errors"
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

func TestObjectDeletesSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetDelete := "Objects=[{Key=a},{Key=b}],Quiet=false"

	providers.MockS3API.
		On("DeleteObjects", mock.MatchedBy(
			func(input *s3.DeleteObjectsInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectsOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.DeleteObjects, "--bucket", targetBucket,
		"--" + flags.Delete, targetDelete,
		"--" + flags.Region, "REG"}
	//call  plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjects", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")

}

func TestObjectDeletesRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	BadDelete := "Objects=[{Key=string},{Key=string}],Quiet=fale"

	providers.MockS3API.
		On("DeleteObjects", mock.MatchedBy(
			func(input *s3.DeleteObjectsInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectsOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.DeleteObjects, "--bucket", targetBucket,
		"--" + flags.Delete, BadDelete,
		"--" + flags.Region, "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjects", 0)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}

func TestObjectDeletesBadJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	BadDelete := "{\"Objects\":[{\"Key\":\"a\"}],\"Quiet\": fale}"

	providers.MockS3API.
		On("DeleteObjects", mock.MatchedBy(
			func(input *s3.DeleteObjectsInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.DeleteObjectsOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.DeleteObjects, "--bucket", targetBucket,
		"--" + flags.Delete, BadDelete,
		"--" + flags.Region, "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjects", 0)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}

func TestObjectDeletesWithoutDelete(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"

	providers.MockS3API.
		On("DeleteObjects", mock.MatchedBy(
			func(input *s3.DeleteObjectsInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("InvalidDelete")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.DeleteObjects, "--bucket", targetBucket,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "DeleteObjects", 0)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.Contains(t, output, "Incorrect Usage.")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}
