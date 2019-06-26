package app

import (
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/version"
	"github.com/urfave/cli"
)

// NewApp Object for the CLI
// Parameter:
//     Name: Name of the App (COS)
// Returns:
//     CLI Application
func NewApp(name string) *cli.App {

	// Generation of the Application with
	// App - CLI Application
	// Name - COS
	// HelpName - COS
	// Version of App - found in /version/version.go
	// Usage Error calls
	// Discard Writer
	app := cli.NewApp()
	app.Name = name
	app.HelpName = name
	app.Version = version.CLIVersion.String()
	app.OnUsageError = OnUsageError
	app.Writer = ioutil.Discard

	// Template to factorize the help section of the commands
	cli.CommandHelpTemplate = CommandHelpTemplate

	// Commands that this CLI Application supports
	app.Commands = cli.Commands{
		commands.CommandBucketCreate,
		commands.CommandBucketDelete,
		commands.CommandBucketGetLocation,
		commands.CommandBucketGetClass,
		commands.CommandBucketHead,
		commands.CommandBucketsList,
		commands.CommandBucketsListExtended,
		commands.CommandObjectGet,
		commands.CommandObjectHead,
		commands.CommandObjectPut,
		commands.CommandObjectDelete,
		commands.CommandObjectDeletes,
		commands.CommandObjectCopy,
		commands.CommandObjectList,
		commands.CommandMPUCreate,
		commands.CommandMPUPartUpload,
		commands.CommandMPUPartUploadCopy,
		commands.CommandMPUPartsList,
		commands.CommandMPUList,
		commands.CommandMPUAbort,
		commands.CommandMPUComplete,
		commands.CommandConfig,
		commands.CommandBucketCorsPut,
		commands.CommandBucketCorsGet,
		commands.CommandBucketCorsDelete,
		commands.CommandDownload,
		commands.CommandUpload,
		commands.CommandWait,
	}

	// Runs to set Usage Text for the flags of all commands
	setUsageText(app.Commands)

	// Return the CLI application
	return app
}

// Set the usage text from the flags for all commands
func setUsageText(commands cli.Commands) {
	// Iterate through commands and establish usage for each
	for idx := range commands {
		commands[idx].UsageText = fromFlagsToUsage(commands[idx].Flags)
		commands[idx].OnUsageError = OnUsageError
		setUsageText(commands[idx].Subcommands)
	}
}

// Set the usage text from the flags
func fromFlagsToUsage(flags []cli.Flag) string {
	// Build a list to contain flag names
	flagsNames := make([]string, 0, len(flags))

	// Iterate through flags and append flag by flag to the list
	for _, flag := range flags {
		flagsNames = append(flagsNames, getFlagUse(flag))
	}
	// Returns the list of flag names
	return strings.Join(flagsNames, " ")
}

// Usage test for the flags
func getFlagUse(flag cli.Flag) string {
	// Reflect of flags hidden or not
	flagRfx := reflect.ValueOf(flag)
	hiddenRfx := flagRfx.FieldByName("Hidden")

	// Initialize flag name place holder
	flagNamePlaceHolder := ""

	// Build descriptions of the flags
	flagDesc := strings.SplitN(flag.String(), "\t", 2)
	if len(flagDesc) > 0 {
		flagNamePlaceHolder = flagDesc[0]
	}

	// Set hidden brackets for the flags
	if hiddenRfx.IsValid() && hiddenRfx.Bool() {
		return flagNamePlaceHolder
	} else {
		return "[" + flagNamePlaceHolder + "]"
	}
}
