package functions

import (
	"io"
	"path/filepath"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/render"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
)

// ObjectGet downloads a file from a bucket.
// Parameter:
//
//	CLI Context Application
//
// Returns:
//
//	Error = zero or non-zero
func ObjectGet(c *cli.Context) (err error) {
	// check the number of arguments
	if c.NArg() > 1 {
		// if the number of arguments is bigger than 1 throw an error
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		// returns stopping any further processing
		return
	}

	// allocate a variable fot the COS Context
	var cosContext *utils.CosContext
	// do a get and check of the COS Context
	if cosContext, err = GetCosContext(c); err != nil {
		// if the get cos context returned error,
		// return the error and stop any further processing
		return
	}

	// Monitor the file
	keepFile := false

	// Download location
	var dstPath string

	// register a deferred function to remove the file if the get did not complete
	defer func() {
		// check that the file download did not complete
		// and that the file name is not empty
		if !keepFile && dstPath != "" {
			// call the context wrapper to the SO remove function
			cosContext.Remove(dstPath)
		}
	}()

	// initialize the input with a pointer to a empty GetObjectInput
	input := new(s3.GetObjectInput)

	// Define the mandatory fields of GetObjectInput,
	// and the flags to be used as source for the input values
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	// Define the optional fields of GetObjectInput,
	// and the flags to be used as source for the input values
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
		fields.VersionId:                  flags.VersionId,
	}

	// populate the input values using the mandatory and optional maps define before
	// and checks that there was no error during the map operation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		// if a error occurred during the mapping,
		// propagate the error and stop any further processing
		return
	}

	// allocate a variable of type io.WriteCloser,
	// that will be used as destination of the download/get operation
	var file io.WriteCloser
	// calls the auxiliary getAndValidateDownloadPath passing
	// the Cos Context, the first argument of command invocation and the key
	// checks if an error occurred or
	// the file value is empty ( when destination exists and user do not confirm overwrite )
	if dstPath, file, err = getAndValidateDownloadPath(cosContext, c.Args().First(),
		aws.StringValue(input.Key), c.IsSet(flags.Force)); err != nil || file == nil {
		// propagate current error value ( can be nil )
		// and stops any further processing
		return
	}
	// register a deferred function to close the writer on the exit of current scope
	defer file.Close()

	// allocate a variable to hold the S3API
	var client s3iface.S3API
	// fetch a client from COS Context overriding the default region if needed
	// also checks for error in the fetch operation
	if client, err = cosContext.GetClient(c.String(flags.Region)); err != nil {
		// if an error occurred,
		// propagate the error and stop any further processing
		return
	}

	// allocate a variable to hold the GetObject result
	var output *s3.GetObjectOutput
	// invokes the GetObject operation using the client and checks the result
	if output, err = client.GetObject(input); err != nil {
		// if an error occurs in the GetObject request call
		// propagate the error and stop any further processing
		return
	}

	// register a function to close the Body once function gets out of scope
	defer output.Body.Close()

	// No need to show progress bar when user wants json output,
	display := cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON))
	if _, ok := display.(*render.TextRender); ok {
		// Create a new progress bar with the total content length in bytes
		bar := pb.New64(*output.ContentLength).SetUnits(pb.U_BYTES)
		bar.Start()
		// Write output to both file and bar
		if _, err = io.Copy(io.MultiWriter(file, bar), output.Body); err != nil {
			return
		}
		bar.Finish()
	} else {
		// copy from the response body to the file defined in the command invocation / default destination
		// and checks no error occurred
		if _, err = io.Copy(file, output.Body); err != nil {
			// if error occurred propagate the error
			return
		}
	}

	// flags the file should not be deleted
	keepFile = true

	// use a wrapper struct to add a new field
	objectOutputWrapper := &render.GetObjectOutputWrapper{
		GetObjectOutput:  output,
		DownloadLocation: &dstPath,
	}

	// render the result in JSON or Textual format
	// depending if the flag JSON was passed in the command invocation
	err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, objectOutputWrapper, nil)

	// Return
	return

}

// getAndValidateDownloadPath is an auxiliary function that
// takes the COS Context, the value of the first argument as destination location and the value of the key
// checks if the destination location is not empty
// if not set and/or is empty
// uses the concatenation of the default download location and the key as the destination location
// checks if the destination location is writable and not already exists,
// if exists ask confirmation before overwrite it
func getAndValidateDownloadPath(
	cosContext *utils.CosContext,
	outfile string,
	key string,
	force bool,
) (
	path string,
	wc utils.WriteCloser,
	err error,
) {
	// identifies if the rule used to get the destination location is the default rule or is an invocation parameter
	usingDefaultRule := false
	if outfile == "" {
		var downloadLocation string
		if downloadLocation, err = cosContext.GetDownloadLocation(); err != nil {
			return
		}
		// path join remove trailing /
		// to prevent losing it, it will be saved before join
		trailing := key[len(key)-1:]
		key = key[:len(key)-1]
		if len(key) > 0 {
			outfile = filepath.Join(downloadLocation, key)
		} else {
			outfile = downloadLocation + string(filepath.Separator)
		}
		outfile += trailing
		usingDefaultRule = true
	}
	// verify if the destination location is a folder by checking the last character
	if filepath.Separator == outfile[len(outfile)-1] {
		err = &errors.ObjectGetError{
			Location:         outfile,
			UsingDefaultRule: usingDefaultRule,
			Cause:            errors.IsDir,
		}
		return
	}

	// checks that the destination location does not already exists
	if fileInfo, innerError := cosContext.GetFileInfo(outfile); innerError == nil {
		// if it exists checks it is not a folder
		if fileInfo.IsDir() {
			err = &errors.ObjectGetError{
				Location:         outfile,
				UsingDefaultRule: usingDefaultRule,
				Cause:            errors.IsDir,
			}
			return
		}
		confirmed := false
		dir := filepath.Dir(outfile)
		if absDir, err := filepath.Abs(dir); err == nil {
			dir = absDir
		}

		if !force {
			// warn that the destination location already exists and ask for confirmation before overwrite
			cosContext.UI.Warn(render.WarningGetObject(filepath.Base(outfile), dir))
			cosContext.UI.Prompt(render.MessageConfirmationContinue(), &terminal.PromptOptions{}).Resolve(&confirmed)
			// if the user does not confirm the overwrite, display a message and exit this function
			if !confirmed {
				cosContext.UI.Say(render.MessageOperationCanceled())
				return
			}
		}
	}

	// open the destination location for write and checks for errors
	if wc, err = cosContext.WriteCloserOpen(outfile); err != nil {
		innerError := err
		err = &errors.ObjectGetError{
			ParentError:      innerError,
			Location:         outfile,
			UsingDefaultRule: usingDefaultRule,
			Cause:            errors.Opening,
		}
		return
	}
	path = outfile
	return
}
