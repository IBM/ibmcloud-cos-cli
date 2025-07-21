package functions

import (
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/render"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// Endpoints - show regions
// Parameter:
//
//	CLI Context Application
//
// Returns:
//
//	Error = zero or non-zero
func Endpoints(c *cli.Context) (err error) {
	if c.NArg() > 0 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		return
	}

	cosContext, err := GetCosContext(c)
	if err != nil {
		return
	}

	flagRegion := c.String(flags.Region)

	region, err := cosContext.GetCurrentRegion(flagRegion)
	if err != nil || region == "" {
		return &errors.EndpointsError{
			Cause: errors.EndpointsRegionEmpty,
		}
	}

	ibmEndpoints, err := utils.NewIBMEndPoints(cosContext.Session.Config, "")
	if err != nil {
		return err
	}

	var output render.RegionEndpointsOutput
	if c.IsSet(flags.ListRegions) {
		regions, err := ibmEndpoints.NewGetAllRegions()
		if err != nil {
			return err
		}

		output.Regions = regions
		return cosContext.GetDisplay(c.String(flags.Output), false).Display(nil, &output, nil)
	}

	output, err = ibmEndpoints.GetAllEndpointsFor(s3.ServiceName, region)
	if err != nil {
		var cause errors.EndpointsErrorCause
		if flagRegion == "" {
			cause = errors.EndpointsStoredRegionInvalid
		} else {
			cause = errors.EndpointsRegionInvalid
		}
		return &errors.EndpointsError{
			Region: region,
			Cause:  cause,
		}
	}

	return cosContext.GetDisplay(c.String(flags.Output), false).Display(nil, &output, nil)
}
