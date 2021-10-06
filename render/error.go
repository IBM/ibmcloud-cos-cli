package render

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/urfave/cli"
)

type ErrorRender struct {
	terminal.UI
}

func NewErrorRender(terminal terminal.UI) *ErrorRender {
	tmp := new(ErrorRender)
	tmp.UI = terminal
	return tmp
}

// DisplayError print the error details in screen
// THE CONCEPT
//  For each specific type a new entry should be added in the switch
//  the switch is evaluated from top to down from left to right,
//  more specific types needs to be evaluated before more generic ones.
//  If type specific entry not found tries to safe cast to "CodeError" interface
//  and produce code driven messages, if safe cast fails fallback to just print error message.
// Still in POC state, most of errors use default strings coming from server
// a more complete error mapping will be added later to make plugin more globalization friendly,
// currently a few error mappings are done to show expected execution flow and hook points
func (er *ErrorRender) DisplayError(errorIn error) (err error) {
	var errorMessage string
	switch typeCheckedError := errorIn.(type) {
	case *errors.ObjectGetError:
		errorMessage = getMessageFromGetObjectError(typeCheckedError)
	case *errors.CommandError:
		errorMessage = getMessageFromCommandError(typeCheckedError)
	case errors.CodeError:
		errorMessage = getMessageByCodeError(typeCheckedError)
	default:
		errorMessage = errorIn.Error()
	}
	er.Failed(errorMessage)
	return
}

func getMessageFromGetObjectError(getObjectError *errors.ObjectGetError) string {
	overRideDefaultLocMessage := T("Override the default location providing an OUTFILE parameter")
	var message string
	switch getObjectError.Cause {
	case errors.IsDir:
		message = T("The download destination '{{.Location}}' is a directory.", getObjectError)
	case errors.Opening:
		message = T("Error opening '{{.Location}}' to write.", getObjectError)
	default:
		message = ""
	}
	if getObjectError.UsingDefaultRule && message != "" {
		message += "\n" + overRideDefaultLocMessage
	}
	return message
}

var commandErrorCausesStrings = map[errors.CommandErrorCause]string{
	errors.BadFlagSyntax:       T("Bad Flag Syntax in '--%s'"),
	errors.NotDefinedFlag:      T("Command does not support Flag '--%s'"),
	errors.InvalidBooleanValue: T("Invalid boolean balue in Flag '--%s'"),
	errors.InvalidBooleanFlag:  T("Invalid boolean flag '--%s'"),
	errors.MissingValue:        T("Flag '--%s' requires a value"),
	errors.InvalidValue:        T("The value in flag '--%s' is invalid"),
	errors.MissingRequiredFlag: T("Mandatory Flag '--%s' is missing"),
	errors.InvalidNArg:         "Unexpected number of arguments in command '%s'.",
	errors.InvalidDisplayValue: "Unsupported output format for command '%s', only ‘JSON’ and ‘TEXT’ are supported.",
}

func getMessageFromCommandError(commandError *errors.CommandError) string {
	message := commandErrorCausesStrings[commandError.Cause]
	switch commandError.Cause {
	case errors.InvalidNArg:
		message = fmt.Sprintf(message, commandError.CLIContext.Command.Name)
	case errors.InvalidDisplayValue:
		message = fmt.Sprintf(message, commandError.CLIContext.Command.Name)
	default:
		message = fmt.Sprintf(message, commandError.Flag)
	}

	buffer := bytes.NewBuffer([]byte{})
	currentWriter := commandError.CLIContext.App.Writer
	defer func() {
		commandError.CLIContext.App.Writer = currentWriter
	}()
	commandError.CLIContext.App.Writer = buffer
	cli.ShowCommandHelp(commandError.CLIContext, commandError.CLIContext.Command.Name)
	buffer.String()

	return message + "\n" + strings.TrimSpace(buffer.String())
}

// more can be added from
//https://cloud.ibm.com/docs/infrastructure/cloud-object-storage-infrastructure?topic=cloud-object-storage-infrastructure-common-error-codes
func getMessageByCodeError(errorIn errors.CodeError) string {
	switch errorIn.Code() {
	case "EmptyStaticCreds":
		return T("Try logging in using 'ibmcloud login'.")
	case "InvalidArgument":
		// If CRN is not registered in config or invalid, we send an actual invalid
		// argument message to request users to configure or re-configure CRN.
		// Otherwise, other errors are related to bad headers or any other invalid
		// arguments on the APIs.
		if !strings.Contains(errorIn.Error(), "Invalid Argument") {
			return errorIn.Error()
		}
		return T("Invalid Argument.  A valid service instance CRN must be configured to create or list buckets.  Verify the CRN using ‘ibmcloud cos config list’.")
	case "InvalidBucketName":
		return T("The specified bucket name is invalid. Bucket names must start and end in alphanumeric characters (from 3 to 63) and are limited to lowercase, numbers, non-consecutive dots, and hyphens.")
	case "BucketAlreadyExists":
		return T("The requested bucket name is not available. The bucket namespace is shared by all users of the system. Select a different name and try again.")
	case "AccessDenied":
		return T("Access to your IBM Cloud account was denied. Log in again by typing ibmcloud login --sso.")
	case "BucketAlreadyOwnedByYou":
		return T("A bucket with the specified name already exists in your account. Create a bucket with a new name.")
	case "NoSuchBucket":
		return T("The specified bucket was not found in your IBM Cloud account. This may be because you provided the wrong region. Provide the bucket's correct region and try again.")
	case "BucketNotEmpty":
		return T("The specified bucket is not empty. Delete all the files in the bucket, then try again.")
	case "EntityTooSmall":
		return T("Your proposed upload is smaller than the minimum allowed size. File parts must be greater than 5 MB in size, except for the last part.")
	case "NoSuchKey":
		return T("The specified key does not exist.")
	case "NoSuchWebsiteConfiguration":
		return T("The specified bucket does not have website configuration.")
	default:
		return errorIn.Error()
	}
}
