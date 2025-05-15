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

func TestListObjectVersionsNoObjectVersions(t *testing.T) {
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
	targetKeyMarker := "OLKeyMarker"
	targetVersionIdMarker := "OLVersionIdMarker"

	var inputCapture *s3.ListObjectVersionsInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectVersionsOutput, last bool) bool)
			pagerSrvMockOVL(inputCapture, funcCapture, 0, nil)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--key-marker", targetKeyMarker,
		"--version-id-marker", targetVersionIdMarker,
		"--region", "REG",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetDelimiter, *inputCapture.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCapture.EncodingType)
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)
	assert.Equal(t, targetKeyMarker, *inputCapture.KeyMarker)
	assert.Equal(t, targetVersionIdMarker, *inputCapture.VersionIdMarker)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found no object versions in bucket '%s'.", targetBucket))
	assert.Contains(t, output, fmt.Sprintf("Found no delete markers in bucket '%s'.", targetBucket))
}

func TestListObjectVersionsOneObjectVersion(t *testing.T) {
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
	targetKeyMarker := "OLKeyMarker"
	targetVersionIdMarker := "OLVersionIdMarker"

	var inputCapture *s3.ListObjectVersionsInput
	pagesSzCapture := make([]int, 0, 0)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectVersionsOutput, last bool) bool)
			pagerSrvMockOVL(inputCapture, funcCapture, 1, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--key-marker", targetKeyMarker,
		"--version-id-marker", targetVersionIdMarker,
		"--region", "REG",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetDelimiter, *inputCapture.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCapture.EncodingType)
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)
	assert.Equal(t, targetKeyMarker, *inputCapture.KeyMarker)
	assert.Equal(t, targetVersionIdMarker, *inputCapture.VersionIdMarker)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found 1 object version in bucket '%s':", targetBucket))
	assert.Contains(t, output, "Name   Version ID     Last Modified (UTC)        Object Size   Is Latest")
	assert.Contains(t, output, "key0   version-id-0")
	assert.Contains(t, output, "true")
	assert.Contains(t, output, fmt.Sprintf("Found no delete markers in bucket '%s'.", targetBucket))
}

func TestListObjectVersionsObjectVersionsTextNoTruncation(t *testing.T) {
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
	targetKeyMarker := "OLKeyMarker"
	targetVersionIdMarker := "OLVersionIdMarker"

	var inputCapture *s3.ListObjectVersionsInput
	pagesSzCapture := make([]int, 0, 0)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectVersionsOutput, last bool) bool)
			pagerSrvMockOVL(inputCapture, funcCapture, 5, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--key-marker", targetKeyMarker,
		"--version-id-marker", targetVersionIdMarker,
		"--region", "REG",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetDelimiter, *inputCapture.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCapture.EncodingType)
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)
	assert.Equal(t, targetKeyMarker, *inputCapture.KeyMarker)
	assert.Equal(t, targetVersionIdMarker, *inputCapture.VersionIdMarker)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found 5 object versions in bucket '%s':", targetBucket))
	assert.Contains(t, output, "Name   Version ID     Last Modified (UTC)        Object Size   Is Latest")
	assert.Contains(t, output, "key0   version-id-0")
	assert.Contains(t, output, "key0   version-id-1")
	assert.Contains(t, output, "key2   version-id-2")
	assert.Contains(t, output, "key2   version-id-3")
	assert.Contains(t, output, "key4   version-id-4")
	assert.Contains(t, output, "true")
	assert.Contains(t, output, "false")
	assert.Contains(t, output, fmt.Sprintf("Found no delete markers in bucket '%s'.", targetBucket))
}

func TestListObjectVersionsObjectVersionsJsonNoTruncation(t *testing.T) {
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
	targetKeyMarker := "OLKeyMarker"
	targetVersionIdMarker := "OLVersionIdMarker"

	var inputCapture *s3.ListObjectVersionsInput
	pagesSzCapture := make([]int, 0, 0)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectVersionsOutput, last bool) bool)
			pagerSrvMockOVL(inputCapture, funcCapture, 5, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--key-marker", targetKeyMarker,
		"--version-id-marker", targetVersionIdMarker,
		"--region", "REG",
		"--output", "json",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetDelimiter, *inputCapture.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCapture.EncodingType)
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)
	assert.Equal(t, targetKeyMarker, *inputCapture.KeyMarker)
	assert.Equal(t, targetVersionIdMarker, *inputCapture.VersionIdMarker)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.NotContains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, "\"DeleteMarkers\": null,")
}

func TestListObjectVersionsOneDeleteMarker(t *testing.T) {
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
	targetKeyMarker := "OLKeyMarker"
	targetVersionIdMarker := "OLVersionIdMarker"

	var inputCapture *s3.ListObjectVersionsInput
	pagesSzCapture := make([]int, 0, 0)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectVersionsOutput, last bool) bool)
			pagerSrvMockOVLDM(inputCapture, funcCapture, 1, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--key-marker", targetKeyMarker,
		"--version-id-marker", targetVersionIdMarker,
		"--region", "REG",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetDelimiter, *inputCapture.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCapture.EncodingType)
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)
	assert.Equal(t, targetKeyMarker, *inputCapture.KeyMarker)
	assert.Equal(t, targetVersionIdMarker, *inputCapture.VersionIdMarker)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found no object versions in bucket '%s'.", targetBucket))
	assert.Contains(t, output, fmt.Sprintf("Found 1 delete marker in bucket '%s':", targetBucket))
	assert.Contains(t, output, "Name   Version ID     Last Modified (UTC)        Is Latest")
	assert.Contains(t, output, "key0   version-id-0")
	assert.Contains(t, output, "true")
}

func TestListObjectVersionsDeleteMarkersTextNoTruncation(t *testing.T) {
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
	targetKeyMarker := "OLKeyMarker"
	targetVersionIdMarker := "OLVersionIdMarker"

	var inputCapture *s3.ListObjectVersionsInput
	pagesSzCapture := make([]int, 0, 0)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectVersionsOutput, last bool) bool)
			pagerSrvMockOVLDM(inputCapture, funcCapture, 5, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--key-marker", targetKeyMarker,
		"--version-id-marker", targetVersionIdMarker,
		"--region", "REG",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetDelimiter, *inputCapture.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCapture.EncodingType)
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)
	assert.Equal(t, targetKeyMarker, *inputCapture.KeyMarker)
	assert.Equal(t, targetVersionIdMarker, *inputCapture.VersionIdMarker)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found no object versions in bucket '%s'.", targetBucket))
	assert.Contains(t, output, fmt.Sprintf("Found 5 delete markers in bucket '%s':", targetBucket))
	assert.Contains(t, output, "Name   Version ID     Last Modified (UTC)        Is Latest")
	assert.Contains(t, output, "key0   version-id-0")
	assert.Contains(t, output, "key0   version-id-1")
	assert.Contains(t, output, "key2   version-id-2")
	assert.Contains(t, output, "key2   version-id-3")
	assert.Contains(t, output, "key4   version-id-4")
	assert.Contains(t, output, "true")
	assert.Contains(t, output, "false")
}

func TestListObjectVersionsDeleteMarkersJsonNoTruncation(t *testing.T) {
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
	targetKeyMarker := "OLKeyMarker"
	targetVersionIdMarker := "OLVersionIdMarker"

	var inputCapture *s3.ListObjectVersionsInput
	pagesSzCapture := make([]int, 0, 0)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectVersionsOutput, last bool) bool)
			pagerSrvMockOVLDM(inputCapture, funcCapture, 5, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Delimiter, targetDelimiter,
		"--" + flags.EncodingType, targetEncodingType,
		"--" + flags.Prefix, targetPrefix,
		"--key-marker", targetKeyMarker,
		"--version-id-marker", targetVersionIdMarker,
		"--region", "REG",
		"--output", "json",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetDelimiter, *inputCapture.Delimiter)
	assert.Equal(t, targetEncodingType, *inputCapture.EncodingType)
	assert.Equal(t, targetPrefix, *inputCapture.Prefix)
	assert.Equal(t, targetKeyMarker, *inputCapture.KeyMarker)
	assert.Equal(t, targetVersionIdMarker, *inputCapture.VersionIdMarker)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.NotContains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, "\"Versions\": null")
}

func TestListObjectVersionsWhenPageBiggerThanMaxRequestMax(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "OLBucketT2"
	targetMaxKeys := int64(3)

	var inputCapture *s3.ListObjectVersionsInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, strconv.FormatInt(targetMaxKeys, 10),
		"--" + flags.PageSize, "100",
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetMaxKeys, *inputCapture.MaxKeys)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
}

func TestListObjectVersionsWhenPageSmallerThanMaxRequestPage(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "OLBucketT3"
	targetMaxKeys := int64(100)

	var inputCapture *s3.ListObjectVersionsInput

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, "1000",
		"--" + flags.PageSize, strconv.FormatInt(targetMaxKeys, 10),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)
	assert.Equal(t, targetMaxKeys, *inputCapture.MaxKeys)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
}

func TestListObjectVersionsPaginate(t *testing.T) {
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

	var inputCapture *s3.ListObjectVersionsInput

	pagesSzCapture := make([]int, 0, 0)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectVersionsOutput, last bool) bool)
			pagerSrvMockOVL(inputCapture, funcCapture, 1000, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, strconv.Itoa(targetMaxItems),
		"--" + flags.PageSize, strconv.Itoa(targetPageSz),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)

	// requested 75 elements split in 20 size pages
	assert.Equal(t, []int{20, 20, 20, 15}, pagesSzCapture)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found %d object versions in bucket '%s'", targetMaxItems, targetBucket))
	assert.Contains(t, output, fmt.Sprintf("Found no delete markers in bucket '%s'.", targetBucket))
}

func TestListObjectVersionsDeleteMarkersPaginate(t *testing.T) {
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

	var inputCapture *s3.ListObjectVersionsInput

	pagesSzCapture := make([]int, 0, 0)

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("ListObjectVersionsPages",
			mock.MatchedBy(
				func(input *s3.ListObjectVersionsInput) bool {
					inputCapture = input
					return true
				}),
			mock.Anything).
		Run(func(args mock.Arguments) {
			funcCapture := args.Get(1).(func(page *s3.ListObjectVersionsOutput, last bool) bool)
			pagerSrvMockOVLDM(inputCapture, funcCapture, 1000, &pagesSzCapture)
		}).
		Return(nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-",
		commands.ObjectVersions,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.MaxItems, strconv.Itoa(targetMaxItems),
		"--" + flags.PageSize, strconv.Itoa(targetPageSz),
		"--region", "REG",
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "ListObjectVersionsPages", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	assert.Equal(t, targetBucket, *inputCapture.Bucket)

	// requested 75 elements split in 20 size pages
	assert.Equal(t, []int{20, 20, 20, 15}, pagesSzCapture)

	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	assert.Contains(t, output, "OK")
	assert.NotContains(t, errors, "FAIL")
	assert.Contains(t, output, fmt.Sprintf("Found no object versions in bucket '%s'.", targetBucket))
	assert.Contains(t, output, fmt.Sprintf("Found %d delete markers in bucket '%s'", targetMaxItems, targetBucket))
}

func pagerSrvMockOVL(input *s3.ListObjectVersionsInput, backFeed func(*s3.ListObjectVersionsOutput, bool) bool,
	hardLimit int, results *[]int) {
	keepGoing := true
	if hardLimit == 0 {
		parts := new(s3.ListObjectVersionsOutput)
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
		parts := buildListObjectVersionsOutput(i)
		parts.Name = input.Bucket
		hardLimit -= i
		keepGoing = backFeed(parts, hardLimit <= 0)
	}
}

func buildListObjectVersionsOutput(howMany int) *s3.ListObjectVersionsOutput {
	result := new(s3.ListObjectVersionsOutput).SetVersions([]*s3.ObjectVersion{})
	var suffix string
	var isLatest bool
	for i := 0; i < howMany; i++ {
		if i%2 == 0 {
			suffix = strconv.Itoa(i)
			isLatest = (i + 1) == howMany // If this is the end of the loop, this is the latest (only) version
		} else {
			suffix = strconv.Itoa(i - 1)
			isLatest = true
		}
		key := "key" + suffix
		versionId := "version-id-" + strconv.Itoa(i)
		obj := new(s3.ObjectVersion).
			SetKey(key).
			SetVersionId(versionId).
			SetIsLatest(isLatest).
			SetLastModified(time.Now()).
			SetSize(int64(2))
		result.Versions = append(result.Versions, obj)
		result.NextKeyMarker = obj.Key
		result.NextVersionIdMarker = obj.VersionId
	}
	return result
}

func pagerSrvMockOVLDM(input *s3.ListObjectVersionsInput, backFeed func(*s3.ListObjectVersionsOutput, bool) bool,
	hardLimit int, results *[]int) {
	keepGoing := true
	if hardLimit == 0 {
		parts := new(s3.ListObjectVersionsOutput)
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
		parts := buildListObjectVersionsOutputDM(i)
		parts.Name = input.Bucket
		hardLimit -= i
		keepGoing = backFeed(parts, hardLimit <= 0)
	}
}

func buildListObjectVersionsOutputDM(howMany int) *s3.ListObjectVersionsOutput {
	result := new(s3.ListObjectVersionsOutput).SetDeleteMarkers([]*s3.DeleteMarkerEntry{})
	var suffix string
	var isLatest bool
	for i := 0; i < howMany; i++ {
		if i%2 == 0 {
			suffix = strconv.Itoa(i)
			isLatest = (i + 1) == howMany // If this is the end of the loop, this is the latest (only) version
		} else {
			suffix = strconv.Itoa(i - 1)
			isLatest = true
		}
		key := "key" + suffix
		versionId := "version-id-" + strconv.Itoa(i)
		obj := new(s3.DeleteMarkerEntry).
			SetKey(key).
			SetVersionId(versionId).
			SetIsLatest(isLatest).
			SetLastModified(time.Now())
		result.DeleteMarkers = append(result.DeleteMarkers, obj)
		result.NextKeyMarker = obj.Key
		result.NextVersionIdMarker = obj.VersionId
	}
	return result
}
