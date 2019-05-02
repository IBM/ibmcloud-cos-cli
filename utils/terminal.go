package utils

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
)

// Locally modified terminal color
// https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/issues/180
func EntityNameColor(in string) string {
	return terminal.EntityNameColor(SaneEntityForColorTerm(in))
}

// Convert from %(!NOZERO) to % to remove the IBMCloud CLI bug
// https://github.com/IBM-Cloud/ibm-cloud-cli-sdk/issues/180
func SaneEntityForColorTerm(in string) string {
	return strings.Replace(in, "%", "%%", -1)
}
