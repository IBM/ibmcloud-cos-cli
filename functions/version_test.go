//go:build unit
// +build unit

package functions_test

import (
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestVersionSunnyPath(t *testing.T) {
	// Reset mocks after the test
	defer providers.MocksRESET()

	// --- Arrange ---
	// Disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		// Capture the exit code
		exitCode = &ec
	}

	// --- Act ----
	// Set OS args for the test
	os.Args = []string{"-", commands.Version}

	// Call the plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// Capture all output and errors
	output := providers.FakeUI.Outputs()
	errors := providers.FakeUI.Errors()

	// Assert that the exit code is zero
	assert.Equal(t, (*int)(nil), exitCode)
	// Assert that the output contains the expected string
	assert.Contains(t, output, "ibmcloud cos cli plugin")
	// Assert that there are no errors
	assert.NotContains(t, errors, "FAIL")
}
