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
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/IBM/ibmcloud-cos-cli/utils"
)

var (
	corsRulesJSONStr = `{
  "CORSRules":[
    {
      "AllowedHeaders":["AllowedHeader1","AllowedHeader2"],
      "AllowedMethods":["AllowedMethod3","AllowedMethod4"],
      "AllowedOrigins":["AllowedOrigin5","AllowedOrigin6"],
      "ExposeHeaders":["ExposeHeader7","ExposeHeader8"],
      "MaxAgeSeconds":75
    }
  ]
}`

	corsRulesSimpleJSONStr = `CORSRules=[
    {
      AllowedHeaders=[AllowedHeader1,AllowedHeader2],
      AllowedMethods=[AllowedMethod3,AllowedMethod4],
      AllowedOrigins=[AllowedOrigin5,AllowedOrigin6],
      ExposeHeaders=[ExposeHeader7,ExposeHeader8],
      MaxAgeSeconds=75
    }
  ]`

	corsRulesObject = new(s3.CORSConfiguration).
			SetCORSRules([]*s3.CORSRule{
			new(s3.CORSRule).
				SetAllowedHeaders([]*string{
					aws.String("AllowedHeader1"),
					aws.String("AllowedHeader2"),
				}).
				SetAllowedMethods([]*string{
					aws.String("AllowedMethod3"),
					aws.String("AllowedMethod4"),
				}).
				SetAllowedOrigins([]*string{
					aws.String("AllowedOrigin5"),
					aws.String("AllowedOrigin6"),
				}).
				SetExposeHeaders([]*string{
					aws.String("ExposeHeader7"),
					aws.String("ExposeHeader8"),
				}).
				SetMaxAgeSeconds(75),
		})
)

func TestBucketCorsPutSunnyPathJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedCorConfig *s3.CORSConfiguration

	providers.MockS3API.
		On("PutBucketCors", mock.MatchedBy(
			func(input *s3.PutBucketCorsInput) bool {
				capturedCorConfig = input.CORSConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketCorsOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketCorsPut, "--bucket", targetBucket, "--region", "REG", "--cors-configuration",
		corsRulesJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketCors", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	// assert json proper parsed
	assert.Equal(t, corsRulesObject, capturedCorConfig)

	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")

}

func TestBucketCorsPutSunnyPathJSONFile(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	targetFileName := "fileMock"

	isClosed := false

	var capturedCorConfig *s3.CORSConfiguration

	providers.MockS3API.
		On("PutBucketCors", mock.MatchedBy(
			func(input *s3.PutBucketCorsInput) bool {
				capturedCorConfig = input.CORSConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketCorsOutput), nil).
		Once()

	providers.MockFileOperations.
		On("ReadSeekerCloserOpen", mock.MatchedBy(func(fileName string) bool {
			return fileName == targetFileName
		})).
		Return(utils.WrapString(corsRulesJSONStr, &isClosed), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketCorsPut, "--bucket", targetBucket, "--region", "REG", "--cors-configuration",
		"file://" + targetFileName}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketCors", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	// assert json proper parsed
	assert.Equal(t, corsRulesObject, capturedCorConfig)

	// assert file is closed
	assert.True(t, isClosed)

	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")

}

func TestBucketCorsPutSunnyPathSimplifiedJSONString(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TARGETBUCKET"

	providers.MockPluginConfig.On("GetString", config.ServiceEndpointURL).Return("", nil)

	var capturedCorConfig *s3.CORSConfiguration

	providers.MockS3API.
		On("PutBucketCors", mock.MatchedBy(
			func(input *s3.PutBucketCorsInput) bool {
				capturedCorConfig = input.CORSConfiguration
				return *input.Bucket == targetBucket
			})).
		Return(new(s3.PutBucketCorsOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketCorsPut, "--bucket", targetBucket, "--region", "REG", "--cors-configuration",
		corsRulesSimpleJSONStr}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketCors", 1)
	// assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	// assert json proper parsed
	assert.Equal(t, corsRulesObject, capturedCorConfig)

	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert OK
	assert.Contains(t, output, "OK")
	// assert Not Fail
	assert.NotContains(t, errors, "FAIL")

}

func TestBucketCorsPutEmptyStaticCreds(t *testing.T) {
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
		On("PutBucketCors", mock.MatchedBy(
			func(input *s3.PutBucketCorsInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(nil, errors.New("EmptyStaticCreds")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketCorsPut, "--bucket", targetBucket, "--region", "REG", "--cors-configuration", "{}"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketCors", 1)
	// assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")

}
func TestBucketCorsPutWithoutCORSConfig(t *testing.T) {
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
		On("PutBucketCors", mock.MatchedBy(
			func(input *s3.PutBucketCorsInput) bool {
				return *input.Bucket == targetBucket
			})).
		Return(nil, errors.New("NoCORSConfig")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.BucketCorsPut, "--bucket", targetBucket,
		"--region", "REG"}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region (since success is last)
	providers.MockS3API.AssertNumberOfCalls(t, "PutBucketCors", 0)
	// assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli
	// capture all output //
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()
	// assert Not OK
	assert.NotContains(t, output, "OK")
	// assert Fail
	assert.Contains(t, errors, "FAIL")

}
