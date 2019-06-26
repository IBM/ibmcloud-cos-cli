//+build unit

package providers

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	terminalHelpers "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/endpoints"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
	"github.com/IBM/ibmcloud-cos-cli/di/providers/mocks"
	"github.com/IBM/ibmcloud-cos-cli/utils"
)

// Test Environment Shared Mocks
var (
	FakeUI            = new(terminalHelpers.FakeUI)
	MockS3API         = new(mocks.S3API)
	MockUploaderAPI   = new(mocks.Uploader)
	MockDownloaderAPI = new(mocks.Downloader)
	MockPluginConfig  = new(mocks.PluginConfig)

	MockRegionResolver = &RegionResolverMock{
		ListKnownRegions: mocks.ListKnownRegions{},
		Resolver:         mocks.Resolver{},
	}

	MockFileOperations = new(mocks.FileOperations)

	ReferenceUploader   = new(s3manager.Uploader)
	ReferenceDownloader = new(s3manager.Downloader)
)

func MocksRESET() {
	FakeUI = new(terminalHelpers.FakeUI)
	MockS3API = new(mocks.S3API)
	MockUploaderAPI = new(mocks.Uploader)
	MockDownloaderAPI = new(mocks.Downloader)
	MockPluginConfig = new(mocks.PluginConfig)

	MockRegionResolver = &RegionResolverMock{
		ListKnownRegions: mocks.ListKnownRegions{},
		Resolver:         mocks.Resolver{},
	}

	MockFileOperations = new(mocks.FileOperations)

	ReferenceUploader = new(s3manager.Uploader)
	ReferenceDownloader = new(s3manager.Downloader)
}

type RegionResolverMock struct {
	mocks.ListKnownRegions
	mocks.Resolver
}

func NewUI() terminal.UI {
	return FakeUI
}

func GetS3APIFn() func(*session.Session) s3iface.S3API {
	return func(*session.Session) s3iface.S3API {
		return MockS3API
	}
}

func GetDownloaderAPIFn() func(svc s3iface.S3API, options ...func(*s3manager.Downloader)) utils.Downloader {
	return func(svc s3iface.S3API, options ...func(*s3manager.Downloader)) utils.Downloader {
		for _, fun := range options {
			fun(ReferenceDownloader)
		}
		return MockDownloaderAPI
	}
}

func GetUploaderAPIFn() func(svc s3iface.S3API, options ...func(output *s3manager.Uploader)) utils.Uploader {
	return func(svc s3iface.S3API, options ...func(output *s3manager.Uploader)) utils.Uploader {
		for _, fun := range options {
			fun(ReferenceUploader)
		}
		return MockUploaderAPI
	}
}

// maybe mock the provider to assert calling parameters
func GetPluginConfig(_ plugin.PluginContext) plugin.PluginConfig {
	return MockPluginConfig
}

func NewSession(_ *aws.Config) (*session.Session, error) {
	return new(session.Session), nil
}

// mocks should intersect before it is effective needed
func NewConfig(_ plugin.PluginContext, _ endpoints.Resolver, _ *BaseConfig) (*aws.Config, error) {
	return nil, nil
}

func NewCOSEndPointsWSClient(_ plugin.PluginContext, _ *BaseConfig) (*RegionResolverMock, error) {
	return MockRegionResolver, nil
}

func GetFileOperations() utils.FileOperations {
	return MockFileOperations
}

type BaseConfig aws.Config

// mocks should intersect before it is effective needed
func GetBaseConfig(_ plugin.PluginContext) *BaseConfig {
	return nil
}
