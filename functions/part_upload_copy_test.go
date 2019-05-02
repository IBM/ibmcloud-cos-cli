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

func TestPartUploadCopySunnyPath(t *testing.T) {
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
	targetPartNumber := "1"
	targetUploadID := "80fds-afdasfa-s32141"

	providers.MockS3API.
		On("UploadPartCopy", mock.MatchedBy(
			func(input *s3.UploadPartCopyInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.UploadPartCopyOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.UploadPartCopy, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CopySource, targetCopySource,
		"--" + flags.PartNumber, targetPartNumber,
		"--" + flags.UploadID, targetUploadID,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "UploadPartCopy", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")

}

func TestPartUploadCopyRainyPath(t *testing.T) {
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
	targetPartNumber := "1"
	targetUploadID := "80fds-afdasfa-s32141"

	providers.MockS3API.
		On("UploadPartCopy", mock.MatchedBy(
			func(input *s3.UploadPartCopyInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("InvalidUploadID")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.UploadPartCopy, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CopySource, targetCopySource,
		"--" + flags.PartNumber, targetPartNumber,
		"--" + flags.UploadID, targetUploadID,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "UploadPartCopy", 1)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Not Fail
	assert.Contains(t, output, "FAIL")

}

func TestPartUploadCopyWithoutCopySource(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetPartNumber := "1"
	targetUploadID := "80fds-afdasfa-s32141"

	providers.MockS3API.
		On("UploadPartCopy", mock.MatchedBy(
			func(input *s3.UploadPartCopyInput) bool {
				return *input.Bucket == targetBucket

			})).
		Return(nil, errors.New("InvalidUploadID")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.UploadPartCopy, "--bucket", targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.PartNumber, targetPartNumber,
		"--" + flags.UploadID, targetUploadID,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "UploadPartCopy", 0)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.Contains(t, output, "Incorrect Usage.")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}
