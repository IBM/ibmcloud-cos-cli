//+build unit

package functions_test

import (
	"encoding/json"
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

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

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
		commands.BucketsExtended,
		"--" + flags.Prefix, targetPrefix,
		"--" + flags.Marker, targetMarker,
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
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

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

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
		commands.BucketsExtended,
		"--" + flags.MaxItems, strconv.FormatInt(targetMaxKeys, 10),
		"--" + flags.PageSize, "100",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	assert.Equal(t, targetMaxKeys, *inputCapture.MaxKeys)
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
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

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

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
		commands.BucketsExtended,
		"--" + flags.MaxItems, "1000",
		"--" + flags.PageSize, strconv.FormatInt(targetMaxKeys, 10),
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	assert.Equal(t, targetMaxKeys, *inputCapture.MaxKeys)
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
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

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

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
		commands.BucketsExtended,
		"--" + flags.MaxItems, strconv.Itoa(targetMaxItems),
		"--" + flags.PageSize, strconv.Itoa(targetPageSz),
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// requested 75 elements split in 20 size pages
	assert.Equal(t, []int{20, 20, 20, 15}, pagesSzCapture)
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
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

func TestListBucketsExtendedCreationTemplateIdText(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	expectedCreationDate := time.Now()

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us", nil)

	providers.MockS3API.
		On("ListBucketsExtendedPages",
			mock.MatchedBy(
				func(input *s3.ListBucketsExtendedInput) bool {
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListBucketsExtendedOutput, last bool) bool)
			mockCreationTemplateIdPager(expectedCreationDate, funcCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.BucketsExtended,
		"--output", "text",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "Location Constraint")
	assert.Contains(t, output, "Creation Date")
	assert.Contains(t, output, "Creation Template ID")
	assert.Contains(t, output, "BucketName")
	assert.Contains(t, output, "CreationTemplateId")
	assert.Contains(t, output, expectedCreationDate.Format("Jan 02, 2006 at 15:04:05"))
	assert.Contains(t, output, "LocationConstraint")
}

func TestListBucketsExtendedCreationTemplateIdJson(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	expectedCreationDate := time.Now()

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	providers.MockPluginConfig.
		On("GetStringWithDefault", "Default Region", mock.AnythingOfType("string")).
		Return("us", nil)

	providers.MockS3API.
		On("ListBucketsExtendedPages",
			mock.MatchedBy(
				func(input *s3.ListBucketsExtendedInput) bool {
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListBucketsExtendedOutput, last bool) bool)
			mockCreationTemplateIdPager(expectedCreationDate, funcCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.BucketsExtended,
		"--output", "json",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListBucketsExtendedPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.NotContains(t, errors, "FAIL")
	parsedOutput := new(s3.ListBucketsExtendedOutput)
	err := json.Unmarshal([]byte(output), &parsedOutput)
	assert.Nil(t, err)
	assert.NotNil(t, parsedOutput.Buckets)
	assert.Equal(t, len(parsedOutput.Buckets), 1)
	bucketEntry := parsedOutput.Buckets[0]
	assert.Contains(t, expectedCreationDate.String(), aws.TimeValue(bucketEntry.CreationDate).String())
	assert.Equal(t, aws.StringValue(bucketEntry.CreationTemplateId), "CreationTemplateId")
	assert.Equal(t, aws.StringValue(bucketEntry.LocationConstraint), "LocationConstraint")
	assert.Equal(t, aws.StringValue(bucketEntry.Name), "BucketName")
}

func mockCreationTemplateIdPager(expectedCreationDate time.Time, helper func(*s3.ListBucketsExtendedOutput, bool) bool) {
	mockOutput := mockCreationTemplateIdBuilder(expectedCreationDate)
	helper(mockOutput, true)
	return
}

func mockCreationTemplateIdBuilder(expectedCreationDate time.Time) *s3.ListBucketsExtendedOutput {
	mockEntry := new(s3.BucketExtended).
		SetCreationDate(expectedCreationDate).
		SetCreationTemplateId("CreationTemplateId").
		SetLocationConstraint("LocationConstraint").
		SetName("BucketName")
	mockList := []*s3.BucketExtended{}
	mockList = append(mockList, mockEntry)
	mockOutput := new(s3.ListBucketsExtendedOutput).
		SetBuckets(mockList)
	return mockOutput
}
