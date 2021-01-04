package version

import "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

// CLIName is the name of the IBM Cloud
// Object Storage CLI
const CLIName = "ibm-cloud-cos-cli"

// CLIVersion is the version of the CLI
var CLIVersion = plugin.VersionType{
	Major: 1,
	Minor: 2,
	Build: 2,
}
