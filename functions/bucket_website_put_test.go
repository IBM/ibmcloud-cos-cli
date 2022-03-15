//+build unit

package functions_test

import (
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

var (
	websiteConfigurationJSONStr = `{
  "ErrorDocument": {
    "Key": "error.html"
  },
  "IndexDocument": {
    "Suffix": "index.html"
  },
  "RoutingRules": [
    {
      "Condition": {
        "HttpErrorCodeReturnedEquals": "400",
        "KeyPrefixEquals": "valid_prefix"
      },
      "Redirect": {
        "HostName" : "routing.hostname.ibm.com",
        "HttpRedirectCode": "301",
        "Protocol": "https",
        "ReplaceKeyPrefixWith": "redirect_prefix",
        "ReplaceKeyWith": "replace_key"
      }
    }
  ]
}`

	websiteConfigurationSimpleJSONStr = `
  ErrorDocument={Key=error.html},
  IndexDocument={Suffix=index.html},
  RoutingRules=[{
    Condition={
      HttpErrorCodeReturnedEquals="400",
      KeyPrefixEquals=valid_prefix
    },
    Redirect={
      HostName=routing.hostname.ibm.com,
      HttpRedirectCode="301",
      Protocol=https,
      ReplaceKeyPrefixWith=redirect_prefix,
      ReplaceKeyWith=replace_key
    }
  }]
`

	websiteConfigurationObject = new(s3.WebsiteConfiguration).
					SetErrorDocument(new(s3.ErrorDocument).SetKey("error.html")).
					SetIndexDocument(new(s3.IndexDocument).SetSuffix("index.html")).
					SetRoutingRules([]*s3.RoutingRule{
			new(s3.RoutingRule).
				SetCondition(
					new(s3.Condition).
						SetHttpErrorCodeReturnedEquals("400").
						SetKeyPrefixEquals("valid_prefix"),
				).
				SetRedirect(
					new(s3.Redirect).
						SetHostName("routing.hostname.ibm.com").
						SetHttpRedirectCode("301").
						SetProtocol("https").
						SetReplaceKeyPrefixWith("redirect_prefix").
						SetReplaceKeyWith("replace_key"),
				),
		})

	websiteConfigurationRedirectJSONStr = `{
    "RedirectAllRequestsTo": {
      "HostName": "redirect.hostname.ibm.com",
      "Protocol": "https"
    }
}`

	websiteConfigurationRedirectSimpleJSONStr = `
  RedirectAllRequestsTo={HostName=redirect.hostname.ibm.com,Protocol=https}
`

	websiteConfigurationRedirectObject = new(s3.WebsiteConfiguration).
						SetRedirectAllRequestsTo(new(s3.RedirectAllRequestsTo).
							SetHostName("redirect.hostname.ibm.com").
							SetProtocol("https"))
)

func TestBucketWebsitePutSunnyPathJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedWebsiteConfiguration *s3.WebsiteConfiguration

	providers.MockS3API.
		On("PutBucketWebsite", mock.MatchedBy(
			func(input *s3.PutBucketWebsiteInput) bool {
				capturedWebsiteConfiguration = input.WebsiteConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketWebsiteOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsitePut, "--bucket", targetBucket, "--region", "REG",
		"--website-configuration", websiteConfigurationJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketWebsite", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, websiteConfigurationObject, capturedWebsiteConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketWebsitePutSunnyPathJSONFile(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	targetFileName := "fileMock"
	isClosed := false

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedWebsiteConfiguration *s3.WebsiteConfiguration

	providers.MockS3API.
		On("PutBucketWebsite", mock.MatchedBy(
			func(input *s3.PutBucketWebsiteInput) bool {
				capturedWebsiteConfiguration = input.WebsiteConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketWebsiteOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(websiteConfigurationJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsitePut, "--bucket", targetBucket, "--region", "REG",
		"--website-configuration", "file://" + targetFileName}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketWebsite", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, websiteConfigurationObject, capturedWebsiteConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketWebsitePutSunnyPathSimplifiedJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedWebsiteConfiguration *s3.WebsiteConfiguration

	providers.MockS3API.
		On("PutBucketWebsite", mock.MatchedBy(
			func(input *s3.PutBucketWebsiteInput) bool {
				capturedWebsiteConfiguration = input.WebsiteConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketWebsiteOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsitePut, "--bucket", targetBucket, "--region", "REG",
		"--website-configuration", websiteConfigurationSimpleJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketWebsite", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, websiteConfigurationObject, capturedWebsiteConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketWebsitePutSunnyPathRedirectRequestsJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedWebsiteConfiguration *s3.WebsiteConfiguration

	providers.MockS3API.
		On("PutBucketWebsite", mock.MatchedBy(
			func(input *s3.PutBucketWebsiteInput) bool {
				capturedWebsiteConfiguration = input.WebsiteConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketWebsiteOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsitePut, "--bucket", targetBucket, "--region", "REG",
		"--website-configuration", websiteConfigurationRedirectJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketWebsite", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, websiteConfigurationRedirectObject, capturedWebsiteConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketWebsitePutSunnyPathRedirectRequestsJSONFile(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"
	targetFileName := "fileMock"
	isClosed := false

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedWebsiteConfiguration *s3.WebsiteConfiguration

	providers.MockS3API.
		On("PutBucketWebsite", mock.MatchedBy(
			func(input *s3.PutBucketWebsiteInput) bool {
				capturedWebsiteConfiguration = input.WebsiteConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketWebsiteOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(websiteConfigurationRedirectJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsitePut, "--bucket", targetBucket, "--region", "REG",
		"--website-configuration", "file://" + targetFileName}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketWebsite", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, websiteConfigurationRedirectObject, capturedWebsiteConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketWebsitePutSunnyPathRedirectRequestsSimplifiedJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedWebsiteConfiguration *s3.WebsiteConfiguration

	providers.MockS3API.
		On("PutBucketWebsite", mock.MatchedBy(
			func(input *s3.PutBucketWebsiteInput) bool {
				capturedWebsiteConfiguration = input.WebsiteConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketWebsiteOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsitePut, "--bucket", targetBucket, "--region", "REG",
		"--website-configuration", websiteConfigurationRedirectSimpleJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketWebsite", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli
	// assert json proper parsed
	assert.Equal(t, websiteConfigurationRedirectObject, capturedWebsiteConfiguration)
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")
}

func TestBucketPutWebsiteWithoutBucket(t *testing.T) {
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
		On("PutBucketWebsite", mock.MatchedBy(
			func(input *s3.PutBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketWebsiteOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsitePut, "--region", "REG",
		"--website-configuration", websiteConfigurationJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketWebsite", 0)
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

func TestBucketPutWebsiteWithoutWebsiteConfiguration(t *testing.T) {
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
		On("PutBucketWebsite", mock.MatchedBy(
			func(input *s3.PutBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketWebsiteOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsitePut, "--bucket", targetBucket, "--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketWebsite", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "Mandatory Flag '--website-configuration' is missing")
}

func TestBucketPutWebsiteWithMalformedJsonWebsiteConfiguration(t *testing.T) {
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
		On("PutBucketWebsite", mock.MatchedBy(
			func(input *s3.PutBucketWebsiteInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketWebsiteOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketWebsitePut, "--bucket", targetBucket, "--region", "REG",
		"--website-configuration", "{\"ErrorDocument\": {\"Key\": \"error.html\"}},"} // trailing comma invalid
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketWebsite", 0)
	// assert exit code is non-zero
	assert.Equal(t, 1, *exitCode) // exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")
	assert.Contains(t, errors, "The value in flag '--website-configuration' is invalid")
}
