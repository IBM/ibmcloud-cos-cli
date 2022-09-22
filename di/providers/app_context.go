//+build !unit

package providers

import (
	gohttp "net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/http"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/trace"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials"
	"github.com/IBM/ibm-cos-sdk-go/aws/endpoints"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
	"github.com/IBM/ibmcloud-cos-cli/aspera"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/utils"
)

func NewUI() terminal.UI {
	return terminal.NewStdUI()
}

func GetS3API(sess *session.Session) s3iface.S3API {
	return s3.New(sess)
}

func GetS3APIFn() func(*session.Session) s3iface.S3API {
	return GetS3API
}

func GetDownloaderAPI(svc s3iface.S3API, options ...func(*s3manager.Downloader)) utils.Downloader {
	return s3manager.NewDownloaderWithClient(svc, options...)
}

func GetDownloaderAPIFn() func(svc s3iface.S3API, options ...func(*s3manager.Downloader)) utils.Downloader {
	return GetDownloaderAPI
}

func GetUploaderAPI(svc s3iface.S3API, options ...func(output *s3manager.Uploader)) utils.Uploader {
	return s3manager.NewUploaderWithClient(svc, options...)
}

func GetUploaderAPIFn() func(svc s3iface.S3API, options ...func(output *s3manager.Uploader)) utils.Uploader {
	return GetUploaderAPI
}

func GetAsperaTransfer(sess *session.Session, apikey string) (utils.AsperaTransfer, error) {
	svc := s3.New(sess)
	return aspera.New(svc, apikey)
}

func GetAsperaTransferFn() func(sess *session.Session, apikey string) (utils.AsperaTransfer, error) {
	return GetAsperaTransfer
}

func GetPluginConfig(ctx plugin.PluginContext) plugin.PluginConfig {
	return ctx.PluginConfig()
}

func NewSession(config *aws.Config) (sess *session.Session, err error) {
	sess, err = session.NewSession(config)
	if err == nil {
		sess.Handlers.Build.PushBackNamed(utils.CLIVersionUserAgentHandler)
	}
	return
}

func NewConfig(ctx plugin.PluginContext, resolver endpoints.Resolver, baseConfig *BaseConfig) (*aws.Config, error) {
	var conf *aws.Config
	if baseConfig == nil {
		conf = new(aws.Config)
	} else {
		conf = (*aws.Config)(baseConfig).Copy()
	}

	conf.DisableRestProtocolURICleaning = aws.Bool(true)

	if hmac, _ := ctx.PluginConfig().GetBoolWithDefault(config.HMACProvided, config.HMACProvidedDefault); hmac {
		id, _ := ctx.PluginConfig().GetString(config.AccessKeyID)
		//if err != nil || id == "" {
		//	return nil, errors.New("error.getting.HMAC.ID")
		//}
		secret, _ := ctx.PluginConfig().GetString(config.SecretAccessKey)
		//if err != nil || secret == "" {
		//	return nil, errors.New("error.getting.HMAC.SECRET")
		//}

		conf.Credentials = credentials.NewStaticCredentials(id, secret, "")
	} else {
		conf.Credentials = utils.NewBxBridgeCredentials(ctx)
	}

	forcePathStyle, err := ctx.PluginConfig().GetBoolWithDefault(config.ForcePathStyle, config.ForcePathStyleDefault)
	if err != nil {
		return nil, err
	}
	conf.WithS3ForcePathStyle(forcePathStyle)

	conf.WithEndpointResolver(resolver)

	return conf, nil
}

func NewCOSEndPointsWSClient(ctx plugin.PluginContext, conf *BaseConfig) (*utils.COSEndPointsWSClient, error) {
	// silently discard error as it will fallback to the default Production Endpoint
	regionsEndPoint, _ := ctx.PluginConfig().GetStringWithDefault(config.RegionsEndpointURL, "")
	// using default config
	return utils.NewIBMEndPoints((*aws.Config)(conf), regionsEndPoint)
}

type FileOperationsImpl struct{}

func (_ *FileOperationsImpl) ReadSeekerCloserOpen(location string) (utils.ReadSeekerCloser, error) {
	return os.Open(location)
}

func (_ *FileOperationsImpl) WriteCloserOpen(location string) (utils.WriteCloser, error) {
	return os.Create(location)
}

func (_ *FileOperationsImpl) GetFileInfo(location string) (os.FileInfo, error) {
	return os.Stat(location)
}

func (_ *FileOperationsImpl) Remove(location string) error {
	return os.Remove(location)
}

func (_ *FileOperationsImpl) GetTotalBytes(location string) (int64, error) {
	var size int64
	err := filepath.Walk(location, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func GetFileOperations() *FileOperationsImpl {
	return new(FileOperationsImpl)
}

type BaseConfig aws.Config

type loggerWrap struct{}

func (_ *loggerWrap) Log(in ...interface{}) {
	trace.Logger.Println(in)
}

func GetBaseConfig(ctx plugin.PluginContext) *BaseConfig {

	region := ctx.CurrentRegion()
	if region.Name != "" {
		config.FallbackRegion = region.Name
	}

	conf := new(aws.Config)
	client := &gohttp.Client{
		Timeout: time.Duration(ctx.HTTPTimeout()) * time.Second,
	}

	if ctx.Trace() == "true" {
		trace.Logger = trace.NewLogger(ctx.Trace())
		conf.LogLevel = aws.LogLevel(aws.LogDebug)
		conf.Logger = new(loggerWrap)

		client.Transport = http.NewTraceLoggingTransport(gohttp.DefaultTransport)
	}
	conf.HTTPClient = client

	return (*BaseConfig)(conf)
}
