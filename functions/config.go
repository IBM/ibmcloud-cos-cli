package functions

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/utils"
)

const (
	invalidType = "-INVALID-"
)

// ConfigOption holds a representation of the config values location and how to display them
type ConfigOption struct {
	// Key config key
	Key string
	// Display label to use for the key
	Display string
	// Default when config does not contains the value what to show
	Default string
	// some function to process the raw value from the config
	PostLoad func(interface{}) string
}

var (
	configOptionLastUpdated = ConfigOption{
		Key:     config.LastUpdated,
		Display: T("Last Updated"),
	}

	configOptionDefaultRegion = ConfigOption{
		Key: config.DefaultRegion,
		// it cannot be set here package load time, in app init current user region will override fallback
		//Default: config.FallbackRegion,
		Display: T("Default Region"),
	}

	configOptionDefaultDownloadLocation = ConfigOption{
		Key:     config.DownloadLocation,
		Default: config.FallbackDownloadLocation,
		Display: T("Download Location"),
	}

	configOptionDefaultCRN = ConfigOption{
		Key:     config.CRN,
		Display: T("CRN"),
	}

	configOptionHMACKey = ConfigOption{
		Key: config.AccessKeyID,
	}

	configOptionHMACSecret = ConfigOption{
		Key:      config.SecretAccessKey,
		PostLoad: mask,
	}

	configOptionAuthenticationMethod = ConfigOption{
		Key:      config.HMACProvided,
		Display:  T("Authentication Method"),
		PostLoad: mapAuthMethod,
		Default:  mapAuthMethod(config.HMACProvidedDefault),
	}

	configOptionURLStyle = ConfigOption{
		Key:      config.ForcePathStyle,
		Display:  T("URL Style"),
		PostLoad: mapURLStyle,
		Default:  mapURLStyle(config.ForcePathStyleDefault),
	}

	configOptionServiceEndpoint = ConfigOption{
		Key:     config.ServiceEndpointURL,
		Display: T("Service Endpoint"),
	}
)

// ConfigList Lists current config values
func ConfigList(c *cli.Context) error {
	configOptionDefaultRegion.Default = config.FallbackRegion

	// list as no args , check if the number of args is ZERO
	if c.NArg() != 0 {
		cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
	}

	// takes the CosContext from application metadata
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	ui := cosContext.UI
	conf := cosContext.Config

	// list of all configurations values to display in the list
	options := []ConfigOption{
		// last updated row
		configOptionLastUpdated,
		// default region row
		configOptionDefaultRegion,
		// default download location row
		configOptionDefaultDownloadLocation,
		// default crn row
		configOptionDefaultCRN,
		// hmac key row
		configOptionHMACKey,
		// hmac secret row
		configOptionHMACSecret,
		// authentication method row
		configOptionAuthenticationMethod,
		// url style
		configOptionURLStyle,
		// service endpoint
		configOptionServiceEndpoint,
	}

	// builds a table with the config values to display
	table := buildTable(ui, conf, options)
	// table table in screen
	table.Print()

	return nil
}

// this function masks secrets, so confidential details are not openly displayed in screen
func mask(input interface{}) string {
	// type switch over input type
	switch v := input.(type) {
	case string:
		// if string replace the content by *
		return strings.Repeat("*", len(v))
	default:
		// if not string not expected situation
		return invalidType
	}
}

// maps the authentication value between the the way it is stored and the way it is displayed
func mapAuthMethod(input interface{}) string {
	// type switch over input type
	switch v := input.(type) {
	case bool:
		// if bool map it to authentication method label
		return boolToAuth(v)
	default:
		// if not bool not expected situation
		return invalidType
	}
}

// convert the boolean HMACProvided to an authentication method name
func boolToAuth(b bool) string {
	// if true is hmac else is iam
	if b {
		return config.HMAC
	}
	return config.IAM
}

// convert an authentication method name to a boolean ( is HMAC )
func authToBool(auth string) (bool, error) {
	// switch over auth method string
	switch strings.ToUpper(auth) {
	case config.HMAC:
		// if hmac, map it to true
		return true, nil
	case config.IAM:
		// if iam, map it false
		return false, nil
	default:
		// if not previous cases rise an error
		return false, errors.New("invalid.method")
	}
}

// mapURLStyle maps the authentication value from stored ForcePathStyle to VHost or Path
func mapURLStyle(input interface{}) string {
	// type switch over input type
	switch v := input.(type) {
	case bool:
		// if bool map it to authentication method label
		return boolToURLStyle(v)
	default:
		// if not bool not expected situation
		return invalidType
	}
}

// boolToURLStyle maps the authentication value from boolean ForcePathStyle to VHost or Path
func boolToURLStyle(b bool) string {
	// if true is hamc else is iam
	if b {
		return config.Path
	}
	return config.VHost
}

// urlStyleToBool maps VHost or Path to ForcePathStyle
func urlStyleToBool(urlStyle string) (bool, error) {
	// switch over auth method string
	switch strings.ToUpper(urlStyle) {
	case strings.ToUpper(config.VHost):
		// if hamc map it to true
		return false, nil
	case strings.ToUpper(config.Path):
		// if iam  map it false
		return true, nil
	default:
		// if not previous cases rise an error
		return false, errors.New("invalid.urlstyle")
	}
}

// build an populate a table with the configuration values to display
func buildTable(ui terminal.UI, pc plugin.PluginConfig, fields []ConfigOption) terminal.Table {
	// sets table headers
	table := ui.Table([]string{T("Key"), T("Value")})
	// iterate across all values definition and add them as table rows
	for _, row := range fields {
		key, value := "", row.Default
		if row.Display != "" {
			key = row.Display
		} else {
			key = row.Key
		}
		if pc.Exists(row.Key) {
			rawValue := pc.Get(row.Key)
			if row.PostLoad != nil {
				value = row.PostLoad(rawValue)
			} else {
				value = fmt.Sprintf("%v", rawValue)
			}
		}
		table.Add(key, value)
	}
	return table
}

// ConfigChangeDefaultRegion allows the user to change the default region for the program to look for a bucket.
func ConfigChangeDefaultRegion(c *cli.Context) error {

	configOptionDefaultRegion.Default = config.FallbackRegion

	// takes the CosContext from application metadata
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	ui := cosContext.UI
	conf := cosContext.Config

	// validate the number of flags and the number of arguments passed to the command
	if c.NumFlags() > 1 || c.NArg() > 0 {
		// if number of flags or number of args do not match show usage message
		cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
	}

	// if list flag is set, display the value of the region and exit
	if c.IsSet(flags.List) {
		options := []ConfigOption{configOptionDefaultRegion}
		// build the table using rows/options
		table := buildTable(ui, conf, options)
		// display table
		table.Print()
		return nil
	}

	var region string
	var err error

	// if region falg is set , use its value
	if c.IsSet(flags.Region) {
		region = c.String(flags.Region)
	} else {
		// else tries to read from the console

		// when ask for a new value of region tries to display current as default
		if region, err = conf.GetStringWithDefault(config.DefaultRegion, config.FallbackRegion); err != nil {
			ui.Failed(T("Unable to load region."))
			return cli.NewExitError("", 1)
		}
		err = ui.Prompt(configOptionDefaultRegion.Display, &terminal.PromptOptions{}).Resolve(&region)
		if err != nil {
			ui.Failed(T("Unable to read new region."))
			return cli.NewExitError("", 1)
		}
	}

	// Set the config file with default region
	err = conf.Set(config.DefaultRegion, region)
	if err != nil {
		ui.Failed(T("Unable to save new region."))
		return cli.NewExitError("", 1)
	}

	// Set the config file with default region
	conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))

	// OK
	ui.Ok()

	// Output the message
	ui.Say(T("Successfully saved default region. The program will look for buckets in the region {{.region}}.",
		map[string]interface{}{"region": terminal.EntityNameColor(region)}))

	// Return
	return nil
}

// ConfigSetDLLocation allows the user to set their default download location (where files will be downloaded)
func ConfigSetDLLocation(c *cli.Context) error {
	// takes the CosContext from application metadata
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	ui := cosContext.UI
	conf := cosContext.Config

	// check the number of args and falgs
	if c.NumFlags() > 1 || c.NArg() > 0 {
		cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
	}

	// if list flag is set, display current download default location and exit
	if c.IsSet(flags.List) {

		options := []ConfigOption{
			// add default download location to table
			configOptionDefaultDownloadLocation,
		}
		// build table
		table := buildTable(ui, conf, options)
		// display table
		table.Print()
		return nil
	}

	var downloadLocation string
	var err error

	// creates a function with current cosContext that will be used to validate the location
	var validated = false
	validateFunc := buildValidateDownloadLocation(cosContext)

	// if ddl flag provided use its value as new default download location
	if c.IsSet(flags.DDL) {
		downloadLocation = c.String(flags.DDL)
	} else {
		// else prompt the user for a new download location
		if downloadLocation, err = conf.GetStringWithDefault(config.DownloadLocation,
			config.FallbackDownloadLocation); err != nil {
			ui.Failed(T("Unable to load default download location."))
			return cli.NewExitError("", 1)
		}
		err = ui.Prompt("Default Download Location", &terminal.PromptOptions{ValidateFunc: validateFunc}).
			Resolve(&downloadLocation)
		if err != nil {
			ui.Failed(T("Unable to read new default download location."))
			return cli.NewExitError("", 1)
		}
		validated = true
	}

	// prompt operation validates the new value of location before accept it
	// but if the value comes from the ddl flag needs to be validated before accepted
	if !validated {
		err = validateFunc(downloadLocation)
		if err != nil {
			ui.Failed(err.Error())
			return cli.NewExitError("", 1)
		}
	}

	// Saving default download location in the config
	ui.Say(T("Saving default download location..."))

	// Set the config file with the download location
	err = conf.Set(config.DownloadLocation, downloadLocation)
	if err != nil {
		ui.Failed(T("Unable to store new download location."))
		return cli.NewExitError("", 1)
	}

	// Set the last update timestamp in the config
	conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))

	// Return OK
	ui.Ok()

	// Output the successful message
	ui.Say(T("Successfully saved download location. New files will be downloaded to '") +
		terminal.EntityNameColor(downloadLocation) + "'.")

	// Return
	return nil
}

// creates a new validation function
func buildValidateDownloadLocation(cosContext *utils.CosContext) func(string) error {
	// return a function that checks if it's a valid download Filepath
	return func(location string) error {
		// retrieve location info
		info, err := cosContext.GetFileInfo(location)

		// check that location exists
		if err != nil && os.IsNotExist(err) {
			// if not exists return error
			return errors.New(T("The specified download location does not exist on your system. Try again with a valid download path."))
		}

		// check if location is a dir
		if err == nil && !info.IsDir() {
			// if not a dir return error
			return errors.New(T("The specified download location is not valid. Try again with a valid download path."))
		}

		return nil
	}
}

// ConfigCRN allows the user to store the CRN into Credentials.json. This is the first setup that any user does,
// so it's a bit more involved
func ConfigCRN(c *cli.Context) error {
	// takes the CosContext from application metadata
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	ui := cosContext.UI
	conf := cosContext.Config

	// validate the number of flags and the number of arguments passed to the command
	if (c.IsSet(flags.List) && c.NumFlags() > 1) || c.NumFlags() > 2 || c.NArg() > 0 {
		cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
	}

	// if list flag set
	// list and exit
	if c.IsSet(flags.List) {
		options := []ConfigOption{
			// add crn row to config table
			configOptionDefaultCRN,
		}
		// build table
		table := buildTable(ui, conf, options)
		// print table on screen
		table.Print()
		return nil
	}

	var crn string
	var err error
	force := false
	if c.IsSet(flags.Force) {
		force = true
	}

	var oldCRN string

	// check if a CRN was set before
	// if yes, warns about current being replaced in the process
	if conf.Exists(config.CRN) {

		oldCRN, err = conf.GetString(config.CRN)
		var lastUpdatedDate string
		lastUpdatedDate, err = conf.GetStringWithDefault(config.LastUpdated, "UNKNOWN")

		ui.Warn(T("WARNING: You have already stored a service instance ID before."))
		ui.Say(T("It was last updated on {{.lastupdated}}.", map[string]interface{}{"lastupdated": lastUpdatedDate}))

	}

	// if crn flag is set uses its value as new crn
	if c.IsSet(flags.CRN) {
		crn = c.String(flags.CRN)
		if !force && oldCRN != "" && crn != oldCRN {
			ui.ChoicesPrompt("Select the CRN", []string{crn, oldCRN},
				&terminal.PromptOptions{}).Resolve(&crn)
		}
	} else {
		// else prompt the user for a new one using the old old one as fallback
		crn = oldCRN
		if err = ui.Prompt("Resource Instance ID CRN: ", nil).Resolve(&crn); err != nil {
			ui.Failed(T("Unable to get new CRN."))
			return cli.NewExitError("", 1)
		}
	}

	// Alerts users that CRN is to be saved in the config file
	ui.Say(T("Saving new Service Instance ID..."))

	// Set the CRN in the confg file
	err = conf.Set(config.CRN, crn)
	if err != nil {
		ui.Failed(T("Unable to store Secret key."))
		return cli.NewExitError("", 1)
	}

	// Set the last update timestamp in the config
	conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))

	// Return OK
	ui.Ok()

	// Output the message
	ui.Say(T("Successfully stored your service instance ID."))

	// Return
	return nil
}

// ConfigAmazonHMAC stores HMAC Credentials in the config file
func ConfigAmazonHMAC(c *cli.Context) error {
	// takes the CosContext from application metadata
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	ui := cosContext.UI
	conf := cosContext.Config

	// validate the number of flags and arguments passed to the command
	if c.NumFlags() > 1 || c.NArg() > 0 {
		cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
	}

	// if list flag set
	// display current values and exit
	if c.IsSet(flags.List) {

		options := []ConfigOption{
			// add hmac key row
			configOptionHMACKey,
			// add hmac secret row
			configOptionHMACSecret,
		}
		// build table
		table := buildTable(ui, conf, options)
		// print table
		table.Print()
		return nil
	}

	// Local initialization of HMAC credentials
	var accessKeyID string
	var secretAccessKey string

	// Ask the user for the secret and access keys, using the IBM Cloud CLI SDK's Prompt functionality.
	ui.Prompt("Access key", nil).Resolve(&accessKeyID)

	// Saves the HMAC access key id in the config file
	err := conf.Set(config.AccessKeyID, accessKeyID)
	if err != nil {
		ui.Failed(T("Unable to store Access key."))
		return cli.NewExitError("", 1)
	}

	// Prompt users for the secret key
	ui.Prompt("Secret key", nil).Resolve(&secretAccessKey)

	// Saves the HMAC secret access key in the config file
	err = conf.Set(config.SecretAccessKey, secretAccessKey)
	if err != nil {
		ui.Failed(T("Unable to store Secret key."))
		return cli.NewExitError("", 1)
	}

	// Set the last update timestamp in the config
	conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))

	// Print OK
	ui.Ok()

	// Output the message
	ui.Say(T("Successfully saved HMAC Credentials to file."))

	// Return
	return nil
}

// ConfigSetAuthMethod allows the user to switch between IAM and HMAC based authentication, by setting "HMACProvided"
// to either true or false.
func ConfigSetAuthMethod(c *cli.Context) error {
	// takes the CosContext from application metadata
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	ui := cosContext.UI
	conf := cosContext.Config

	// validates the number of arguments and flags passed to to command
	if c.NumFlags() > 1 || c.NArg() > 0 {
		cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
	}

	// if list flag is set, display current value and exit
	if c.IsSet(flags.List) {
		options := []ConfigOption{
			// add auth method row to config table
			configOptionAuthenticationMethod,
		}
		// build table
		table := buildTable(ui, conf, options)
		// print table
		table.Print()
		return nil
	}

	var authMethod string
	var authBool bool
	var err error

	// if method flag is set use it
	if c.IsSet(flags.Method) {
		authMethod = c.String(flags.Method)
	} else {
		// else prompt for new value, using previous as default
		if authBool, err = conf.GetBoolWithDefault(config.HMACProvided, config.HMACProvidedDefault); err != nil {
			ui.Failed(T("Unable to load config method."))
			return cli.NewExitError("", 1)
		}
		authMethod = boolToAuth(authBool)

		err = ui.ChoicesPrompt("Select the Authentication Method", []string{config.IAM, config.HMAC},
			&terminal.PromptOptions{}).Resolve(&authMethod)

		if err != nil {
			ui.Failed(T("Unable to read new Authentication Method."))
			return cli.NewExitError("", 1)
		}
	}

	// maps user input to the value to be stored
	authBool, err = authToBool(authMethod)
	if err != nil {
		ui.Failed(T("Unable to parse authentication method."))
		return cli.NewExitError("", 1)
	}

	// Set HMAC provide in the config file
	err = conf.Set(config.HMACProvided, authBool)
	if err != nil {
		ui.Failed(T("Unable to switch authentication method."))
		return cli.NewExitError("", 1)
	}

	// Set the last update timestamp in the config
	conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))

	// Print OK
	ui.Ok()

	// Output the message
	ui.Say(T("Successfully switched to {{.Auth}}-based authentication. The program will access your Cloud Object Storage account using your {{.Auth}} Credentials.",
		map[string]interface{}{"Auth": terminal.EntityNameColor(boolToAuth(authBool))}))

	// Return
	return nil

}

// ConfigSetRegionsEndpointURL override the value of the default regions endpoint
func ConfigSetRegionsEndpointURL(c *cli.Context) error {
	// takes the CosContext from application metadata
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	ui := cosContext.UI
	conf := cosContext.Config

	// validates the number of args and flags
	if c.NumFlags() != 1 || c.NArg() > 0 {
		cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
	}
	// if list flag is set,
	// list current value and exit
	if c.IsSet(flags.List) {
		options := []ConfigOption{
			{
				Key: config.RegionsEndpointURL,
			},
		}
		table := buildTable(ui, conf, options)
		table.Print()
		return nil
	}

	var regionsURL string
	var err error

	// if url flag is set, use its value as new regions endpoint value
	if c.IsSet(flags.URL) {
		regionsURL = c.String(flags.URL)

		err = conf.Set(config.RegionsEndpointURL, regionsURL)
		if err != nil {
			ui.Failed(T("Unable to save new regions endpoint URL."))
			return cli.NewExitError("", 1)
		}

		conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))
		ui.Ok()
		ui.Say(T("Successfully updated regions endpoint URL to {{.URL}}.",
			map[string]interface{}{"URL": terminal.EntityNameColor(regionsURL)}))
	}

	// if clear falg set, clear the override value, falling back to original default value
	if c.IsSet(flags.Clear) {

		err = conf.Erase(config.RegionsEndpointURL)
		if err != nil {
			ui.Failed(T("Unable to clear regions endpoint URL."))
			return cli.NewExitError("", 1)
		}

		// Set the last update timestamp in the config
		conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))

		// OK
		ui.Ok()

		// Output the successfully message
		ui.Say(T("Successfully cleared regions endpoint URL."))
	}

	// Return
	return nil
}

// ConfigSetURLStyle allows the user to switch between VHost and Path URL Styles
func ConfigSetURLStyle(c *cli.Context) error {
	// takes the CosContext from application metadata
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	ui := cosContext.UI
	conf := cosContext.Config

	// validates the number of arguments and flags passed to to command
	if c.NumFlags() > 1 || c.NArg() > 0 {
		cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
	}

	// if list flag is set, display current value and exit
	if c.IsSet(flags.List) {
		options := []ConfigOption{
			// add auth method row to config table
			configOptionURLStyle,
		}
		// build table
		table := buildTable(ui, conf, options)
		// print table
		table.Print()
		return nil
	}

	var urlStyle string
	var forcePathStyle bool
	var err error

	// if method flag is set use it
	if c.IsSet(flags.Style) {
		urlStyle = c.String(flags.Style)
	} else {
		// else prompt for new value, using previous as default
		if forcePathStyle, err = conf.GetBoolWithDefault(config.ForcePathStyle, config.ForcePathStyleDefault); err != nil {
			ui.Failed(T("Unable to load current url style."))
			return cli.NewExitError("", 1)
		}
		urlStyle = boolToURLStyle(forcePathStyle)

		err = ui.ChoicesPrompt("Select the URL Style", []string{config.VHost, config.Path},
			&terminal.PromptOptions{}).Resolve(&urlStyle)

		if err != nil {
			ui.Failed(T("Unable to read new URL Style value."))
			return cli.NewExitError("", 1)
		}
	}

	// maps user input to the value to be stored
	forcePathStyle, err = urlStyleToBool(urlStyle)
	if err != nil {
		ui.Failed(T("Unable to parse URL style."))
		return cli.NewExitError("", 1)
	}

	// Set ForcePathStyle in the config file
	err = conf.Set(config.ForcePathStyle, forcePathStyle)
	if err != nil {
		ui.Failed(T("Unable to switch URL style."))
		return cli.NewExitError("", 1)
	}

	// Set the last update timestamp in the config
	conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))

	// Print OK
	ui.Ok()

	// Output the message
	ui.Say(T("Successfully saved S3 {{.Style}} style hosting option. CLI calls will use {{.Style}} style requests.",
		map[string]interface{}{"Style": terminal.EntityNameColor(boolToURLStyle(forcePathStyle))}))

	// Return
	return nil
}

// ConfigSetEndpointURL Sets the Service Endpoint and related actions
func ConfigSetEndpointURL(c *cli.Context) error {
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	ui := cosContext.UI
	conf := cosContext.Config

	if c.NumFlags() > 1 || c.NArg() > 0 {
		cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
	}

	if c.IsSet(flags.List) {
		options := []ConfigOption{
			{
				Key: config.ServiceEndpointURL,
			},
		}
		table := buildTable(ui, conf, options)
		table.Print()
		return nil
	}

	var serviceEndpoint string
	var err error

	if c.IsSet(flags.URL) {
		serviceEndpoint = c.String(flags.URL)

		err = conf.Set(config.ServiceEndpointURL, serviceEndpoint)
		if err != nil {
			ui.Failed(T("Unable to save new Service Endpoint URL."))
			return cli.NewExitError("", 1)
		}

		conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))

		ui.Ok()
		ui.Say(T("Successfully updated service endpoint URL.",
			map[string]interface{}{"URL": terminal.EntityNameColor(serviceEndpoint)}))
	}

	if c.IsSet(flags.Clear) {

		err = conf.Erase(config.ServiceEndpointURL)
		if err != nil {
			ui.Failed(T("Unable to clear service endpoint URL."))
			return cli.NewExitError("", 1)
		}

		conf.Set(config.LastUpdated, time.Now().Local().Format(config.StandardTimeFormat))

		ui.Ok()

		ui.Say(T("Successfully cleared service endpoint URL"))
	}

	return nil
}
