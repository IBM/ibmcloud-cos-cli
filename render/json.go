package render

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"
)

const (
	jsonOmitEmptyString = `json:",omitempty"`
)

type JSONRender struct {
	terminal terminal.UI
	encoder  *json.Encoder
}

func NewJSONRender(terminal terminal.UI) *JSONRender {
	tmp := new(JSONRender)
	tmp.terminal = terminal
	tmp.encoder = json.NewEncoder(tmp.terminal.Writer())
	tmp.encoder.SetIndent(" ", " ")
	return tmp
}

func (jsr *JSONRender) Display(input interface{}, output interface{}, additionalParameters map[string]interface{}) error {
	// Turn off escape HTML characters
	jsr.encoder.SetEscapeHTML(false)
	switch castedOutput := output.(type) {
	case *s3.AbortMultipartUploadOutput:
		var castInput *s3.AbortMultipartUploadInput
		var ok bool
		if castInput, ok = input.(*s3.AbortMultipartUploadInput); !ok {
			return badCastError
		}
		output = AbortMultipartUploadOutput{
			Bucket: castInput.Bucket,
			Key:    castInput.Key,
		}
	case *s3.CopyObjectOutput:
		output = new(CopyObjectOutput)
		structMap(output, castedOutput)
	case *s3.CreateBucketOutput:
		var castInput *s3.CreateBucketInput
		var ok bool
		if castInput, ok = input.(*s3.CreateBucketInput); !ok {
			return badCastError
		}
		var location string
		if castInput.CreateBucketConfiguration != nil {
			location = aws.StringValue(castInput.CreateBucketConfiguration.LocationConstraint)
		}
		class := renderClass(getClassFromLocationConstraint(location))
		region := getRegionFromLocationConstraint(location)
		output = CreateBucketOutput{
			Bucket: castInput.Bucket,
			Class:  &class,
			Region: &region,
		}
	case *s3.CreateMultipartUploadOutput:
		output = new(CreateMultipartUploadOutput)
		structMap(output, castedOutput)
	case *s3.CompleteMultipartUploadOutput:
		output = new(CompleteMultipartUploadOutput)
		structMap(output, castedOutput)
	case *s3.DeleteBucketOutput:
		var castInput *s3.DeleteBucketInput
		var ok bool
		if castInput, ok = input.(*s3.DeleteBucketInput); !ok {
			return badCastError
		}
		output = DeleteBucketOutput{
			Bucket: castInput.Bucket,
		}
	case *s3.DeleteBucketCorsOutput:
		var castInput *s3.DeleteBucketCorsInput
		var ok bool
		if castInput, ok = input.(*s3.DeleteBucketCorsInput); !ok {
			return badCastError
		}
		output = DeleteBucketCorsOutput{
			Bucket: castInput.Bucket,
		}
	case *GetObjectOutputWrapper:
		output = new(GetObjectOutput)
		structMap(output, castedOutput)
	case *s3.HeadBucketOutput:
		var castInput *s3.HeadBucketInput
		var ok bool
		if castInput, ok = input.(*s3.HeadBucketInput); !ok {
			return badCastError
		}
		var region string
		if regionKey, found := additionalParameters["region"]; found {
			if region, ok = regionKey.(string); !ok {
				return badCastError
			}
		}
		output = HeadBucketOutput{
			Bucket: castInput.Bucket,
			Region: &region,
		}
	case *s3.HeadObjectOutput:
		output = new(HeadObjectOutput)
		structMap(output, castedOutput)
	case *s3.ListBucketsExtendedOutput:
		tmp := new(ListBucketsExtendedOutput)
		structMap(tmp, castedOutput)
		if tmp.IsTruncated != nil && *tmp.IsTruncated {
			tmp.NextMarker = tmp.Buckets[len(tmp.Buckets)-1].Name
		}
		output = tmp
	case *s3.ListMultipartUploadsOutput:
		output = new(ListMultipartUploadsOutput)
		structMap(output, castedOutput)
	case *s3.ListObjectsOutput:
		output = new(ListObjectsOutput)
		structMap(output, castedOutput)
	case *s3.ListObjectsV2Output:
		output = new(ListObjectsV2Output)
		structMap(output, castedOutput)
	case *s3.ListPartsOutput:
		output = new(ListPartsOutput)
		structMap(output, castedOutput)
	case *s3.PutBucketCorsOutput:
		var castInput *s3.PutBucketCorsInput
		var ok bool
		if castInput, ok = input.(*s3.PutBucketCorsInput); !ok {
			return badCastError
		}
		output = PutBucketCorsOutput{
			Bucket: castInput.Bucket,
		}
	case *s3manager.UploadOutput:
		output = new(UploadOutput)
		structMap(output, castedOutput)
	case *s3.UploadPartOutput:
		output = new(UploadPartOutput)
		structMap(output, castedOutput)
	case *s3.UploadPartCopyOutput:
		output = new(UploadPartCopyOutput)
		structMap(output, castedOutput)
	// TODO: Do we really need to define a class here if we default anyway?
	default:
		return jsr.encoder.Encode(output)
	}
	return jsr.encoder.Encode(output)
}

type AbortMultipartUploadOutput struct {
	Bucket *string `json:",omitempty"`
	Key    *string `json:",omitempty"`
}

type CreateBucketOutput struct {
	Bucket *string `json:",omitempty"`
	Region *string `json:",omitempty"`
	Class  *string `json:",omitempty"`
}

type DeleteBucketCorsOutput struct {
	Bucket *string `json:",omitempty"`
}

type DeleteBucketOutput struct {
	Bucket  *string `json:",omitempty"`
	Deleted *bool   `json:",omitempty"`
}

type HeadBucketOutput struct {
	Bucket *string `json:",omitempty"`
	Region *string `json:",omitempty"`
}

type PutBucketCorsOutput struct {
	Bucket *string `json:",omitempty"`
}

type CopyObjectOutput struct {
	CopyObjectResult    *s3.CopyObjectResult `json:",omitempty"`
	CopySourceVersionId *string              `json:",omitempty"`
	VersionId           *string              `json:",omitempty"`
}

type CreateMultipartUploadOutput struct {
	Bucket   *string `json:",omitempty"`
	UploadId *string `json:",omitempty"`
	Key      *string `json:",omitempty"`
}

type CompleteMultipartUploadOutput struct {
	ETag      *string `json:",omitempty"`
	Bucket    *string `json:",omitempty"`
	Location  *string `json:",omitempty"`
	Key       *string `json:",omitempty"`
	VersionId *string `json:",omitempty"`
}

type DeleteObjectOutput struct {
	DeleteMarker *bool   `json:",omitempty"`
	VersionId    *string `json:",omitempty"`
}

type GetObjectOutput struct {
	AcceptRanges            *string            `json:",omitempty"`
	ContentLength           *int64             `json:",omitempty"`
	ETag                    *string            `json:",omitempty"`
	ContentType             *string            `json:",omitempty"`
	LastModified            *time.Time         `json:",omitempty"`
	Metadata                map[string]*string `json:",omitempty"`
	PartsCount              *int64             `json:",omitempty"`
	CacheControl            *string            `json:",omitempty"`
	ContentDisposition      *string            `json:",omitempty"`
	ContentEncoding         *string            `json:",omitempty"`
	ContentLanguage         *string            `json:",omitempty"`
	ContentRange            *string            `json:",omitempty"`
	MissingMeta             *int64             `json:",omitempty"`
	StorageClass            *string            `json:",omitempty"`
	TagCount                *int64             `json:",omitempty"`
	VersionId               *string            `json:",omitempty"`
	WebsiteRedirectLocation *string            `json:",omitempty"`
	Expiration              *string            `json:",omitempty"`
	DownloadLocation        *string            `json:",omitempty"`
}

// create a wrapper struct around GetObjectOutput to include a
// new field DownloadLocation
type GetObjectOutputWrapper struct {
	*s3.GetObjectOutput
	DownloadLocation *string `json:",omitempty"`
}

type HeadObjectOutput struct {
	AcceptRanges            *string            `json:",omitempty"`
	CacheControl            *string            `json:",omitempty"`
	ContentDisposition      *string            `json:",omitempty"`
	ContentEncoding         *string            `json:",omitempty"`
	ContentLanguage         *string            `json:",omitempty"`
	ContentLength           *int64             `json:",omitempty"`
	ContentRange            *string            `json:",omitempty"`
	ContentType             *string            `json:",omitempty"`
	ETag                    *string            `json:",omitempty"`
	LastModified            *time.Time         `json:",omitempty"`
	Metadata                map[string]*string `json:",omitempty"`
	MissingMeta             *int64             `json:",omitempty"`
	PartsCount              *int64             `json:",omitempty"`
	StorageClass            *string            `json:",omitempty"`
	VersionId               *string            `json:",omitempty"`
	WebsiteRedirectLocation *string            `json:",omitempty"`
}
type ListBucketsExtendedOutput struct {
	Buckets     []*s3.BucketExtended `json:",omitempty"`
	Owner       *s3.Owner            `json:",omitempty"`
	IsTruncated *bool                `json:",omitempty"`
	Marker      *string              `json:",omitempty"`
	MaxKeys     *int64               `json:",omitempty"`
	Prefix      *string              `json:",omitempty"`
	NextMarker  *string              `json:",omitempty"`
}

type ListMultipartUploadsOutput struct {
	Uploads            []*s3.MultipartUpload `json:",omitempty"`
	Bucket             *string               `json:",omitempty"`
	CommonPrefixes     []*s3.CommonPrefix    `json:",omitempty"`
	Delimiter          *string               `json:",omitempty"`
	EncodingType       *string               `json:",omitempty"`
	IsTruncated        *bool                 `json:",omitempty"`
	KeyMarker          *string               `json:",omitempty"`
	MaxUploads         *int64                `json:",omitempty"`
	NextKeyMarker      *string               `json:",omitempty"`
	NextUploadIdMarker *string               `json:",omitempty"`
	Prefix             *string               `json:",omitempty"`
	UploadIdMarker     *string               `json:",omitempty"`
}

type ListObjectsOutput struct {
	Contents       []*s3.Object       `json:",omitempty"`
	CommonPrefixes []*s3.CommonPrefix `json:",omitempty"`
	Delimiter      *string            `json:",omitempty"`
	EncodingType   *string            `json:",omitempty"`
	IsTruncated    *bool              `json:",omitempty"`
	Marker         *string            `json:",omitempty"`
	MaxKeys        *int64             `json:",omitempty"`
	Name           *string            `json:",omitempty"`
	Prefix         *string            `json:",omitempty"`
	NextMarker     *string            `json:",omitempty"`
}

type ListObjectsV2Output struct {
	Contents              []*s3.Object       `json:",omitempty"`
	CommonPrefixes        []*s3.CommonPrefix `json:",omitempty"`
	ContinuationToken     *string            `json:",omitempty"`
	Delimiter             *string            `json:",omitempty"`
	EncodingType          *string            `json:",omitempty"`
	IsTruncated           *bool              `json:",omitempty"`
	KeyCount              *int64             `json:",omitempty"`
	MaxKeys               *int64             `json:",omitempty"`
	Name                  *string            `json:",omitempty"`
	NextContinuationToken *string            `json:",omitempty"`
	Prefix                *string            `json:",omitempty"`
	StartAfter            *string            `json:",omitempty"`
}

type ListPartsOutput struct {
	Initiator            *s3.Initiator `json:",omitempty"`
	Owner                *s3.Owner     `json:",omitempty"`
	Parts                []*s3.Part    `json:",omitempty"`
	StorageClass         *string       `json:",omitempty"`
	Bucket               *string       `json:",omitempty"`
	Key                  *string       `json:",omitempty"`
	IsTruncated          *bool         `json:",omitempty"`
	MaxParts             *int64        `json:",omitempty"`
	Name                 *string       `json:",omitempty"`
	NextPartNumberMarker *int64        `json:",omitempty"`
	PartNumberMarker     *int64        `json:",omitempty"`
	UploadId             *string       `json:",omitempty"`
}

type UploadOutput struct {
	Location string `json:",omitempty"`
	UploadID string `json:",omitempty"`
}
type UploadPartOutput struct {
	ETag *string `json:",omitempty"`
}

type UploadPartCopyOutput struct {
	CopyPartResult      *s3.CopyPartResult `json:",omitempty"`
	CopySourceVersionId *string            `json:",omitempty"`
}

func structMap(destination, origin interface{}) {
	// reflect destination and follow pointer to be writable
	rflxDestination := reflect.ValueOf(destination).Elem()
	// reflect origin and follow pointer so match writer
	rflxOrigin := reflect.ValueOf(origin).Elem()
	// iterate over all fields of struct
	for f := 0; f < rflxDestination.NumField(); f++ {
		// grab destination field by index
		fieldDV := rflxDestination.Field(f)
		// grab source field with same name as destination
		fieldOV := rflxOrigin.FieldByName(rflxDestination.Type().Field(f).Name)
		// check that the field exists and that destination is writable
		if fieldOV.IsValid() && fieldDV.CanSet() {
			// check that the field is not a pointer to empty string
			if value, ok := fieldOV.Interface().(*string); ok && value != nil && *value == "" {
				continue
			}
			// copy from origin to destination
			fieldDV.Set(fieldOV)
		}
	}
}
