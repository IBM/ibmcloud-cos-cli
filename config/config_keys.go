package config

import (
	"path/filepath"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/config_helpers"
)

// Configuration persistence keys
const (
	// CRN constant
	CRN = "CRN"

	// Last time the config is updated
	LastUpdated = "Last Updated"

	// Download Location constant
	DownloadLocation = "Download Location"

	// Region constant
	DefaultRegion = "Default Region"

	// HMAC authentication constants
	HMACProvided              = "HMACProvided"
	LabelAuthenticationMethod = "Authentication Method"
	AccessKeyID               = "AccessKeyID"
	SecretAccessKey           = "SecretAccessKey"

	// Regions Endpoint URL constant
	RegionsEndpointURL = "RegionsEndpointURL"

	// ForcePathStyle
	ForcePathStyle = "ForcePathStyle"
	LabelURLStyle  = "URL Style"
)

// CLI App Context Metadata Keys
const (
	CosContextKey = "CosContext"
)

var (
	// Variables to use across the package
	FallbackRegion           = "us-geo"
	FallbackDownloadLocation = filepath.Join(config_helpers.UserHomeDir(), "Downloads")

	// Standard time format
	StandardTimeFormat = "Monday, January 02 2006 at 15:04:05"
)

const (
	HMACProvidedDefault = false
	// Current Authentication Methods the CLI supports
	IAM  = "IAM"
	HMAC = "HMAC"
)

const (
	ForcePathStyleDefault = false
	// Current Bucket URL Styles that the CLI supports
	VHost = "VHost"
	Path  = "Path"
)
