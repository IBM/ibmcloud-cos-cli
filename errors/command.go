//go:generate stringer -type=CommandErrorCause

package errors

import (
	"fmt"

	"github.com/urfave/cli"
)

type CommandErrorCause int

const CommandErrorCode = "CommandError"

// Enum Failing flag parse
const (
	_ CommandErrorCause = iota
	BadFlagSyntax
	NotDefinedFlag
	InvalidBooleanValue
	InvalidBooleanFlag
	MissingValue
	InvalidValue
	MissingRequiredFlag
	InvalidNArg
)

type CommandError struct {
	IError     error
	CLIContext *cli.Context
	Flag       string
	Cause      CommandErrorCause
}

func (_ *CommandError) Code() string {
	return CommandErrorCode
}

func (ce *CommandError) Error() string {
	if ce.IError != nil {
		return ce.IError.Error()
	}
	return fmt.Sprintf("%s: flag='%s' cause=%s", ce.Code(), ce.Flag, ce.Cause)
}
