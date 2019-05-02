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

func TestObjectCopySunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetCopySource := "SourceBucket/SourceKey"

	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.CopyObjectOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.CopyObject, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CopySource, targetCopySource,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")

}

func TestObjectCopyRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetCopySource := "SourceBucket/SourceKey"
	badKey := "NoSuchKey"

	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("NoSuchKey")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.CopyObject, "--bucket", targetBucket,
		"--" + flags.CopySource, targetCopySource,
		"--" + flags.Key, badKey,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 1)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Not Fail
	assert.Contains(t, output, "FAIL")

}

func TestObjectCopyWithoutCopySource(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	badKey := "NoSuchKey"

	providers.MockS3API.
		On("CopyObject", mock.MatchedBy(
			func(input *s3.CopyObjectInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("NoSuchKey")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.CopyObject, "--bucket", targetBucket,
		"--" + flags.Key, badKey,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CopyObject", 0)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.Contains(t, output, "Incorrect Usage.")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}
