//go:build unit
// +build unit

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

func TestListObjectsV2HappyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "ObjectV2ListingBucket"
	targetDelimiter := "OLDelimiter"
	targetEncodingType := "OLEncodingType"
	targetPrefix := "OLPrefix"

	var inputCaptureV2 *s3.ListObjectsV2Input

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectsV2Pages",
			mock.MatchedBy(
				func(inputV2 *s3.ListObjectsV2Input) bool {
					inputCaptureV2 = inputV2
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCaptureV2 := args.Get(1).(func(page *s3.ListObjectsV2Output, last bool) bool)
			pagerSrvMockOLV2(inputCaptureV2, funcCaptureV2, 0, nil)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ListObjectsV2,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--region", "REG",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectsV2Pages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCaptureV2.Bucket)
	assert.Equal(t, targetDelimiter, *inputCaptureV2.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCaptureV2.EncodingType)
	assert.Equal(t, targetPrefix, *inputCaptureV2.Prefix)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found no objects in bucket '%s'.", targetBucket))

}

func TestObjectsListV2WhenPageBiggerThanMaxRequestMax(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "ObjectV2ListingBucketT2"
	targetMaxKeys := int64(3)

	var inputCaptureV2 *s3.ListObjectsV2Input

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectsV2Pages",
			mock.MatchedBy(
				func(inputV2 *s3.ListObjectsV2Input) bool {
					inputCaptureV2 = inputV2
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ListObjectsV2,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, strconv.FormatInt(targetMaxKeys, 10),
		"--" + flags.PageSize, "100",
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectsV2Pages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCaptureV2.Bucket)
	assert.Equal(t, targetMaxKeys, *inputCaptureV2.MaxKeys)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")

}

func TestWhenPageSmallerThanMaxRequestPageV2(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "ObjectV2ListingBucketT3"
	targetMaxKeys := int64(100)

	var inputCaptureV2 *s3.ListObjectsV2Input

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectsV2Pages",
			mock.MatchedBy(
				func(inputV2 *s3.ListObjectsV2Input) bool {
					inputCaptureV2 = inputV2
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ListObjectsV2,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, "1000",
		"--" + flags.PageSize, strconv.FormatInt(targetMaxKeys, 10),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectsV2Pages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCaptureV2.Bucket)
	assert.Equal(t, targetMaxKeys, *inputCaptureV2.MaxKeys)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")

}

func TestObjectsListV2Paginate(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "ObjectV2ListingBucketT4"
	targetPageSz := 20
	targetMaxItems := 75

	var inputCaptureV2 *s3.ListObjectsV2Input

	pagesSzCapture := make([]int, 0, 0)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectsV2Pages",
			mock.MatchedBy(
				func(inputV2 *s3.ListObjectsV2Input) bool {
					inputCaptureV2 = inputV2
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCaptureV2 := args.Get(1).(func(page *s3.ListObjectsV2Output, last bool) bool)
			pagerSrvMockOLV2(inputCaptureV2, funcCaptureV2, 1000, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ListObjectsV2,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, strconv.Itoa(targetMaxItems),
		"--" + flags.PageSize, strconv.Itoa(targetPageSz),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectsV2Pages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCaptureV2.Bucket)

	// requested 75 elements split in 20 size pages
	assert.Equal(t, []int{20, 20, 20, 15}, pagesSzCapture)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found %d objects in bucket '%s'", targetMaxItems,
		targetBucket))

}

func pagerSrvMockOLV2(inputV2 *s3.ListObjectsV2Input, backFeed func(*s3.ListObjectsV2Output, bool) bool,
	hardLimit int, results *[]int) {
	keepGoing := true
	if hardLimit == 0 {
		parts := new(s3.ListObjectsV2Output)
		parts.Name = inputV2.Bucket
		backFeed(parts, hardLimit <= 0)
		return
	}
	for hardLimit > 0 && keepGoing {
		i := int(aws.Int64Value(inputV2.MaxKeys))
		*results = append(*results, i)
		if i == 0 || i > hardLimit {
			i = hardLimit
		}
		parts := buildListObjectsV2Output(i)
		parts.Name = inputV2.Bucket
		hardLimit -= i
		keepGoing = backFeed(parts, hardLimit <= 0)
	}
}

func buildListObjectsV2Output(howMany int) *s3.ListObjectsV2Output {
	result := new(s3.ListObjectsV2Output).SetContents([]*s3.Object{})
	for i := 0; i < howMany; i++ {
		obj := new(s3.Object).
			SetKey("key" + strconv.Itoa(i)).
			SetLastModified(time.Now()).
			SetSize(int64(2))
		result.Contents = append(result.Contents, obj)
		// result.NextMarker = obj.Key
	}
	return result
}
