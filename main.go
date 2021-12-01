package main

import (
	"os"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/urfave/cli"
)

// Start of the COS CLI Plugin Main program
func main() {
	// var to hold exit code
	exitCode := 0

	// urfave-cli process some exit codes before bubble them up,
	// this wrapper will help the exit error bubble up until the error render
	cli.OsExiter = wrapExiter(&exitCode)
	// Start starts the plugin -
	// Calling to the IBM Cloud CLI SDK
	plugin.Start(new(cos.Plugin))

	// Exit with an exit code
	os.Exit(exitCode)
}

func wrapExiter(exitCodeHolder *int) func(int) {
	return func(exitCode int) {
		if exitCodeHolder != nil && *exitCodeHolder == 0 {
			*exitCodeHolder = exitCode
		}
	}
}
