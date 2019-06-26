package app

import (
	"regexp"

	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/urfave/cli"
)

var (
	BadFlagSyntaxRegex       = regexp.MustCompile(`bad flag syntax: ([\w\-]+)$`)
	NotDefinedFlagRegex      = regexp.MustCompile(`flag provided but not defined: -([\w\-]+)$`)
	InvalidBooleanValueRegex = regexp.MustCompile(`invalid boolean value .+? for -([\w\-]+): .+$`)
	InvalidBooleanFlagRegex  = regexp.MustCompile(`invalid boolean flag ([\w\-]+): .+$`)
	MissingValueRegex        = regexp.MustCompile(`flag needs an argument: -([\w\-]+)$`)
	InvalidValueRegex        = regexp.MustCompile(`invalid value ".+?" for flag -([\w\-]+): .+$`)
)

func OnUsageError(context *cli.Context, err error, _ bool) error {
	result := new(errors.CommandError)
	result.CLIContext = context
	result.IError = err
	switch errorStr := err.Error(); {
	case BadFlagSyntaxRegex.MatchString(errorStr):
		result.Cause = errors.BadFlagSyntax
		result.Flag = BadFlagSyntaxRegex.FindStringSubmatch(errorStr)[1]
	case NotDefinedFlagRegex.MatchString(errorStr):
		result.Cause = errors.NotDefinedFlag
		result.Flag = NotDefinedFlagRegex.FindStringSubmatch(errorStr)[1]
	case InvalidBooleanValueRegex.MatchString(errorStr):
		result.Cause = errors.InvalidBooleanValue
		result.Flag = InvalidBooleanValueRegex.FindStringSubmatch(errorStr)[1]
	case InvalidBooleanFlagRegex.MatchString(errorStr):
		result.Cause = errors.InvalidBooleanFlag
		result.Flag = InvalidBooleanFlagRegex.FindStringSubmatch(errorStr)[1]
	case MissingValueRegex.MatchString(errorStr):
		result.Cause = errors.MissingValue
		result.Flag = MissingValueRegex.FindStringSubmatch(errorStr)[1]
	case InvalidValueRegex.MatchString(errorStr):
		result.Cause = errors.InvalidValue
		result.Flag = InvalidValueRegex.FindStringSubmatch(errorStr)[1]
	default:
		return err
	}
	return result
}
