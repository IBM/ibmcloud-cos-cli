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

func TestListObjectsHappyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "OLBucket"
	targetDelimiter := "OLDelimiter"
	targetEncodingType := "OLEncodingType"
	targetPrefix := "OLPrefix"
	targetMarker := "OLMarker"

	var inputCapture *s3.ListObjectsInput

	providers.MockS3API.
		On("ListObjectsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectsOutput, last bool) bool)
			pagerSrvMock(inputCapture, funcCapture, 0, nil)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Objects,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--" + flags.Marker, targetMarker,
		"--region", "REG",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetDelimiter, *inputCapture.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCapture.EncodingType)
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found no objects in bucket '%s'.", targetBucket))

}

func TestObjectsListWhenPageBiggerThanMaxRequestMax(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "OLBucketT2"
	targetMaxKeys := int64(3)

	var inputCapture *s3.ListObjectsInput

	providers.MockS3API.
		On("ListObjectsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Objects,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, strconv.FormatInt(targetMaxKeys, 10),
		"--" + flags.PageSize, "100",
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetMaxKeys, *inputCapture.MaxKeys)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")

}

func TestWhenPageSmallerThanMaxRequestPage(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "OLBucketT3"
	targetMaxKeys := int64(100)

	var inputCapture *s3.ListObjectsInput

	providers.MockS3API.
		On("ListObjectsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Objects,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, "1000",
		"--" + flags.PageSize, strconv.FormatInt(targetMaxKeys, 10),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetMaxKeys, *inputCapture.MaxKeys)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")

}

func TestObjectsListPaginate(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "OLBucketT4"
	targetPageSz := 20
	targetMaxItems := 75

	var inputCapture *s3.ListObjectsInput

	pagesSzCapture := make([]int, 0, 0)

	providers.MockS3API.
		On("ListObjectsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectsOutput, last bool) bool)
			pagerSrvMock(inputCapture, funcCapture, 1000, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.Objects,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, strconv.Itoa(targetMaxItems),
		"--" + flags.PageSize, strconv.Itoa(targetPageSz),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectsPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)

	// requested 75 elements split in 20 size pages
	assert.Equal(t, []int{20, 20, 20, 15}, pagesSzCapture)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found %d objects in bucket '%s'", targetMaxItems,
		targetBucket))

}

func pagerSrvMock(input *s3.ListObjectsInput, backFeed func(*s3.ListObjectsOutput, bool) bool,
	hardLimit int, results *[]int) {
	keepGoing := true
	if hardLimit == 0 {
		parts := new(s3.ListObjectsOutput)
		parts.Name = input.Bucket
		backFeed(parts, hardLimit <= 0)
		return
	}
	for hardLimit > 0 && keepGoing {
		i := int(aws.Int64Value(input.MaxKeys))
		*results = append(*results, i)
		if i == 0 || i > hardLimit {
			i = hardLimit
		}
		parts := buildListObjectsOutput(i)
		parts.Name = input.Bucket
		hardLimit -= i
		keepGoing = backFeed(parts, hardLimit <= 0)
	}
}

func buildListObjectsOutput(howMany int) *s3.ListObjectsOutput {
	result := new(s3.ListObjectsOutput).SetContents([]*s3.Object{})
	for i := 0; i < howMany; i++ {
		obj := new(s3.Object).
			SetKey("key" + strconv.Itoa(i)).
			SetLastModified(time.Now()).
			SetSize(int64(2))
		result.Contents = append(result.Contents, obj)
		result.NextMarker = obj.Key
	}
	return result
}
