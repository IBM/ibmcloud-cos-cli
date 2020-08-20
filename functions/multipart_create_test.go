//+build unit

package functions_test

import (
	"errors"
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

func TestMPUCreateSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetCacheControl := "TargetCacheControl"
	targetContentDisposition := "TargetContentDisposition"
	targetContentEncoding := "TargetContentEncoding"
	targetContentLanguage := "TargetContentLanguage"
	targetContentType := "TargetContentType"
	targetMetadata := "key1=value1,key2=value2,key3=value3"

	targetMetadataAsMap := map[string]*string{
		"key1": aws.String("value1"),
		"key2": aws.String("value2"),
		"key3": aws.String("value3"),
	}

	targetUploadID := "TargetUploadID"

	var capturedInput *s3.CreateMultipartUploadInput

	providers.MockS3API.
		On("CreateMultipartUpload", mock.MatchedBy(
			func(input *s3.CreateMultipartUploadInput) bool {
				capturedInput = input
				return true
			})).
		Return(
			new(s3.CreateMultipartUploadOutput).
				SetBucket(targetBucket).
				SetKey(targetKey).
				SetUploadId(targetUploadID), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.MultipartUploadCreate,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CacheControl, targetCacheControl,
		"--" + flags.ContentDisposition, targetContentDisposition,
		"--" + flags.ContentEncoding, targetContentEncoding,
		"--" + flags.ContentLanguage, targetContentLanguage,
		"--" + flags.ContentType, targetContentType,
		"--" + flags.Metadata, targetMetadata,
		"--" + flags.Region, "us",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CreateMultipartUpload", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	// assert request match cli parameters
	assert.Equal(t, *capturedInput.Bucket, targetBucket)
	assert.Equal(t, *capturedInput.Key, targetKey)
	assert.Equal(t, *capturedInput.CacheControl, targetCacheControl)
	assert.Equal(t, *capturedInput.ContentDisposition, targetContentDisposition)
	assert.Equal(t, *capturedInput.ContentEncoding, targetContentEncoding)
	assert.Equal(t, *capturedInput.ContentLanguage, targetContentLanguage)
	assert.Equal(t, *capturedInput.ContentType, targetContentType)
	assert.Equal(t, capturedInput.Metadata, targetMetadataAsMap)

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")

	assert.Contains(t, output, targetBucket)
	assert.Contains(t, output, targetKey)
	assert.Contains(t, output, targetUploadID)

}

func TestMPUCreateRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetCacheControl := "TargetCacheControl"
	targetContentDisposition := "TargetContentDisposition"
	targetContentEncoding := "TargetContentEncoding"
	targetContentLanguage := "TargetContentLanguage"
	targetContentType := "TargetContentType"
	targetMetadata := "key1=value1,key2=value2,key3=value3"

	targetMetadataAsMap := map[string]*string{
		"key1": aws.String("value1"),
		"key2": aws.String("value2"),
		"key3": aws.String("value3"),
	}

	var capturedInput *s3.CreateMultipartUploadInput

	providers.MockS3API.
		On("CreateMultipartUpload", mock.MatchedBy(
			func(input *s3.CreateMultipartUploadInput) bool {
				capturedInput = input
				return true
			})).
		Return(
			nil, errors.New("NoSuchKey")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.MultipartUploadCreate,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.CacheControl, targetCacheControl,
		"--" + flags.ContentDisposition, targetContentDisposition,
		"--" + flags.ContentEncoding, targetContentEncoding,
		"--" + flags.ContentLanguage, targetContentLanguage,
		"--" + flags.ContentType, targetContentType,
		"--" + flags.Metadata, targetMetadata,
		"--" + flags.Region, "us",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CreateMultipartUpload", 1)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli

	// assert request match cli parameters
	assert.Equal(t, *capturedInput.Bucket, targetBucket)
	assert.Equal(t, *capturedInput.Key, targetKey)
	assert.Equal(t, *capturedInput.CacheControl, targetCacheControl)
	assert.Equal(t, *capturedInput.ContentDisposition, targetContentDisposition)
	assert.Equal(t, *capturedInput.ContentEncoding, targetContentEncoding)
	assert.Equal(t, *capturedInput.ContentLanguage, targetContentLanguage)
	assert.Equal(t, *capturedInput.ContentType, targetContentType)
	assert.Equal(t, capturedInput.Metadata, targetMetadataAsMap)

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}

func TestMPUCreateWithoutKey(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"

	providers.MockS3API.
		On("CreateMultipartUpload", mock.MatchedBy(
			func(input *s3.CreateMultipartUploadInput) bool {
				return true
			})).
		Return(
			nil, errors.New("NoSuchKey")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.MultipartUploadCreate,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Region, "us",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "CreateMultipartUpload", 0)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}
