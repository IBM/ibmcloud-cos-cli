//+build unit

package functions_test

import (
	"errors"
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

func TestBucketWebsiteGetSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("GetBucketWebsite", mock.MatchedBy(
			func(input *s3.GetBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketWebsiteOutput).
			SetErrorDocument(&s3.ErrorDocument{Key: aws.String("error.html")}).
			SetIndexDocument(&s3.IndexDocument{Suffix: aws.String("index.html")}).
			SetRedirectAllRequestsTo(new(s3.RedirectAllRequestsTo)).
			SetRoutingRules([]*s3.RoutingRule{}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsiteGet, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketWebsite", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	assert.Contains(t, output, "Index Suffix: index.html")
	assert.Contains(t, output, "Error Document: error.html")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketWebsiteGetJsonSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("GetBucketWebsite", mock.MatchedBy(
			func(input *s3.GetBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketWebsiteOutput).
			SetErrorDocument(&s3.ErrorDocument{Key: aws.String("error.html")}).
			SetIndexDocument(&s3.IndexDocument{Suffix: aws.String("index.html")}).
			SetRedirectAllRequestsTo(new(s3.RedirectAllRequestsTo)).
			SetRoutingRules([]*s3.RoutingRule{}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsiteGet, "--bucket", targetBucket, "--region", "REG", "--output", "json"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketWebsite", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "\"ErrorDocument\": ")
	assert.Contains(t, output, "\"Key\": \"error.html\"")
	assert.Contains(t, output, "\"IndexDocument\": ")
	assert.Contains(t, output, "\"Suffix\": \"index.html\"")
	assert.Contains(t, output, "\"RedirectAllRequestsTo\":")
	assert.Contains(t, output, "\"RoutingRules\": []")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketWebsiteGetWithoutBucket(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("GetBucketWebsite", mock.MatchedBy(
			func(input *s3.GetBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.GetBucketWebsiteOutput).
			SetErrorDocument(&s3.ErrorDocument{Key: aws.String("error.html")}).
			SetIndexDocument(&s3.IndexDocument{Suffix: aws.String("index.html")}).
			SetRedirectAllRequestsTo(new(s3.RedirectAllRequestsTo)).
			SetRoutingRules([]*s3.RoutingRule{}), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsiteGet, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketWebsite", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
}

func TestBucketWebsiteGetNoWebsite(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	providers.MockS3API.
		On("GetBucketWebsite", mock.MatchedBy(
			func(input *s3.GetBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(nil, errors.New("NoSuchWebsiteConfiguration")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsiteGet, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "GetBucketWebsite", 1)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
}
