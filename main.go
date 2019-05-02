package main

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibmcloud-cos-cli/cos"
)

// Start of the COS CLI Plugin Main program
func main() {
	// Start starts the plugin -
	// Calling to the IBM Cloud CLI SDK
	plugin.Start(new(cos.Plugin))
}
