package functions

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/IBM/ibmcloud-cos-cli/config/fields"

	"github.com/IBM/ibmcloud-cos-cli/config/flags"

	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibmcloud-cos-cli/config"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"

	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/urfave/cli"
)

// BucketClass allows the user to get the Class of a specific bucket
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketClass(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI

	// Initialize GetBucketLocationInput
	input := new(s3.GetBucketLocationInput)

	// Required parameter and no optional parameter for GetBucketLocation
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Check through user inputs for validation
	err := MapToSDKInput(c, input, mandatory, map[string]string{})
	if err != nil {
		ui.Failed(T("Incorrect Usage."))
		cli.ShowCommandHelp(c, c.Command.Name)
		ui.Say("")
		return cli.NewExitError("", 1)
	}

	// Operate Bucket Location Coordinator
	output, err := getBucketLocationCoordinator(cosContext, input)
	// Error handling
	if err != nil {
		if strings.Contains(err.Error(), "EmptyStaticCreds") {
			ui.Failed(err.Error() + "\n" + T("Try logging in using 'ibmcloud login'."))
		} else {
			ui.Failed(err.Error())
		}
		return cli.NewExitError("", 1)
	}

	// Decode the location from the Location Constraint
	_, class := locationDecoder(aws.StringValue(output.LocationConstraint))

	// Success
	ui.Ok()

	// Output the successful message
	ui.Say(T("Details about bucket {{.details}}:",
		map[string]interface{}{"details": utils.EntityNameColor(aws.StringValue(input.Bucket))}))
	ui.Say(T("Class: {{.class}}", map[string]interface{}{"class": utils.EntityNameColor(Class(class).String())}))

	// Return
	return nil
}

// BucketLocation allows the user to get the location of a specific bucket given the bucket name and the region.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func BucketLocation(c *cli.Context) error {

	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI

	// Initialize GetBucketLocationInput
	input := new(s3.GetBucketLocationInput)

	// Required parameter for GetBucketLocation
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Required parameter for GetBucketLocation
	options := map[string]string{}

	// Check through user inputs for validation
	err := MapToSDKInput(c, input, mandatory, options)
	if err != nil {
		ui.Failed(T("Incorrect Usage."))
		cli.ShowCommandHelp(c, c.Command.Name)
		ui.Say("")
		return cli.NewExitError("", 1)
	}

	// Execute the GetBucketLocation COordinator to get class
	output, err := getBucketLocationCoordinator(cosContext, input)
	// Error Handling
	if err != nil {
		if strings.Contains(err.Error(), "EmptyStaticCreds") {
			ui.Failed(err.Error() + "\n" + T("Try logging in using 'ibmcloud login'."))
		} else {
			ui.Failed(err.Error())
		}
		return cli.NewExitError("", 1)
	}

	// Decode the Location Constraint to obtain class
	region, class := locationDecoder(aws.StringValue(output.LocationConstraint))

	// Success
	ui.Ok()

	// Output the successful message
	ui.Say(T("Details about bucket {{.details}}:",
		map[string]interface{}{"details": utils.EntityNameColor(aws.StringValue(input.Bucket))}))
	ui.Say(T("Region: {{.region}}", map[string]interface{}{"region": utils.EntityNameColor(region)}))
	ui.Say(T("Class: {{.class}}", map[string]interface{}{"class": utils.EntityNameColor(Class(class).String())}))

	// Return
	return nil

}

// GetLocationWrapper builds a struct for GetBucketLocationOuput and Error
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

	// Set background channell for GetLocationWrapper
	channel := make(chan GetLocationWrapper)

	// Initialize waitgroup for current known COS regions
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

// GetBucketLocatioWorker is assigned to grab the assigned region of the bucket
func getBucketLocationWorker(ctx context.Context, region string, input *s3.GetBucketLocationInput,
	clientGen func(string) s3iface.S3API, resultChannel chan<- GetLocationWrapper, waitGroup *sync.WaitGroup) {

	var getLocationWrapper GetLocationWrapper

	defer waitGroup.Done()

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
	client := clientGen(region)

	// Make a GetBucketLocation request with context and input
	output, err := client.GetBucketLocationWithContext(ctx, input)

	// Assign output and erro to global getLocationWrapper variables
	getLocationWrapper.Result = output
	getLocationWrapper.Error = err

	// Pass wrapper to the channel
	resultChannel <- getLocationWrapper
}

// locationDecoder - decodes region/class from the location constraint string
func locationDecoder(locationConstraint string) (string, string) {
	// Regex to find a region match
	regionDetails := utils.RegionDecoderRegex.FindStringSubmatch(locationConstraint)
	if regionDetails != nil {
		region, geo, class := regionDetails[1], regionDetails[2], regionDetails[3]
		return region + geo, class
	}

	// Return empty if not found
	return "", ""
}
