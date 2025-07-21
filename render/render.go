package render

import (
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

// GetBucketClassOutput type alias to separate get location from get class
type GetBucketClassOutput s3.GetBucketLocationOutput

// DownloadOutput type to wrap s3 download result
type DownloadOutput struct {
	_          struct{}
	TotalBytes int64
}

// AsperaUpload for displaying different results for json and text
type AsperaUploadOutput struct {
	TotalBytes int64
}

// Display type - JSON or Text
type Display interface {
	Display(interface{}, interface{}, map[string]interface{}) error
}

type Regions struct {
	ServiceType string   `json:",omitempty"`
	Region      []string `json:",omitempty"`
}

type RegionEndpointsOutput struct {
	ServiceType *string                      `json:",omitempty"`
	Region      *string                      `json:",omitempty"`
	Endpoints   *PublicPrivateDirectEPOutput `json:",omitempty"`
	Regions     []*Regions                   `json:",omitempty"`
}

type PublicPrivateDirectEPOutput struct {
	Public  map[string]*string `json:",omitempty"`
	Private map[string]*string `json:",omitempty"`
	Direct  map[string]*string `json:",omitempty"`
}
