package utils

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/render"
)

// CosContext Struct
// Cloud Object Storage Context struct includes
// Terminal User Interface object
// Plugin Config
// AWS S3 Session
// S3 Interface API Client
// List of available Regions for COS
// File Operations
type CosContext struct {
	UI               terminal.UI
	Config           plugin.PluginConfig
	Session          *session.Session
	ListKnownRegions ListKnownRegions
	JSONRender       *render.JSONRender
	TextRender       *render.TextRender
	ErrorRender      *render.ErrorRender
	ClientGen        func(*session.Session) s3iface.S3API
	DownloaderGen    func(svc s3iface.S3API, options ...func(*s3manager.Downloader)) Downloader
	UploaderGen      func(svc s3iface.S3API, options ...func(output *s3manager.Uploader)) Uploader

	FileOperations
}

// GetDownloadLocation loads the default download location from the plugin configuration
func (c *CosContext) GetDownloadLocation() (string, error) {
	return c.Config.GetStringWithDefault(config.DownloadLocation, config.FallbackDownloadLocation)
}

// GetCurrentRegion if the overrideRegion is empty
// loads the default region from the plugin configuration
// if not empty just return it
func (c *CosContext) GetCurrentRegion(overrideRegion string) (string, error) {
	if overrideRegion == "" {
		tmp, err := c.Config.GetStringWithDefault(config.DefaultRegion, config.FallbackRegion)
		return strings.TrimSpace(tmp), err
	}
	return overrideRegion, nil
}

// GetClient generates an S3 client to make requests through Go SDK
// - if an override region is passed, set the session region to be the passed value
// - if empty string is passed, load the default region from configuration
func (c *CosContext) GetClient(overrideRegion string) (client s3iface.S3API, err error) {
	var region string
	if region, err = c.GetCurrentRegion(overrideRegion); err != nil {
		return
	}
	serviceEndpoint, err := c.GetServiceEndpoint()
	if err != nil {
		return
	}
	splitURL := strings.Split(serviceEndpoint, ".")
	var cloudIndex int
	for i, n := range splitURL {
		if n == "cloud" {
			cloudIndex = i
			break
		}
	}
	if cloudIndex != 0 {
		region = splitURL[cloudIndex-1]
	}
	cfg := new(aws.Config).WithRegion(region).WithEndpoint(serviceEndpoint)
	sess := c.Session.Copy(cfg)
	client = c.ClientGen(sess)
	return
}

// GetDownloader gets a Downloader for the region passed as argument,
// if no region passed it uses the configuration default
func (c *CosContext) GetDownloader(overrideRegion string,
	options ...func(output *s3manager.Downloader)) (downloader Downloader, err error) {
	var client s3iface.S3API
	if client, err = c.GetClient(overrideRegion); err != nil {
		return
	}
	downloader = c.DownloaderGen(client, options...)
	return
}

// GetUploader gets a Uploader for the region passed as argument,
// if no region passed it uses the configuration default
func (c *CosContext) GetUploader(overrideRegion string,
	options ...func(output *s3manager.Uploader)) (uploader Uploader, err error) {
	var client s3iface.S3API
	if client, err = c.GetClient(overrideRegion); err != nil {
		return
	}
	uploader = c.UploaderGen(client, options...)
	return
}

// GetDisplay generates output either in text or json format
func (c *CosContext) GetDisplay(output string, isJSON bool) render.Display {

	// parse output string
	if output != "" {
		if strings.EqualFold(output, "json") {
			return c.JSONRender
		} else if strings.EqualFold(output, "text") {
			return c.TextRender
		} else {
			fmt.Println("error")
			// force error here
		}
	}

	// If output is not used, check for the boolean value
	if isJSON || output == "json" {
		return c.JSONRender
	}
	return c.TextRender
}

// GetServiceEndpoint returns the config file value of the Service End Point URL
func (c *CosContext) GetServiceEndpoint() (string, error) {
	return c.Config.GetString(config.ServiceEndpointURL)
}

// FileOperations interface to support
// ReadSeekerCloserOpen,
// WriteCloserOpen,
// GetFileInfo, and
// Remove
type FileOperations interface {
	ReadSeekerCloserOpen(location string) (ReadSeekerCloser, error)
	WriteCloserOpen(location string) (WriteCloser, error)
	GetFileInfo(location string) (os.FileInfo, error)
	Remove(location string) error
}

// ReadSeekerCloser a FileOperations interface
type ReadSeekerCloser interface {
	io.ReadSeeker
	io.Closer
}

// interface to hold any type that supports operations
// Write
// Close
// WriteAt
type WriteCloser interface {
	io.WriteCloser
	io.WriterAt
}

// Uploader interface that represents SDK operations for the Uploader,
// helper for dependency injection and facilitates the testing
type Uploader interface {
	Upload(input *s3manager.UploadInput,
		options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	UploadWithContext(ctx aws.Context, input *s3manager.UploadInput,
		opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	UploadWithIterator(ctx aws.Context, iter s3manager.BatchUploadIterator,
		opts ...func(*s3manager.Uploader)) error
}

// Downloader interface that represents SDK operations for the Downloader,
// helper for dependency injection and facilitates the testing
type Downloader interface {
	Download(w io.WriterAt, input *s3.GetObjectInput,
		options ...func(*s3manager.Downloader)) (n int64, err error)
	DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput,
		options ...func(*s3manager.Downloader)) (n int64, err error)
	DownloadWithIterator(ctx aws.Context, iter s3manager.BatchDownloadIterator,
		opts ...func(*s3manager.Downloader)) error
}
