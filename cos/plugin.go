package cos

import (
	"reflect"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/app"
	"github.com/IBM/ibmcloud-cos-cli/di/injectors"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/version"
	"github.com/urfave/cli"
)

// Cloud Object Storage Namespace as IBMCloud CLI subcommand
const (
	COSNameSpace = "cos"
)

// Plugin Struct
type Plugin struct{}

// Run directs the urfave CLI to route all the CLI's actions through urfave.
func (_ *Plugin) Run(context plugin.PluginContext, args []string) {

	// Define the namespace for this CLI App
	nameSpace := strings.Split(context.CommandNamespace(), " ")

	// Name of the CLI (cos)
	name := strings.Join([]string{context.CLIName(), nameSpace[0]}, " ")

	// Generate a new CLI App with the name
	cliApp := app.NewApp(name)

	// Initialize COS Context
	ctx, err := injectors.InitializeCosContext(context)

	// Error handling
	if err != nil {
		panic(err)
	}

	cliApp.Metadata = make(map[string]interface{})
	cliApp.Metadata[config.CosContextKey] = ctx

	cmd := append(nameSpace, args...)

	// CLI app is executed
	err = cliApp.Run(cmd)

	//if err != nil {
	//	panic(err)
	//}
}

// GetMetadata of the plugin
func (_ *Plugin) GetMetadata() plugin.PluginMetadata {
	// COS CLI App
	cliApp := app.NewApp(COSNameSpace)

	// Pass in the namespace and COS CLI App commands
	ns, cmds := mapCommands(cliApp.Commands, COSNameSpace)

	// Build Plugin Metadata for COS
	cosPlugin := plugin.PluginMetadata{
		Name:    "cloud-object-storage",
		Version: version.CLIVersion,
		Namespaces: append(ns, plugin.Namespace{
			Name:        COSNameSpace,
			Description: T("Interact with IBM Cloud Object Storage services"),
		}),
		Commands: cmds,
	}

	// Return COS Plugin Metadata
	return cosPlugin
}

// mapCommands iterates through all the commands supported for COS
func mapCommands(commands cli.Commands, ns string) (nsr []plugin.Namespace, cmdsr []plugin.Command) {
	//
	nsr, cmdsr = []plugin.Namespace{}, []plugin.Command{}
	// Iterate through list of commands to append command and its contents
	for _, command := range commands {
		nsItx, cmdsItx := mapCommand(command, ns)
		nsr = append(nsr, nsItx...)
		cmdsr = append(cmdsr, cmdsItx...)
	}
	return
}

// mapCommand iterates each command with its name and description under the parent
// Cloud Object Storage namespace
func mapCommand(command cli.Command, ns string) ([]plugin.Namespace, []plugin.Command) {
	//
	if len(command.Subcommands) > 0 {
		nsItx, cmdsItx := mapCommands(command.Subcommands, ns+" "+command.Name)
		namespace := plugin.Namespace{
			ParentName:  ns,
			Name:        command.Name,
			Description: command.Description,
		}
		return append(nsItx, namespace), cmdsItx
	}

	// Command definition with Name, Description, Usage, Flags and whether it is hidden or not
	cmd := plugin.Command{
		Namespace:   ns,
		Name:        command.Name,
		Description: command.Description,
		Usage:       command.Name + " " + command.UsageText + " " + command.ArgsUsage,
		Flags:       mapFlags(command.Flags),
		Hidden:      command.Hidden,
	}

	// Return blank plugin namespace and command
	return nil, []plugin.Command{cmd}
}

// Map through all the flags for the COS commands
func mapFlags(flags []cli.Flag) []plugin.Flag {
	// Build a list of flags per command
	result := make([]plugin.Flag, 0, len(flags))

	// Append list of flags to the command plugin list
	for _, flag := range flags {
		result = append(result, mapFlag(flag))
	}

	// Return list of flags for the command
	return result
}

// Map through each flag for a command
func mapFlag(flag cli.Flag) plugin.Flag {
	flagRfx := reflect.ValueOf(flag)

	description := ""
	usageRfx := flagRfx.FieldByName("Usage")
	if usageRfx.IsValid() {
		description = usageRfx.String()
	}

	valueRfx := flagRfx.FieldByName("Value")
	hasValue := valueRfx.IsValid()

	// Return the command flag with its name, description
	// and whether it requires value or not
	return plugin.Flag{
		Name:        flag.GetName(),
		Description: description,
		HasValue:    hasValue,
	}
}
