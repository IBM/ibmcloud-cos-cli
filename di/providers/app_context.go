//+build !unit

package providers

import (
	"io"
	gohttp "net/http"
	"os"
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

	conf.WithS3ForcePathStyle(true)
	conf.DisableRestProtocolURICleaning = aws.Bool(true)

	hmac, _ := ctx.PluginConfig().GetBoolWithDefault(config.HMACProvided, false)

	if hmac {

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

func (_ *FileOperationsImpl) WriteCloserOpen(location string) (io.WriteCloser, error) {
	return os.Create(location)
}

func (_ *FileOperationsImpl) GetFileInfo(location string) (os.FileInfo, error) {
	return os.Stat(location)
}

func (_ *FileOperationsImpl) Remove(location string) error {
	return os.Remove(location)
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

	trace.Logger = trace.NewLogger(ctx.Trace())

	conf := new(aws.Config)
	conf.LogLevel = aws.LogLevel(aws.LogDebug)
	conf.Logger = new(loggerWrap)

	client := &gohttp.Client{
		Transport: http.NewTraceLoggingTransport(gohttp.DefaultTransport),
		Timeout:   time.Duration(ctx.HTTPTimeout()) * time.Second,
	}

	conf.HTTPClient = client

	return (*BaseConfig)(conf)
}
