package functions

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

// DeepCopyIntoUsingJSON makes a deep copy from source to destination
// by marshaling the source in a temporary buffer and unmarshal back from the buffer to the destination
func DeepCopyIntoUsingJSON(destination, source interface{}) (err error) {
	var bts []byte
	if bts, err = json.Marshal(source); err != nil {
		return
	}
	err = json.Unmarshal(bts, destination)
	return
}

// GetCosContext boilerplate to retrieve the CosContext from metadata in a more correct way
func GetCosContext(cliContext *cli.Context) (*utils.CosContext, error) {
	// check the metadata map contains the cos context key
	if value, found := cliContext.App.Metadata[config.CosContextKey]; found {
		// if key found
		// tries type assertion to CosCostext and returns it
		if result, typeMatch := value.(*utils.CosContext); typeMatch {
			return result, nil
		}
		// if type assertion fails return error
		return nil, awserr.New("metadata.coscontext.invalidtype", "Metadata Value does not match expected Type", nil)

	}
	// if key not present return error
	return nil, awserr.New("metadata.coscontext.notfound", "Metadata Value does not constains CosContext Key", nil)
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
			// if not present return an error
			ce := new(errors.CommandError)
			ce.Flag = flagName
			ce.CLIContext = cliContext
			ce.Cause = errors.MissingRequiredFlag
			return ce
		}
		// for each field maps the flag content to S3 input field
		err := populateField(cliContext, flagName, destinationRflx, fieldName)
		// check if the maping returned an error
		// if error occurred stop processing the remaining fields/flags
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
			// check if the mapping returned an error
			// if error occurred stop processing the remaining fields/flags
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// populate field , grabs the value from the cli context and maps it to the S3 input
func populateField(cliContext *cli.Context,
	flagName string,
	destinationRflx reflect.Value,
	fieldName string) (err error) {
	// allocates a new error of type command error
	commandError := new(errors.CommandError)
	// set the command error context
	commandError.CLIContext = cliContext
	// sets the error flag name
	commandError.Flag = flagName
	// sets the cause to be invalid value
	commandError.Cause = errors.InvalidValue

	// retrieves the field by name using reflection
	fieldRflx := destinationRflx.FieldByName(fieldName)

	// checks the field exists
	if !fieldRflx.IsValid() {
		err = awserr.New("parameter.mapping.invalid", "Field '"+fieldName+"' does not exist", nil)
		return
	}

	// checks if the field is of kind interface
	if fieldRflx.Type().Kind() == reflect.Interface {
		// create a value of type io.ReadSeeker and gets it reflection
		readSeekerRflx := reflect.ValueOf(new(io.ReadSeeker)).Elem()
		// checks if the io.ReadSeeker inplements the type of the field interface
		if readSeekerRflx.Type().Implements(fieldRflx.Type()) {
			// allocates a var to hold the cos context
			var cosContext *utils.CosContext
			// populate the cos context variable and checks if the fetch had an error
			if cosContext, err = GetCosContext(cliContext); err != nil {
				// if an error occurred , propagate it and stop any further processing
				return
			}
			// allocate a io.ReadSeeker
			var readSeeker io.ReadSeeker
			// call the cos context wrapper of the os open operation
			// and checks if error occurred
			if readSeeker, err = cosContext.ReadSeekerCloserOpen(cliContext.String(flagName)); err != nil {
				// if error occurred
				// sets the command internal error to be the IO error
				// return the command error and stop processing
				commandError.IError = err
				err = commandError
				return
			}
			// if open succeeds, assign the open result to the struct field and exit function
			fieldRflx.Set(reflect.ValueOf(readSeeker))
			return
		}
	}

	// checks the kind of the field of the s3 input is map
	if fieldRflx.Kind() == reflect.Map {
		// if the kind is map , allocate a new pointer to map using reflection
		ptr := reflect.New(fieldRflx.Type())
		// set a reflection build map at the pointer destination
		ptr.Elem().Set(reflect.MakeMap(fieldRflx.Type()))
		// since used parseJSONinFile instead of parseJSON it alsos accept the format file://
		// tries to parse the content of the flag as a map
		err = parseJSONinFile(cliContext, ptr.Interface(), cliContext.String(flagName))
		// if no error thrown set map build using reflection as the s3 input field value
		if err == nil {
			fieldRflx.Set(ptr.Elem())
		} else {
			commandError.IError = err
			err = commandError
		}
		return
	}

	// gets the type of the reflected field
	fieldTypeRflx := fieldRflx.Type()
	// checks if the type is of kinf of pointer
	if fieldRflx.Kind() == reflect.Ptr {
		// if pinter uses the Elem() operation to get the pointed to Type
		fieldTypeRflx = fieldTypeRflx.Elem()
	}

	// uses reflection to build a new value of the type
	v := reflect.New(fieldTypeRflx)
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
		// if field is type pointer to int64
		*f, err = parseInt64(cliContext.String(flagName))
	case *int:
		// if field is type pointer to int
		var tmp int64
		tmp, err = parseInt64(cliContext.String(flagName))
		*f = int(tmp)
	case *time.Time:
		// if field type is pointer to time , parse the time value using parseTime
		*f, err = parseTime(cliContext.String(flagName))
	case *s3.Delete:
		// if type is pointer to Delete , use golang json decoder to map value
		err = parseJSONinFile(cliContext, f, cliContext.String(flagName))
	// case *s3.CreateBucketConfiguration:
	// 	err = parseJSONinFile(f, cliContext.String(flagName))
	case *s3.CORSConfiguration:
		// if type is pointer to CORSConfiguration , use golang json decoder to map value
		err = parseJSONinFile(cliContext, f, cliContext.String(flagName))
	case *s3.CompletedMultipartUpload:
		// if type is pointer to CompletedMultipartUpload , use golang json decoder to map value
		err = parseJSONinFile(cliContext, f, cliContext.String(flagName))
	case *s3.AccessControlPolicy:
		// if type is pointer to AccessControlPolicy, use golang json decoder to map value
		err = parseJSONinFile(cliContext, f, cliContext.String(flagName))
	default:
		panic("INVALID TYPE -- not mapped type yet")
	}

	// check if an error occurred
	if err != nil {
		// if and error occurred set it to be the command error internal error
		commandError.IError = err
		err = commandError
		// exit function to avoid further processing
		// and propagate the error to caller
		return
	}

	// check if the kind of the field is a pointer To Type
	if fieldRflx.Kind() != reflect.Ptr {
		// if not pointer to,
		// uses Elem to follow the value to the pointed value
		v = v.Elem()
	}

	// assigns value populated previous to the struct field
	fieldRflx.Set(v)

	return
}

// parseInt64 convert int64 value to string
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
func parseJSONinFile(cliContext *cli.Context, i interface{}, input string) (err error) {
	trimmed := strings.TrimSpace(input)
	var content string
	if strings.HasPrefix(trimmed, filePrefix) {
		fileName := trimmed[len(filePrefix):]
		var cosContext *utils.CosContext
		if cosContext, err = GetCosContext(cliContext); err != nil {
			return
		}
		var rsc utils.ReadSeekerCloser
		if rsc, err = cosContext.ReadSeekerCloserOpen(fileName); err != nil {
			return
		}
		defer rsc.Close()
		var fileBytes []byte
		if fileBytes, err = ioutil.ReadAll(rsc); err != nil {
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
	max      *int64 // the upper bound of the sum of the cardinality of all pages
	total    int64  // sum of cardinality of pages retrieved so far
	pageSize *int64 // the page size
	nextPage *int64 // the number of elements to request in next page
}

// NewPaginationHelper creates a new PaginationHelper from context using the specified flags
// return the PaginationHelper itself and a pointer to Page Size
func NewPaginationHelper(ctx *cli.Context, maxFlagName, pageSizeFlagName string) (*PaginationHelper, *int64, error) {
	// allocate a new command error
	commandError := new(errors.CommandError)
	// set command context to be current context
	commandError.CLIContext = ctx
	// set command cause to be invalid value
	commandError.Cause = errors.InvalidValue
	// allocate a new PaginationHelper
	pgHelper := &PaginationHelper{
		total: 0,
	}
	// check if the invocation sets the upper limit of all pages
	if ctx.IsSet(maxFlagName) {
		// if it is set parse it from flag and check for error on converting
		if value, err := parseInt64(ctx.String(maxFlagName)); err != nil {
			commandError.Flag = maxFlagName
			commandError.IError = err
			return nil, nil, commandError
		} else {
			// if set and valid set max of pagination helper
			pgHelper.max = &value
		}
	}
	// check if page size flag was set in the command invocation
	if ctx.IsSet(pageSizeFlagName) {
		// if it is set parse it from flag and check for error on converting
		if value, err := parseInt64(ctx.String(pageSizeFlagName)); err != nil {
			commandError.Flag = pageSizeFlagName
			commandError.IError = err
			return nil, nil, commandError
		} else {
			// if set and valid set pagesize of pagination helper
			pgHelper.pageSize = &value
		}
	}
	pgHelper.initNextPageSize()
	return pgHelper, pgHelper.nextPage, nil
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
