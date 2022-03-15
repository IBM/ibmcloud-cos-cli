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
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
)

func TestGetBucketClassSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}
	// set Fake regions
	fakeRegions := []string{"r1", "r2", "r3"}
	// set RegionResolverMock to return fake regions
	providers.MockRegionResolver.ListKnownRegions.On("GetAllRegions").Return(fakeRegions, nil)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	// for the first N-1 calls to get location return error
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(nil, aws.ErrMissingRegion).
		Times(len(fakeRegions) - 1)
	// for last call of get location return a valid result
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(new(s3.GetBucketLocationOutput).SetLocationConstraint("r1-vault"), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", "bucket-class-get", "--bucket", "bucket"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketLocationWithContext", len(fakeRegions))
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
	// assert class: Vault
	assert.Contains(t, output, "Vault")

}

func TestGetBucketClassRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}
	// set Fake regions
	fakeRegions := []string{"r1", "r2", "r3"}
	// set RegionResolverMock to return fake regions
	providers.MockRegionResolver.ListKnownRegions.On("GetAllRegions").Return(fakeRegions, nil)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	// for the first N-1 calls to get location return error
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(nil, aws.ErrMissingRegion).
		Times(len(fakeRegions) - 1)
	// for last call of get location return a valid result
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(nil, errors.New("BadLocation")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", "bucket-class-get", "--bucket", "bucket"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketLocationWithContext", len(fakeRegions))
	// assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	// assert Not class: Vault
	assert.NotContains(t, output, "Vault")

}

func TestGetBucketLocationSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}
	// set Fake regions
	fakeRegions := []string{"r1", "r2", "r3"}
	// set RegionResolverMock to return fake regions
	providers.MockRegionResolver.ListKnownRegions.On("GetAllRegions").Return(fakeRegions, nil)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	// for the first N-1 calls to get location return error
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(nil, aws.ErrMissingRegion).
		Times(len(fakeRegions) - 1)
	// for last call of get location return a valid result
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(new(s3.GetBucketLocationOutput).SetLocationConstraint("TARGETREGION-cold"), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", "bucket-location-get", "--bucket", "bucket"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketLocationWithContext", len(fakeRegions))
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
	// assert class: Vault
	assert.Contains(t, output, "Cold Vault")
	// assert region: r3
	assert.Contains(t, output, "TARGETREGION")

}

func TestGetBucketLocationRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}
	// set Fake regions
	fakeRegions := []string{"r1", "r2", "r3"}
	// set RegionResolverMock to return fake regions
	providers.MockRegionResolver.ListKnownRegions.On("GetAllRegions").Return(fakeRegions, nil)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	// for the first N-1 calls to get location return error
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(nil, aws.ErrMissingRegion).
		Times(len(fakeRegions) - 1)
	// for last call of get location return a valid result
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(nil, errors.New("BadLocation")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", "bucket-location-get", "--bucket", "bucket"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketLocationWithContext", len(fakeRegions))
	// assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	// assert Not class: Vault
	assert.NotContains(t, output, "Cold Vault")
	// assert Not region: r3
	assert.NotContains(t, output, "TARGETREGION")

}

func TestGetBucketLocationWithoutBucket(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}
	// set Fake regions
	fakeRegions := []string{"r1", "r2", "r3"}
	// set RegionResolverMock to return fake regions
	providers.MockRegionResolver.ListKnownRegions.On("GetAllRegions").Return(fakeRegions, nil)
	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)
	// for the first N-1 calls to get location return error
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(nil, aws.ErrMissingRegion).
		Times(len(fakeRegions) - 1)
	// for last call of get location return a valid result
	providers.MockS3API.
		On("GetBucketLocationWithContext", mock.Anything, mock.Anything).
		Return(nil, errors.New("BadLocation")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", "bucket-location-get"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketLocationWithContext", 0)
	// assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	// assert Not class: Vault
	assert.NotContains(t, output, "Cold Vault")
	// assert Not region: r3
	assert.NotContains(t, output, "TARGETREGION")

}
