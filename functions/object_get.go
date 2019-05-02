package functions

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/urfave/cli"
)

// Constant message
const (
	overRideDefaultLocMessage = "Override the default location providing an OUTFILE parameter"
)

// ObjectGet downloads a file from a bucket.
// Parameter:
//   	CLI Context Application
// Returns:
//  	Error = zero or non-zero
func ObjectGet(c *cli.Context) error {

	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Monitor the file
	keepFile := true

	var dstPath string
	// Remove files when open later in the process
	defer func() {
		if !keepFile && dstPath != "" {
			cosContext.Remove(dstPath)
		}
	}()

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Set GetObjectInput
	input := new(s3.GetObjectInput)

	// Required parameter for GetObject
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	// No optional parameter for GetObject
	options := map[string]string{
		fields.IfMatch:                    flags.IfMatch,
		fields.IfModifiedSince:            flags.IfModifiedSince,
		fields.IfNoneMatch:                flags.IfNoneMatch,
		fields.IfUnmodifiedSince:          flags.IfUnmodifiedSince,
		fields.Range:                      flags.Range,
		fields.ResponseCacheControl:       flags.ResponseCacheControl,
		fields.ResponseContentDisposition: flags.ResponseContentDisposition,
		fields.ResponseContentEncoding:    flags.ResponseContentEncoding,
		fields.ResponseContentLanguage:    flags.ResponseContentLanguage,
		fields.ResponseContentType:        flags.ResponseContentType,
		fields.ResponseExpires:            flags.ResponseExpires,
	}

	// Validate User Inputs and Retrieve Region
	region, err := ValidateUserInputsAndSetRegion(c, input, mandatory, options, conf)
	if err != nil {
		ui.Failed(err.Error())
		cli.ShowCommandHelp(c, c.Command.Name)
		//return cli.NewExitError(err.Error(), 1)
	}

	var outFile bool
	outFile, dstPath, err = getDwnLdPath(c, input)
	if err != nil || dstPath == "" {
		return err
	}

	// Creates an empty file placeholder
	file, err := cosContext.WriteCloserOpen(dstPath)
	if err != nil {
		ui.Failed(badFilePathOpen(dstPath, outFile))
		return cli.NewExitError("", 1)
	}

	// Delays closing the file until download is complete
	defer file.Close()
	keepFile = false

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Downloading object..."))

	// Download the object by calling GetObject
	resp, err := client.GetObject(input)
	// Error handling
	if err != nil {
		if strings.Contains(err.Error(), "EmptyStaticCreds") {
			ui.Failed(err.Error() + "\n" + T("Try logging in using 'ibmcloud login'."))
		} else {
			ui.Failed(err.Error())
		}
		os.Remove(dstPath)
		return cli.NewExitError("", 1)
	}
	// Retrieve the body content
	defer resp.Body.Close()
	io.Copy(file, resp.Body)
	keepFile = true

	// Success
	ui.Ok()
	ui.Say(T("Successfully downloaded '{{.Key}}' from bucket '{{.Bucket}}'",
		map[string]interface{}{"Key": utils.EntityNameColor(*input.Key),
			"Bucket": utils.EntityNameColor(*input.Bucket)}))

	ui.Say(FormatFileSize(*resp.ContentLength) + T(" downloaded."))

	// Return
	return nil
}

func getDwnLdPath(c *cli.Context, input *s3.GetObjectInput) (bool, string, error) {

	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	outFile := false
	var dstPath string
	endsInSeparator := false
	var dirPath string

	// Checks if outfile is present
	if c.NArg() == 1 {
		// Break down the outfile input
		outFile = true
		dstPath = c.Args()[0]
		endsInSeparator = dstPath[len(dstPath)-1] == filepath.Separator
		dstPath = filepath.Clean(dstPath)
		dirPath = filepath.Dir(dstPath)

	} else {
		// Outfile is not present, check Default Download Location in config file
		outFile = false

		// Assign Default Download location to a string
		downloadsPath, err := conf.GetStringWithDefault(config.DownloadLocation, config.FallbackDownloadLocation)
		if err != nil || downloadsPath == "" {
			ui.Failed(T("Unable to get Default Download Location."))
			return false, "", cli.NewExitError("", 1)
		}

		// Get the path to the Default Download Location
		if pathInfo, err := cosContext.GetFileInfo(downloadsPath); err != nil || !pathInfo.IsDir() {
			str := T("The download directory '{{.dl}}' is invalid",
				map[string]interface{}{"dl": downloadsPath})
			str += "\n" + T("Set a valid download location using 'cos config --ddl'")
			ui.Failed(str)
			return false, "", cli.NewExitError("", 1)
		}

		// Check if the filepath has a separator
		endsInSeparator = (*((*input).Key))[len(*((*input).Key))-1] == filepath.Separator

		// Store the destination path and directory path for download
		dstPath = filepath.Join(downloadsPath, *input.Key)
		dirPath = filepath.Dir(dstPath)

	}
	// The file path has a separator, exit
	if endsInSeparator {
		ui.Failed(badFilePathIsDir(dstPath+string([]rune{filepath.Separator}), outFile))
		return false, "", cli.NewExitError("", 1)
	}
	// Set the path info
	if pathInfo, err := cosContext.GetFileInfo(dirPath); err != nil || !pathInfo.IsDir() {
		ui.Failed(badDirPath(dirPath, outFile))
		return false, "", cli.NewExitError("", 1)
	}
	// Checks to see if a file with the same name as the file to be downloaded exists already or not. Warns the user
	// to prevent accidental overwriting.
	if pathInfo, err := cosContext.GetFileInfo(dstPath); err == nil {
		// check if the path is itself a directory
		if pathInfo.IsDir() {
			ui.Failed(badFilePathIsDir(dstPath, outFile))
			return false, "", cli.NewExitError("", 1)
		}
		// Alert user whether they want to overwrite existing file
		resolve := false
		ui.Warn(T("WARNING: An object with the name '{{.file}}' already exists at '{{.dl}}'.",
			map[string]interface{}{"file": filepath.Base(dstPath), "dl": dirPath}))
		ui.Prompt(T("Are you sure you would like to overwrite it?"),
			new(terminal.PromptOptions)).Resolve(&resolve)

		// Cancel the operation if user denies overwriting or exits the prompt
		if !resolve {
			ui.Say(T("Operation canceled."))
			return false, "", nil //cli.NewExitError("", 0)
		}
	}
	return outFile, dstPath, nil
}

// badDirPath tells user that the download directory does not exist
// Parameter:
// 		destination directory path (string)
//		outfile set as parameter (boolean)
// Return:
//		invalid download directory (string)
func badDirPath(dstDirPath string, outFileSet bool) string {
	str := T("The download directory '{{.dl}}' is invalid.", map[string]interface{}{"dl": dstDirPath})
	if !outFileSet {
		str += "\n" + T(overRideDefaultLocMessage)
	}
	return str
}

// badFilePathIsDir tells users that the download destination itself is a directory
// Parameter:
// 		destination file path (string)
//		outfile set as parameter (boolean)
// Return:
//		download directory (string)
func badFilePathIsDir(dstFilePath string, outFileSet bool) string {
	str := T("The download destination '{{.dl}}' is a directory.",
		map[string]interface{}{"dl": dstFilePath})
	if !outFileSet {
		str += "\n" + T(overRideDefaultLocMessage)
	}
	return str
}

// badFilePathOpen prompts user with an error opening the destination file
// Parameter:
// 		destination file path (string)
//		outfile set as parameter (boolean)
// Return:
//		fail opening the file to write (string)
func badFilePathOpen(dstFilePath string, outFileSet bool) string {
	str := T("Error opening '{{.dl}}' to write ", map[string]interface{}{"dl": dstFilePath})
	if !outFileSet {
		str += "\n" + T(overRideDefaultLocMessage)
	}
	return str
}
