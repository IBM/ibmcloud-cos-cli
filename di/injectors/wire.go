//+build wireinject

package injectors

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws/endpoints"
	"github.com/IBM/ibmcloud-cos-cli/render"
	"github.com/google/wire"

	"github.com/IBM/ibmcloud-cos-cli/di/providers"

	"github.com/IBM/ibmcloud-cos-cli/utils"
)

func InitializeCosContext(_ plugin.PluginContext) (*utils.CosContext, error) {
	wire.Build(
		utils.CosContext{},
		providers.NewUI,
		render.NewTextRender,
		render.NewJSONRender,
		render.NewErrorRender,
		providers.GetS3APIFn,
		providers.GetDownloaderAPIFn,
		providers.GetUploaderAPIFn,
		providers.GetAsperaTransferFn,
		providers.GetPluginConfig,
		providers.NewSession,
		providers.NewConfig,
		providers.NewCOSEndPointsWSClient,
		wire.Bind(new(utils.ListKnownRegions), new(utils.COSEndPointsWSClient)),
		wire.Bind(new(endpoints.Resolver), new(utils.COSEndPointsWSClient)),
		providers.GetFileOperations,
		wire.Bind(new(utils.FileOperations), new(providers.FileOperationsImpl)),
		providers.GetBaseConfig,
	)

	return nil, nil
}
