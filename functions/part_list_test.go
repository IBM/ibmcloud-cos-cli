//+build unit

package functions_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
)

func TestListPartsHappyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "PLBucket"
	targetKey := "PLKey"
	targetUploadId := "PLUploadId"
	targetPartNumberMarker := 1

	var inputCapture *s3.ListPartsInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListPartsPages",
			mock.MatchedBy(
				func(input *s3.ListPartsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListPartsOutput, last bool) bool)
			pagerSrvMockPL(inputCapture, funcCapture, 0, nil)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Parts,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.UploadID, targetUploadId,
		"--" + flags.PartNumberMarker, strconv.Itoa(targetPartNumberMarker),
		"--region", "REG",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListPartsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetKey, *inputCapture.Key)
	assert.Equal(t, targetUploadId, *inputCapture.UploadId)
	assert.Equal(t, int64(targetPartNumberMarker), *inputCapture.PartNumberMarker)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found no parts in the multipart upload for '%s'.", targetKey))

}

func TestPartsListWhenPageBiggerThanMaxRequestMax(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "PLBucketT2"
	targetKey := "PLKeyT2"
	targetUploadId := "PLUploadIDT2"
	targetMaxParts := int64(3)

	var inputCapture *s3.ListPartsInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListPartsPages",
			mock.MatchedBy(
				func(input *s3.ListPartsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Parts,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.UploadID, targetUploadId,
		"--" + flags.MaxItems, strconv.FormatInt(targetMaxParts, 10),
		"--" + flags.PageSize, "100",
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListPartsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetMaxParts, *inputCapture.MaxParts)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")

}

func TestListPartsWhenPageSmallerThanMaxRequestPage(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "PLBucketT3"
	targetKey := "PLKeyT3"
	targetUploadId := "PLUploadIDT3"
	targetMaxParts := int64(100)

	var inputCapture *s3.ListPartsInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListPartsPages",
			mock.MatchedBy(
				func(input *s3.ListPartsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Parts,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.UploadID, targetUploadId,
		"--" + flags.MaxItems, "1000",
		"--" + flags.PageSize, strconv.FormatInt(targetMaxParts, 10),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListPartsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetMaxParts, *inputCapture.MaxParts)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")

}

func TestPartsListPaginate(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "PLBucketT4"
	targetKey := "PLKeyT4"
	targetUploadId := "PLUploadIDT4"
	targetPageSz := 20
	targetMaxItems := 75

	var inputCapture *s3.ListPartsInput

	pagesSzCapture := make([]int, 0, 0)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListPartsPages",
			mock.MatchedBy(
				func(input *s3.ListPartsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListPartsOutput, last bool) bool)
			pagerSrvMockPL(inputCapture, funcCapture, 1000, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Parts,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.UploadID, targetUploadId,
		"--" + flags.MaxItems, strconv.Itoa(targetMaxItems),
		"--" + flags.PageSize, strconv.Itoa(targetPageSz),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListPartsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)

	// requested 75 elements split in 20 size pages
	assert.Equal(t, []int{20, 20, 20, 15}, pagesSzCapture)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found %d parts in the multipart upload for '%s'", targetMaxItems,
		targetKey))

}

func pagerSrvMockPL(input *s3.ListPartsInput, backFeed func(*s3.ListPartsOutput, bool) bool,
	hardLimit int, results *[]int) {
	keepGoing := true
	if hardLimit == 0 {
		parts := new(s3.ListPartsOutput)
		parts.Bucket = input.Bucket
		parts.Key = input.Key
		backFeed(parts, hardLimit <= 0)
		return
	}
	for hardLimit > 0 && keepGoing {
		i := int(aws.Int64Value(input.MaxParts))
		*results = append(*results, i)
		if i == 0 || i > hardLimit {
			i = hardLimit
		}
		parts := buildListPartsOutput(i)
		parts.Bucket = input.Bucket
		parts.Key = input.Key
		hardLimit -= i
		keepGoing = backFeed(parts, hardLimit <= 0)
	}
}

func buildListPartsOutput(howMany int) *s3.ListPartsOutput {
	result := new(s3.ListPartsOutput).SetParts([]*s3.Part{})
	for i := 0; i < howMany; i++ {
		part := new(s3.Part).
			SetLastModified(time.Now()).
			SetETag("ETag" + strconv.Itoa(i)).
			SetPartNumber(int64(i)).
			SetSize(int64(2))
		result.Parts = append(result.Parts, part)
		result.NextPartNumberMarker = part.PartNumber
	}
	return result
}
