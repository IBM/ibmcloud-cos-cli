package fields

// Definition of constants for fields the COS CLI uses to match
// fields with the flags the end users specify for the commands
const (
	// Those are the fields we pass to the S3 Client of the IBM COS SDK
	// for Go.  For more information on each field, check the following
	// library, github.com/IBM/ibm-cos-sdk-go, in the s3 service folder
	// where its api.go details the APIs, their structs and parameters
	// they use prior to making requests to the server side.
	Body                        = "Body"
	Bucket                      = "Bucket"
	CacheControl                = "CacheControl"
	Concurrency                 = "Concurrency"
	ContentDisposition          = "ContentDisposition"
	ContentEncoding             = "ContentEncoding"
	ContentLanguage             = "ContentLanguage"
	ContentLength               = "ContentLength"
	ContentMD5                  = "ContentMD5"
	ContentType                 = "ContentType"
	CopySource                  = "CopySource"
	CopySourceIfMatch           = "CopySourceIfMatch"
	CopySourceIfModifiedSince   = "CopySourceIfModifiedSince"
	CopySourceIfNoneMatch       = "CopySourceIfNoneMatch"
	CopySourceIfUnmodifiedSince = "CopySourceIfUnmodifiedSince"
	CopySourceRange             = "CopySourceRange"
	CORSConfiguration           = "CORSConfiguration"
	Delete                      = "Delete"
	Delimiter                   = "Delimiter"
	EncodingType                = "EncodingType"
	IBMServiceInstanceID        = "IBMServiceInstanceId"
	IfMatch                     = "IfMatch"
	IfModifiedSince             = "IfModifiedSince"
	IfNoneMatch                 = "IfNoneMatch"
	IfUnmodifiedSince           = "IfUnmodifiedSince"
	Key                         = "Key"
	KeyMarker                   = "KeyMarker"
	LeavePartsOnErrors          = "LeavePartsOnError"
	Marker                      = "Marker"
	MaxKeys                     = "MaxKeys"
	MaxUploadParts              = "MaxUploadParts"
	Metadata                    = "Metadata"
	MetadataDirective           = "MetadataDirective"
	MultipartUpload             = "MultipartUpload"
	PartNumber                  = "PartNumber"
	PartNumberMarker            = "PartNumberMarker"
	PartSize                    = "PartSize"
	Prefix                      = "Prefix"
	Range                       = "Range"
	ResponseCacheControl        = "ResponseCacheControl"
	ResponseContentDisposition  = "ResponseContentDisposition"
	ResponseContentEncoding     = "ResponseContentEncoding"
	ResponseContentLanguage     = "ResponseContentLanguage"
	ResponseContentType         = "ResponseContentType"
	ResponseExpires             = "ResponseExpires"
	UploadID                    = "UploadId"
	UploadIDMarker              = "UploadIdMarker"
)
