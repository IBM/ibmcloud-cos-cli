package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/ibm-cos-sdk-go/aws"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// FormatFileSize function lift from
// https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
// outputs a human readable representation of the value using multiples of 1024
func FormatFileSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// FormatErrorMessage formats some common error messages that the Amazon AWS SDK returns, that users may find when
// using the CLI. It is incomplete, but if an error is not accounted for here, the full non-formatted error will be
// printed.
func FormatErrorMessage(errorString string) string {
	switch {
	case strings.Contains(errorString, "InvalidBucketName"):
		return T("The specified bucket name is invalid. Bucket names must start and end in alphanumeric characters (from 3 to 63) and are limited to lowercase, numbers, non-consecutive dots, and hyphens.")
	case strings.Contains(errorString, "BucketAlreadyExists"):
		return T("The requested bucket name is not available. The bucket namespace is shared by all users of the system. Select a different name and try again.")
	case strings.Contains(errorString, "AccessDenied"):
		return T("Access to your IBM Cloud account was denied. Log in again by typing ibmcloud login --sso.")
	case strings.Contains(errorString, "BucketAlreadyOwnedByYou"):
		return T("A bucket with the specified name already exists in your account. Create a bucket with a new name.")
	case strings.Contains(errorString, "NoSuchBucket"):
		return T("The specified bucket was not found in your IBM Cloud account. This may be because you provided the wrong region for the bucket to delete. Provide the bucket's correct region and try again.")
	case strings.Contains(errorString, "BucketNotEmpty"):
		return T("The specified bucket is not empty. Delete all the files in the bucket, then try again.")
	case strings.Contains(errorString, "EntityTooSmall"):
		return T("Your proposed upload is smaller than the minimum allowed size. File parts must be greater than 5 MB in size, except for the last part.")
	case strings.Contains(errorString, "NoSuchKey"):
		return T("The specified bucket was not found in your account. Ensure that you have set the correct region with the region flag.")
	default:
		return errorString
	}
}

// GetRegion helper to reduce the boilerplate of getting the region
func GetRegion(cliContext *cli.Context, pluginConfig plugin.PluginConfig) (string, error) {
	// checks if current cli contest has region flag set
	if cliContext.IsSet(flags.Region) {
		// if flag is set returns its value
		return cliContext.String(flags.Region), nil
	}
	// if falg is not set
	// Gets config default region
	// if config default region is not set returns cli login region
	// if cli login region is not set fall back to us ( -geo )
	return pluginConfig.GetStringWithDefault(config.DefaultRegion, config.FallbackRegion)

}

// GetCosContext boilerplate to retrieve the CosContext from metadat in a more correct way
func GetCosContext(cliContext *cli.Context) (*utils.CosContext, error) {
	// check the metdata map contais the cos context key
	value, ok := cliContext.App.Metadata[config.CosContextKey]
	// if key present
	// tries type assertion to CosCostext and returns it
	if ok {
		result, ok := value.(*utils.CosContext)
		if ok {
			return result, nil
		}
		// if type assertion fails return error
		return nil, errors.New("metadata.coscontext.invalidtype")

	}
	// if key not present return error
	return nil, errors.New("metadata.coscontext.notfound")
}

// MapToSDKInput - validates the user flags and their contents against
// SDK interfaces to ensure types in the parameters passed in are valid
// Parameters:
//		CLI Context Application
//		Destination Interface (API Inputs such as DeleteObjectsInput)
//		Required Parameters for a CLI command
//		Optional Parameters for a CLI command
// Return:
//		Error - nil or a valid error
func MapToSDKInput(cliContext *cli.Context, destination interface{},
	mandatoryFields map[string]string, optionalFields map[string]string) error {

	// gets a new reflect value from the original value
	destinationRflx := reflect.ValueOf(destination)

	// check the kind of the relfect value,
	if destinationRflx.Kind() == reflect.Ptr {
		// follow the pointer to the actual value pointed to
		destinationRflx = destinationRflx.Elem()
	} else {
		// panics if not a pointer to ...
		panic("not a pointer ... ")
	}

	// Iterate all through mandatory fields
	// if the value is mandatory and not present return error
	for fieldName, flagName := range mandatoryFields {
		// check the cli context contains the mandatory flag
		if !cliContext.IsSet(flagName) {
			// if not presnet return an error
			return errors.New("missing flag " + flagName)
		}
		// for each field maps the flag content to S3 input field
		err := populateField(cliContext, flagName, destinationRflx, fieldName)
		// check if the maping returned an error
		// if error ocurred stop processing the remaining fields/flags
		if err != nil {
			return err
		}
	}

	// Iterate all through optional fields
	for fieldName, flagName := range optionalFields {
		// check if the flag is populated
		if cliContext.IsSet(flagName) {
			// if flag populated
			// maps the flag content to the S3 input field
			err := populateField(cliContext, flagName, destinationRflx, fieldName)
			// check if the maping returned an error
			// if error ocurred stop processing the remaining fields/flags
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateUserInputsAndSetRegion helper to reduce the boilerplate of getting the region and mapping the fields
func ValidateUserInputsAndSetRegion(c *cli.Context, input interface{}, mandatory map[string]string,
	options map[string]string, p plugin.PluginConfig) (string, error) {

	// apply the mandatory and optional mappings to the S3 input using the flags content as source
	err := MapToSDKInput(c, input, mandatory, options)
	// if mapping fails, return an incorrect usage
	if err != nil {
		return "", cli.NewExitError(incorrectUsage, 1)
	}

	// Set region from either user's input or default region in conf
	region, err := GetRegion(c, p)
	if err != nil || region == "" {
		return "", cli.NewExitError(noRegion, 1)
	}

	return region, nil
}

// populate field , grabs the value from the cli context and maps it to the S3 input
func populateField(cliContext *cli.Context,
	flagName string,
	destinationRflx reflect.Value,
	fieldName string) (err error) {
	fieldRflx := destinationRflx.FieldByName(fieldName)

	// checks the kind of the field of the s3 input is map
	if fieldRflx.Kind() == reflect.Map {
		// if the kind is map , allocate a new pointer to map using reflection
		ptr := reflect.New(fieldRflx.Type())
		// set a reflection build map at the pointer destination
		ptr.Elem().Set(reflect.MakeMap(fieldRflx.Type()))
		// since used parseJSONinFile instead of parseJSON it alsos accept the format file://
		// tries to parse the content of the flag as a map
		err = parseJSONinFile(ptr.Interface(), cliContext.String(flagName))
		// if no error thrown set map build using reflection as the s3 input field value
		if err == nil {
			fieldRflx.Set(ptr.Elem())
		}
		return
	}

	// checks the kind of the field of the s3 input is pointer to ...
	if fieldRflx.Kind() != reflect.Ptr {
		panic("assumed all sdk struct fields are pointers, something is not ok")
	}
	// grabs the type that the pointer points to
	fieldPointerToType := fieldRflx.Type().Elem()
	// uses reflection to build a new value of the type
	v := reflect.New(fieldPointerToType)
	// switch over the types,
	// each type as its mapper
	switch f := v.Interface().(type) {
	case *string:
		// if field is type pointer to string , uses directly the value of the flag
		*f = cliContext.String(flagName)
	case *bool:
		// if field is type pointer to bool , uses directly the value of the flag
		*f = cliContext.Bool(flagName)
	case *int64:
		// if field is type pointer to int64 , uses directly the value of the flag
		*f = cliContext.Int64(flagName)
	case *time.Time:
		// if field type is pointer to time , parse the time value using parseTime
		*f, err = parseTime(cliContext.String(flagName))
	case *s3.Delete:
		// if type is pointer to Delete , use golang json decoder to map value
		err = parseJSONinFile(f, cliContext.String(flagName))
	// case *s3.CreateBucketConfiguration:
	// 	err = parseJSONinFile(f, cliContext.String(flagName))
	case *s3.CORSConfiguration:
		// if type is pointer to CORSConfiguration , use golang json decoder to map value
		err = parseJSONinFile(f, cliContext.String(flagName))
	case *s3.CompletedMultipartUpload:
		// if type is pointer to CompletedMultipartUpload , use golang json decoder to map value
		err = parseJSONinFile(f, cliContext.String(flagName))
	case *s3.AccessControlPolicy:
		// if type is pointer to AccessControlPolicy, use golang json decoder to map value
		err = parseJSONinFile(f, cliContext.String(flagName))
	default:
		panic("INVALID TYPE -- not mapped type yet")
	}

	if err == nil {
		fieldRflx.Set(v)
	}

	return
}

// Convert int64 value to string
func parseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

// parseTime parses different timestamp formats we support
func parseTime(value string) (tm time.Time, err error) {
	// List of supported timestamp formats under ISO8601
	timeFormats := []string{
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		"2006-01-02",
		"Monday, January 02, 2006 at 15:04:05",
		"Monday, January 2, 2006 at 15:04:05",
	}

	// Check if the value passed is a UNIX / Epoch based timestamp
	i, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		tm = time.Unix(i, 0)
		return
	}

	// Iterate through each time format to find a match
	for _, format := range timeFormats {
		tm, err = time.Parse(format, value)
		if err == nil {
			return
		}
	}

	// Return error
	return
}

// parseJSON - parses JSON input user provides
func parseJSON(i interface{}, input string) (err error) {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) > 0 {
		if trimmed[0] == '{' {
			err = json.Unmarshal([]byte(input), i)
		} else {
			var jsonBlob []byte
			jsonBlob, err = utils.GetAsJson(input)
			if err != nil {
				return err
			}
			err = json.Unmarshal(jsonBlob, i)
		}
	}
	return
}

// Constant prefix when opening a user input file
const filePrefix = `file://`

// parseJSONinFile - parse JSON user provides inside a file - collateral
// effect consumes files in the simplified version of json ...
func parseJSONinFile(i interface{}, input string) (err error) {
	trimmed := strings.TrimSpace(input)
	var content string
	if strings.HasPrefix(trimmed, filePrefix) {
		fileName := trimmed[len(filePrefix):]
		// gets a new one, not shared with the app context in production
		// in testing is shared with across app due to implementation details and being easier to assert on calling
		fileOp := providers.GetFileOperations()
		var rsc utils.ReadSeekerCloser
		rsc, err = fileOp.ReadSeekerCloserOpen(fileName)
		if err != nil {
			return
		}
		defer rsc.Close()
		var fileBytes []byte
		fileBytes, err = ioutil.ReadAll(rsc)
		if err != nil {
			return
		}
		content = string(fileBytes)
	} else {
		content = trimmed
	}
	return parseJSON(i, content)
}

/// /// ///
/// Pagination Helper ///
/// /// ///

// PaginationHelper struct to hep with paging handling
type PaginationHelper struct {
	max      *int64
	total    int64
	pageSize *int64
	nextPage *int64
}

// NewPaginationHelper creates a new PaginationHelper from context using the specified flags
// return the PaginationHelper itself and a pointer to Page Size
func NewPaginationHelper(ctx *cli.Context, maxFlagName, pageSizeFlagName string) (*PaginationHelper, *int64) {
	pgHelper := &PaginationHelper{
		total: 0,
	}
	if ctx.IsSet(maxFlagName) {
		pgHelper.max = aws.Int64(ctx.Int64(maxFlagName))
	}
	if ctx.IsSet(pageSizeFlagName) {
		pgHelper.pageSize = aws.Int64(ctx.Int64(pageSizeFlagName))
	}
	pgHelper.initNextPageSize()
	return pgHelper, pgHelper.nextPage
}

// UpdateTotal updates total return if max reached
func (p *PaginationHelper) UpdateTotal(delta int) bool {
	p.total += int64(delta)
	p.refreshNextPageSize()
	return p.Continue()
}

// Continue checks if max reached and more pages can be request
func (p *PaginationHelper) Continue() bool {
	if p.max != nil {
		return *p.max > p.total
	}
	return true
}

// initNextPageSize sets the initial page size
func (p *PaginationHelper) initNextPageSize() {
	// when maxitems not set and pagesize not set
	// do not set current request value, use server defaults
	if p.max == nil && p.pageSize == nil {
		p.nextPage = nil
		return
	}
	// when pagesize is set
	// page size will be upper boundary of the number of items to get in current request
	if p.max == nil && p.pageSize != nil {
		p.nextPage = aws.Int64(*p.pageSize)
		return
	}
	// when no page size set, but max items defined pass the number of missing items as upper boundary
	if p.max != nil && p.pageSize == nil {
		p.nextPage = aws.Int64(*p.max)
		return
	}
	// when pagesize and max request set
	// use the mininum between both as upper boundary
	if *p.max > *p.pageSize {
		p.nextPage = aws.Int64(*p.pageSize)
	} else {
		p.nextPage = aws.Int64(*p.max)
	}
}

// refreshNextPageSize adjust the size of the page
func (p *PaginationHelper) refreshNextPageSize() {
	if p.nextPage != nil && p.max != nil {
		remaining := *p.max - p.total
		if *p.nextPage > remaining {
			*p.nextPage = remaining
		}
	}
}

// GetTotal returns the total number
func (p *PaginationHelper) GetTotal() int64 {
	return p.total
}

// GetNextPage returns the total number
func (p *PaginationHelper) GetNextPage() int64 {
	return *p.nextPage
}
