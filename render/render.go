package render

import "github.com/IBM/ibm-cos-sdk-go/service/s3"

// GetBucketClassOutput type alias to separate get location from get class
type GetBucketClassOutput s3.GetBucketLocationOutput

//
// DownloadOutput type to wrap s3 download result
type DownloadOutput struct {
	_          struct{}
	TotalBytes int64
}

// Display type - JSON or Text
type Display interface {
	Display(interface{}, interface{}, map[string]interface{}) error
}
