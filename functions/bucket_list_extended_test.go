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

func TestListBucketsExtendedHappyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetPrefix := "OLPrefix"
	targetMarker := "OLMarker"

	var inputCapture *s3.ListBucketsExtendedInput

	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us", nil)

	providers.MockS3API.
		On("ListBucketsExtendedPages",
			mock.MatchedBy(
				func(input *s3.ListBucketsExtendedInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ListBucketsExtended,
		"--" + flags.Prefix, targetPrefix,
		"--" + flags.Marker, targetMarker,
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetPrefix, *inputCapture.Prefix)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found no bucket"))
}

func TestBucketsListExtendedWhenPageBiggerThanMaxRequestMax(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetMaxKeys := int64(3)

	var inputCapture *s3.ListBucketsExtendedInput

	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us", nil)

	providers.MockS3API.
		On("ListBucketsExtendedPages",
			mock.MatchedBy(
				func(input *s3.ListBucketsExtendedInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ListBucketsExtended,
		"--" + flags.MaxItems, strconv.FormatInt(targetMaxKeys, 10),
		"--" + flags.PageSize, "100",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetMaxKeys, *inputCapture.MaxKeys)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")
}

func TestLBEWhenPageSmallerThanMaxRequestPage(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetMaxKeys := int64(100)

	var inputCapture *s3.ListBucketsExtendedInput

	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us", nil)

	providers.MockS3API.
		On("ListBucketsExtendedPages",
			mock.MatchedBy(
				func(input *s3.ListBucketsExtendedInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ListBucketsExtended,
		"--" + flags.MaxItems, "1000",
		"--" + flags.PageSize, strconv.FormatInt(targetMaxKeys, 10),
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetMaxKeys, *inputCapture.MaxKeys)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")
}

func TestBucketsListExtendedPaginate(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetPageSz := 20
	targetMaxItems := 75

	var inputCapture *s3.ListBucketsExtendedInput

	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us", nil)

	pagesSzCapture := make([]int, 0, 0)

	providers.MockS3API.
		On("ListBucketsExtendedPages",
			mock.MatchedBy(
				func(input *s3.ListBucketsExtendedInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListBucketsExtendedOutput, last bool) bool)
			pagerBktSrvMock(inputCapture, funcCapture, 1000, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ListBucketsExtended,
		"--" + flags.MaxItems, strconv.Itoa(targetMaxItems),
		"--" + flags.PageSize, strconv.Itoa(targetPageSz),
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	// requested 75 elements split in 20 size pages
	assert.Equal(t, []int{20, 20, 20, 15}, pagesSzCapture)

	output := providers.FakeUI.Outputs()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, output, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found %d buckets", targetMaxItems))

}

func pagerBktSrvMock(input *s3.ListBucketsExtendedInput, backFeed func(*s3.ListBucketsExtendedOutput, bool) bool,
	hardLimit int, results *[]int) {
	keepGoing := true
	for hardLimit > 0 && keepGoing {
		i := int(aws.Int64Value(input.MaxKeys))
		*results = append(*results, i)
		if i == 0 || i > hardLimit {
			i = hardLimit
		}
		parts := buildListBucketsExtendedOutput(i)
		hardLimit -= i
		keepGoing = backFeed(parts, hardLimit <= 0)
	}
}

func buildListBucketsExtendedOutput(howMany int) *s3.ListBucketsExtendedOutput {
	result := new(s3.ListBucketsExtendedOutput).SetBuckets([]*s3.BucketExtended{})
	for i := 0; i < howMany; i++ {
		bucket := new(s3.BucketExtended).
			SetName("bucket" + strconv.Itoa(i)).
			SetCreationDate(time.Now()).
			SetLocationConstraint("us-east")
		result.Buckets = append(result.Buckets, bucket)
	}
	return result
}
