package functions

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/ibmcloud-cos-cli/errors"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/render"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// BucketClassLocation - retrieves class or location of a certain bucket
// Parameter:
//     	CLI Context Application
// Returns:
//  	Error if triggered
func BucketClassLocation(c *cli.Context) (err error) {
	// check the number of arguments
	if c.NArg() > 0 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		return
	}

	// Load COS Context
	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}

	// Initialize GetBucketLocationInput
	input := new(s3.GetBucketLocationInput)

	// Required parameter and no optional parameter for GetBucketLocation
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Optional parameters
	options := map[string]string{}

	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	// Operate Bucket Location Coordinator
	var output *s3.GetBucketLocationOutput
	output, err = getBucketLocationCoordinator(cosContext, input)
	if err != nil {
		return
	}
	// Render processing for class or location of bucket
	var result interface{} = output
	if c.Command.Name == "get-bucket-class" {
		result = (*render.GetBucketClassOutput)(output)
	}

	// Display output
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, result, nil)

	// Return
	return
}

// GetLocationWrapper builds a struct for GetBucketLocationOutput and Error
type GetLocationWrapper struct {
	Result *s3.GetBucketLocationOutput
	Error  error
}

// GetBucketLocationCoordinator handles allocations of locations assigned to workers
func getBucketLocationCoordinator(cosContext *utils.CosContext,

	// Set the input of GetBucketLocation
	input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error) {

	// Populate listKnownRegions with supported COS regions
	listKnownRegions, err := cosContext.ListKnownRegions.GetAllRegions()
	if err != nil {
		return nil, err
	}

	// Cancel context/function set when running in background
	cancelAbleCtx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	// Set background channel for GetLocationWrapper
	channel := make(chan GetLocationWrapper)

	// Initialize wait group for current known COS regions
	var wg sync.WaitGroup
	wg.Add(len(listKnownRegions))

	// Iterate through each known region and run a GetBucketLocation thread on it
	for _, region := range listKnownRegions {
		go getBucketLocationWorker(cancelAbleCtx, region, input, cosContext.GetClient, channel, &wg)
	}

	// Execute with the wait group enabled
	go func() {
		wg.Wait()
		close(channel)
	}()

	// Range through each channel to parse result
	for getLocationWrapper := range channel {
		if getLocationWrapper.Error == nil {
			return getLocationWrapper.Result, nil
		}
		err = getLocationWrapper.Error
	}

	// Return
	return nil, err
}

// GetBucketLocationWorker is assigned to grab the assigned region of the bucket
func getBucketLocationWorker(ctx context.Context, region string, input *s3.GetBucketLocationInput,
	clientGen func(string) (s3iface.S3API, error), resultChannel chan<- GetLocationWrapper, waitGroup *sync.WaitGroup) {

	// Wrapper of Get Bucket Location
	var getLocationWrapper GetLocationWrapper

	// Defer WG to be done
	defer waitGroup.Done()

	// Defer for recovery
	defer func() {
		recover := recover()
		if nil != recover {
			if err, ok := recover.(error); ok {
				getLocationWrapper.Error = err
			} else {
				getLocationWrapper.Error = fmt.Errorf("unexpected panic : %v", recover)
			}
			resultChannel <- getLocationWrapper
		}
	}()

	// Generate client in the region
	var client s3iface.S3API
	client, getLocationWrapper.Error = clientGen(region)
	if getLocationWrapper.Error == nil {
		// Make a GetBucketLocation request with context and input
		getLocationWrapper.Result, getLocationWrapper.Error = client.GetBucketLocationWithContext(ctx, input)
	}

	// Pass wrapper to the channel
	resultChannel <- getLocationWrapper
}
