package utils

import (
	"io"
	"os"

	"github.com/IBM/ibmcloud-cos-cli/config"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3iface"
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
	ClientGen        func(*session.Session) s3iface.S3API
	ListKnownRegions ListKnownRegions

	FileOperations
}

// GetClient generates an S3 client to make requests through Go SDK
func (c *CosContext) GetClient(region string) s3iface.S3API {
	cfg := new(aws.Config).WithRegion(region)
	sess := c.Session.Copy(cfg)
	return c.ClientGen(sess)
}

// IsHMAC - Check if the auth the user uses is HMAC
func (c *CosContext) IsHMAC() bool {
	isHMAC, err := c.Config.GetBoolWithDefault(config.HMACProvided, false)
	if err != nil {
		panic(err)
	}
	return isHMAC
}

// FileOperations interface to support
// ReadSeekerCloserOpen,
// WriteCloserOpen,
// GetFileInfo, and
// Remove
type FileOperations interface {
	ReadSeekerCloserOpen(location string) (ReadSeekerCloser, error)
	WriteCloserOpen(location string) (io.WriteCloser, error)
	GetFileInfo(location string) (os.FileInfo, error)
	Remove(location string) error
}

// ReadSeekerCloser a FileOperations interface
type ReadSeekerCloser interface {
	io.ReadSeeker
	io.Closer
}
