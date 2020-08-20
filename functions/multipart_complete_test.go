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
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
)

func TestMPUCompleteSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetUploadID := "TargetUploadID"
	targetMultipartUpload := "Parts=[{ETag=etag1,PartNumber=1},{ETag=etag2,PartNumber=2}]"
	targetMultipartUploadObject := new(s3.CompletedMultipartUpload).SetParts(
		[]*s3.CompletedPart{
			new(s3.CompletedPart).SetETag("etag1").SetPartNumber(1),
			new(s3.CompletedPart).SetETag("etag2").SetPartNumber(2),
		},
	)

	var capturedInput *s3.CompleteMultipartUploadInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("CompleteMultipartUpload", mock.MatchedBy(
			func(input *s3.CompleteMultipartUploadInput) bool {
				capturedInput = input
				return true
			})).
		Return(
			new(s3.CompleteMultipartUploadOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.MultipartUploadComplete,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.UploadID, targetUploadID,
		"--" + flags.MultipartUpload, targetMultipartUpload,
		"--" + flags.Region, "us",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CompleteMultipartUpload", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	// assert request match cli parameters
	assert.Equal(t, *capturedInput.Bucket, targetBucket)
	assert.Equal(t, *capturedInput.Key, targetKey)
	assert.Equal(t, *capturedInput.UploadId, targetUploadID)
	assert.Equal(t, *capturedInput.MultipartUpload, *targetMultipartUploadObject)

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")
}

func TestMPUCompleteRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetUploadID := "TargetUploadID"
	targetMultipartUpload := "Parts=[{ETag=etag1,PartNumber=1},{ETag=etag2,PartNumber=2}]"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("CompleteMultipartUpload", mock.MatchedBy(
			func(input *s3.CompleteMultipartUploadInput) bool {
				return true
			})).
		Return(
			nil, errors.New("InvalidUploadId")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.MultipartUploadComplete,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.UploadID, targetUploadID,
		"--" + flags.MultipartUpload, targetMultipartUpload,
		"--" + flags.Region, "us",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CompleteMultipartUpload", 1)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}

func TestMPUCompleteWithoutUploadID(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetMultipartUpload := "Parts=[{ETag=etag1,PartNumber=1},{ETag=etag2,PartNumber=2}]"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("CompleteMultipartUpload", mock.MatchedBy(
			func(input *s3.CompleteMultipartUploadInput) bool {
				return true
			})).
		Return(
			nil, errors.New("InvalidUploadId")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.MultipartUploadComplete,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.MultipartUpload, targetMultipartUpload,
		"--" + flags.Region, "us",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CompleteMultipartUpload", 0)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}
