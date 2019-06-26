//+build unit

package functions_test

import (
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

func TestURLStyleSetVHostByFlag(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	providers.
		MockPluginConfig.
		On("Set", config.ForcePathStyle, false).
		Return(nil)

	providers.
		MockPluginConfig.
		On("Set", config.LastUpdated, mock.AnythingOfType("string")).
		Return(nil)

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.Config, commands.URLStyle,
		"--" + flags.Style, config.VHost,
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	assert.Equal(t, (*int)(nil), exitCode)
	providers.MockPluginConfig.AssertExpectations(t)
}

func TestURLStyleSetPathByFlag(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	providers.
		MockPluginConfig.
		On("Set", config.ForcePathStyle, true).
		Return(nil)

	providers.
		MockPluginConfig.
		On("Set", config.LastUpdated, mock.AnythingOfType("string")).
		Return(nil)

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.Config, commands.URLStyle,
		"--" + flags.Style, config.Path,
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	assert.Equal(t, (*int)(nil), exitCode)
	providers.MockPluginConfig.AssertExpectations(t)
}

func TestURLStyleSetVHostByPrompt(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	providers.
		MockPluginConfig.
		On("GetBoolWithDefault", config.ForcePathStyle, config.ForcePathStyleDefault).
		Return(config.ForcePathStyleDefault, nil)

	providers.
		MockPluginConfig.
		On("Set", config.ForcePathStyle, false).
		Return(nil)

	providers.
		MockPluginConfig.
		On("Set", config.LastUpdated, mock.AnythingOfType("string")).
		Return(nil)

	providers.FakeUI.Inputs("1")

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.Config, commands.URLStyle}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	assert.Equal(t, (*int)(nil), exitCode)
	providers.MockPluginConfig.AssertExpectations(t)

	expectedPrompt := `Select the URL Style
1. VHost
2. Path
Enter a number (1)`
	assert.Contains(t, providers.FakeUI.Outputs(), expectedPrompt)
}

func TestURLStyleSetPathByPrompt(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	providers.
		MockPluginConfig.
		On("GetBoolWithDefault", config.ForcePathStyle, config.ForcePathStyleDefault).
		Return(config.ForcePathStyleDefault, nil)

	providers.
		MockPluginConfig.
		On("Set", config.ForcePathStyle, true).
		Return(nil)

	providers.
		MockPluginConfig.
		On("Set", config.LastUpdated, mock.AnythingOfType("string")).
		Return(nil)

	providers.FakeUI.Inputs("2")

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.Config, commands.URLStyle}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	assert.Equal(t, (*int)(nil), exitCode)
	providers.MockPluginConfig.AssertExpectations(t)

	expectedPrompt := `Select the URL Style
1. VHost
2. Path
Enter a number (1)`
	assert.Contains(t, providers.FakeUI.Outputs(), expectedPrompt)
}

func TestURLStyleListBeforeSet(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	providers.
		MockPluginConfig.
		On("Exists", config.ForcePathStyle).
		Return(false)

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.Config, commands.URLStyle,
		"--" + flags.List,
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	assert.Equal(t, (*int)(nil), exitCode)
	providers.MockPluginConfig.AssertExpectations(t)

	expectedOutput := `Key         Value   
URL Style   ` + config.VHost
	assert.Contains(t, providers.FakeUI.Outputs(), expectedOutput)
}

func TestURLStyleListAferSetVHost(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	providers.
		MockPluginConfig.
		On("Exists", config.ForcePathStyle).
		Return(true)

	providers.
		MockPluginConfig.
		On("Get", config.ForcePathStyle).
		Return(false)

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.Config, commands.URLStyle,
		"--" + flags.List,
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	assert.Equal(t, (*int)(nil), exitCode)
	providers.MockPluginConfig.AssertExpectations(t)

	expectedOutput := `Key         Value   
URL Style   ` + config.VHost
	assert.Contains(t, providers.FakeUI.Outputs(), expectedOutput)
}

func TestURLStyleListAferSetPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	providers.
		MockPluginConfig.
		On("Exists", config.ForcePathStyle).
		Return(true)

	providers.
		MockPluginConfig.
		On("Get", config.ForcePathStyle).
		Return(true)

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.Config, commands.URLStyle,
		"--" + flags.List,
	}

	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	assert.Equal(t, (*int)(nil), exitCode)
	providers.MockPluginConfig.AssertExpectations(t)

	expectedOutput := `Key         Value   
URL Style   ` + config.Path
	assert.Contains(t, providers.FakeUI.Outputs(), expectedOutput)
}
