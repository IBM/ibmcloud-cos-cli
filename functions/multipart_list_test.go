//+build unit

package functions_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

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

func TestMultiPartListHappyPathNoUploads(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "MPULBucket"
	targetDelimiter := "MPULDelimiter"
	targetEncodingType := "MPULEncodingType"
	targetPrefix := "MPULPrefix"
	targetKeyMarker := "MPULKeyMarker"
	targetUploadIdMarker := "MPULUploadIdMarker"

	var inputCapture *s3.ListMultipartUploadsInput

	providers.MockS3API.
		On("ListMultipartUploadsPages",
			mock.MatchedBy(
				func(input *s3.ListMultipartUploadsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListMultipartUploadsOutput, last bool) bool)
			paggerSrvMock(inputCapture, funcCapture, 0, nil)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.MultipartUploads,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--" + flags.KeyMarker, targetKeyMarker,
		"--" + flags.UploadIDMarker, targetUploadIdMarker,
		"--region", "REG",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListMultipartUploadsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetDelimiter, *inputCapture.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCapture.EncodingType)
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)
	assert.Equal(t, targetKeyMarker, *inputCapture.KeyMarker)
	assert.Equal(t, targetUploadIdMarker, *inputCapture.UploadIdMarker)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found no multipart uploads in bucket '%s'.", targetBucket))

}

func TestMultiPartListWhenPageBiggerThanMaxRequestMax(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "MPULBucketT2"
	targetMaxUploads := int64(3)

	var inputCapture *s3.ListMultipartUploadsInput

	providers.MockS3API.
		On("ListMultipartUploadsPages",
			mock.MatchedBy(
				func(input *s3.ListMultipartUploadsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.MultipartUploads,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, strconv.FormatInt(targetMaxUploads, 10),
		"--" + flags.PageSize, "100",
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListMultipartUploadsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetMaxUploads, *inputCapture.MaxUploads)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")

}

func TestMultiPartListWhenPageSmallerThanMaxRequestPage(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "MPULBucketT3"
	targetMaxUploads := int64(100)

	var inputCapture *s3.ListMultipartUploadsInput

	providers.MockS3API.
		On("ListMultipartUploadsPages",
			mock.MatchedBy(
				func(input *s3.ListMultipartUploadsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.MultipartUploads,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, "1000",
		"--" + flags.PageSize, strconv.FormatInt(targetMaxUploads, 10),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListMultipartUploadsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetMaxUploads, *inputCapture.MaxUploads)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")

}

func TestMultiPartListPaginate(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "MPULBucketT4"
	targetPageSz := 20
	targetMaxItems := 75

	var inputCapture *s3.ListMultipartUploadsInput

	pagesSzCapture := make([]int, 0, 0)

	providers.MockS3API.
		On("ListMultipartUploadsPages",
			mock.MatchedBy(
				func(input *s3.ListMultipartUploadsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListMultipartUploadsOutput, last bool) bool)
			paggerSrvMock(inputCapture, funcCapture, 1000, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.MultipartUploads,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, strconv.Itoa(targetMaxItems),
		"--" + flags.PageSize, strconv.Itoa(targetPageSz),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListMultipartUploadsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)

	// requested 75 elements split in 20 size pages
	assert.Equal(t, []int{20, 20, 20, 15}, pagesSzCapture)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found %d multipart uploads in bucket '%s'", targetMaxItems,
		targetBucket))

}

func paggerSrvMock(input *s3.ListMultipartUploadsInput, backFeed func(*s3.ListMultipartUploadsOutput, bool) bool,
	hardLimit int, results *[]int) {
	keepGoing := true
	if hardLimit == 0 {
		parts := new(s3.ListMultipartUploadsOutput)
		parts.Bucket = input.Bucket
		backFeed(parts, hardLimit <= 0)
		return
	}
	for hardLimit > 0 && keepGoing {
		i := int(aws.Int64Value(input.MaxUploads))
		*results = append(*results, i)
		if i == 0 || i > hardLimit {
			i = hardLimit
		}
		parts := buildListMultipartUploadsOutput(i)
		hardLimit -= i
		parts.Bucket = input.Bucket
		keepGoing = backFeed(parts, hardLimit <= 0)
	}
}

func buildListMultipartUploadsOutput(howMany int) *s3.ListMultipartUploadsOutput {
	result := new(s3.ListMultipartUploadsOutput).SetUploads([]*s3.MultipartUpload{})
	for i := 0; i < howMany; i++ {
		mpu := new(s3.MultipartUpload).
			SetInitiated(time.Now()).
			SetUploadId("uploadid" + strconv.Itoa(i)).
			SetKey("key" + strconv.Itoa(i))
		result.Uploads = append(result.Uploads, mpu)
		result.NextKeyMarker = mpu.Key
		result.NextUploadIdMarker = mpu.UploadId
	}
	return result
}
