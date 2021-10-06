package render

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
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

func (jsr *JSONRender) Display(_ interface{}, output interface{}, _ map[string]interface{}) error {
	// Turn off escape HTML characters
	jsr.encoder.SetEscapeHTML(false)
	switch castedOutput := output.(type) {
	case *s3.AbortMultipartUploadOutput:
		return nil
	case *s3.CopyObjectOutput:
		output = new(CopyObjectOutput)
		structMap(output, castedOutput)
	case *s3.CreateBucketOutput:
		return nil
	case *s3.CreateMultipartUploadOutput:
		output = new(CreateMultipartUploadOutput)
		structMap(output, castedOutput)
	case *s3.CompleteMultipartUploadOutput:
		output = new(CompleteMultipartUploadOutput)
		structMap(output, castedOutput)
	case *s3.DeleteBucketOutput:
		return nil
	case *s3.DeleteBucketCorsOutput:
		return nil
	case *s3.DeleteObjectOutput:
		return nil
	case *s3.DeleteObjectsOutput:
		return nil
	case *s3.GetObjectOutput:
		output = new(GetObjectOutput)
		structMap(output, castedOutput)
	// "Pretty-print" JSON just leaves out fields and creates a tedious maintenance burden.
	// If a user wants JSON, just give them everything, even with nulls, empty strings, or empty lists
	// If this violates OneCloud Compliance, say so on the PR review.
	//case *s3.GetBucketWebsiteOutput:
	//output = new(GetBucketWebsiteOutput)
	//structMap(output, castedOutput)
	case *s3.HeadBucketOutput:
		return nil
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
	case *s3.ListPartsOutput:
		output = new(ListPartsOutput)
		structMap(output, castedOutput)
	case *s3.PutObjectOutput:
		output = new(PutObjectOutput)
		structMap(output, castedOutput)
	case *s3.PutBucketCorsOutput:
		return nil
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

type CopyObjectOutput struct {
	CopyObjectResult *s3.CopyObjectResult `json:",omitempty"`
}

type CreateMultipartUploadOutput struct {
	Bucket   *string `json:",omitempty"`
	UploadId *string `json:",omitempty"`
	Key      *string `json:",omitempty"`
}

type CompleteMultipartUploadOutput struct {
	ETag     *string `json:",omitempty"`
	Bucket   *string `json:",omitempty"`
	Location *string `json:",omitempty"`
	Key      *string `json:",omitempty"`
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
	WebsiteRedirectLocation *string            `json:",omitempty"`
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

type PutObjectOutput struct {
	ETag *string `json:",omitempty"`
}

type UploadOutput struct {
	Location string `json:",omitempty"`
	UploadID string `json:",omitempty"`
}
type UploadPartOutput struct {
	ETag *string `json:",omitempty"`
}

type UploadPartCopyOutput struct {
	CopyPartResult *s3.CopyPartResult `json:",omitempty"`
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
