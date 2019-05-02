package utils

import (
	"github.com/IBM/ibm-cos-sdk-go/aws/request"
	"github.com/IBM/ibmcloud-cos-cli/version"
)

// User Agent Handler
var CLIVersionUserAgentHandler = request.NamedHandler{
	Name: "plugin.CliVersionUserAgentHandler",
	Fn:   AddCLIVersionToUserAgent,
}

// AddCLIVersionToUserAgent - Build a request to include
// CLI in the header
func AddCLIVersionToUserAgent(r *request.Request) {
	// Concatentate both CLI name and version
	cliUAgent := version.CLIName + "/" + version.CLIVersion.String()
	// Grab the HTTP request header's User-Agent
	uAgent := r.HTTPRequest.Header.Get("User-Agent")
	if len(uAgent) > 0 {
		uAgent = cliUAgent + ", " + uAgent
	}
	// Set the CLI into User Agent header
	r.HTTPRequest.Header.Set("User-Agent", uAgent)
}
