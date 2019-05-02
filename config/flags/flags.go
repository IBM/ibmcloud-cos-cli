package flags

import (
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/urfave/cli"
)

var (
	FlagBody = cli.StringFlag{
		Name:  Body,
		Usage: T("Object data location (`FILE_PATH`)."),
	}

	FlagBucket = cli.StringFlag{
		Name:   Bucket,
		Usage:  T("The name (`BUCKET_NAME`) of the bucket."),
		Hidden: true,
	}

	FlagCacheControl = cli.StringFlag{
		Name:  CacheControl,
		Usage: T("Specifies `CACHING_DIRECTIVES` for the request/reply chain."),
	}

	FlagClass = cli.StringFlag{
		Name:  Class,
		Usage: T("The name (`CLASS_NAME`) of the Class."),
	}

	FlagContentDisposition = cli.StringFlag{
		Name:  ContentDisposition,
		Usage: T("Specifies presentational information (`DIRECTIVES`)."),
	}

	FlagContentEncoding = cli.StringFlag{
		Name:  ContentEncoding,
		Usage: T("Specifies what content encodings (`CONTENT_ENCODING`) have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field."),
	}

	FlagContentLanguage = cli.StringFlag{
		Name:  ContentLanguage,
		Usage: T("The `LANGUAGE` the content is in."),
	}

	FlagContentLength = cli.Int64Flag{
		Name:  ContentLength,
		Usage: T("`SIZE` of the body in bytes. This parameter is useful when the size of the body cannot be determined automatically."),
	}

	FlagContentMD5 = cli.StringFlag{
		Name:  ContentMD5,
		Usage: T("The base64-encoded 128-bit `MD5` digest of the data."),
	}

	FlagContentType = cli.StringFlag{
		Name:  ContentType,
		Usage: T("A standard `MIME` type describing the format of the object data."),
	}

	FlagCopySourceIfMatch = cli.StringFlag{
		Name:  CopySourceIfMatch,
		Usage: T("Copies the object if its entity tag (Etag) matches the specified tag (`ETAG`)."),
	}

	FlagCopySourceIfModifiedSince = cli.StringFlag{
		Name:  CopySourceIfModifiedSince,
		Usage: T("Copies the object if it has been modified since the specified time (`TIMESTAMP`)."),
	}

	FlagCopySourceIfNoneMatch = cli.StringFlag{
		Name:  CopySourceIfNoneMatch,
		Usage: T("Copies the object if its entity tag (ETag) is different than the specified tag (`ETAG`)."),
	}

	FlagCopySourceIfUnmodifiedSince = cli.StringFlag{
		Name:  CopySourceIfUnmodifiedSince,
		Usage: T("Copies the object if it hasn't been modified since the specified time (`TIMESTAMP`)."),
	}

	FlagCopySourceRange = cli.StringFlag{
		Name:  CopySourceRange,
		Usage: T("The range of bytes to copy from the source object. The range value must use the form bytes=first-last, where the first and last are the zero-based byte offsets to copy. For example, bytes=0-9 indicates that you want to copy the first ten bytes of the source. You can copy a range only if the source object is greater than 5 MB."),
	}

	FlagCopySource = cli.StringFlag{
		Name:   CopySource,
		Usage:  T("(`SOURCE`) The name of the source bucket and key name of the source object, separated by a slash (/). Must be URL-encoded."),
		Hidden: true,
	}

	FlagCorsConfiguration = cli.StringFlag{
		Name:  CorsConfiguration,
		Usage: T("The `VALUE` of CorsConfiguration to set."),
	}

	FlagCreateBucketConfiguration = cli.StringFlag{
		Name:  CreateBucketConfiguration,
		Usage: T("The `VALUE` of CreateBucketConfiguration to set."),
	}

	FlagDelete = cli.StringFlag{
		Name:   Delete,
		Usage:  T("The `VALUE` of Delete to set. Syntax: Objects=[{Key=string},{Key=string}],Quiet=boolean"),
		Hidden: true,
	}

	FlagDelimiter = cli.StringFlag{
		Name:  Delimiter,
		Usage: T("A `DELIMITER` is a character you use to group keys."),
	}

	FlagEncodingType = cli.StringFlag{
		Name:  EncodingType,
		Usage: T("Requests to encode the object keys in the response and specifies the encoding `METHOD` to use."),
	}

	FlagMultipartUpload = cli.StringFlag{
		Name:  MultipartUpload,
		Usage: T("The `VALUE` of MultipartUpload to set."),
	}

	FlagForce = cli.BoolFlag{
		Name:  Force,
		Usage: T("The operation will do not ask for confirmation."),
	}

	FlagIfMatch = cli.StringFlag{
		Name:  IfMatch,
		Usage: T("Return the object only if its entity tag (ETag) is the same as the `ETAG` specified, otherwise return a 412 (precondition failed)."),
	}

	FlagIfModifiedSince = cli.StringFlag{
		Name:  IfModifiedSince,
		Usage: T("Return the object only if it has been modified since the specified `TIMESTAMP`, otherwise return a 304 (not modified)."),
	}

	FlagIfNoneMatch = cli.StringFlag{
		Name:  IfNoneMatch,
		Usage: T("Return the object only if its entity tag (ETag) is different from the `ETAG` specified, otherwise return a 304 (not modified)."),
	}

	FlagIfUnmodifiedSince = cli.StringFlag{
		Name:  IfUnmodifiedSince,
		Usage: T("Return the object only if it has not been modified since the specified `TIMESTAMP`, otherwise return a 412 (precondition failed)."),
	}

	FlagKey = cli.StringFlag{
		Name:   Key,
		Usage:  T("The `KEY` of the object."),
		Hidden: true,
	}

	FlagMarker = cli.StringFlag{
		Name:   Marker,
		Usage:  T("Specifies the `KEY` to start with when listing objects in a bucket."),
		Hidden: true,
	}

	FlagMaxItems = cli.Int64Flag{
		Name:  MaxItems,
		Usage: T("The total `NUMBER` of items to return in the command's output."),
	}

	FlagMetadataDirective = cli.StringFlag{
		Name:  MetadataDirective,
		Usage: T("Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request. `DIRECTIVE` values: COPY, REPLACE."),
	}

	FlagMetadata = cli.StringFlag{
		Name:  Metadata,
		Usage: T("A `MAP` of metadata to store. Syntax: KeyName1=string,KeyName2=string"),
	}

	FlagStartingTokenStr = cli.StringFlag{
		Name:  StartingToken,
		Usage: T("A `TOKEN` to specify where to start paginating. This is the NextToken from a previously truncated response."),
	}

	FlagStartingTokenInt64 = cli.Int64Flag{
		Name:  StartingToken,
		Usage: T("A `TOKEN` to specify where to start paginating. This is the NextToken from a previously truncated response."),
	}

	FlagPageSize = cli.Int64Flag{
		Name:  PageSize,
		Usage: T("The `SIZE` of each page to get in the service call. This does not affect the number of items returned in the command's output. Setting a smaller page size results in more calls to the AWS service, retrieving fewer items in each call. This can help prevent the service calls from timing out."),
		Value: 1000,
	}

	FlagPartNumber = cli.Int64Flag{
		Name:   PartNumber,
		Usage:  T("Part `NUMBER` of part being uploaded. This is a positive integer between 1 and 10,000."),
		Value:  1,
		Hidden: true,
	}

	FlagPartNumberMarker = cli.Int64Flag{
		Name:   PartNumberMarker,
		Usage:  T("Part number `VALUE` after which listing begins"),
		Value:  1,
		Hidden: true,
	}

	FlagPrefix = cli.StringFlag{
		Name:  Prefix,
		Usage: T("Limits the response to keys that begin with the specified `PREFIX`."),
	}

	FlagRange = cli.StringFlag{
		Name:  Range,
		Usage: T("Downloads the specified `RANGE` bytes of an object. For more information about the HTTP Range header, go to http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35."),
	}

	FlagRegion = cli.StringFlag{
		Name:  Region,
		Usage: T("The `REGION` where the bucket is present. If this flag is not provided, the program will use the default option specified in config."),
	}

	FlagResponseCacheControl = cli.StringFlag{
		Name:  ResponseCacheControl,
		Usage: T("Sets the Cache-Control `HEADER` of the response."),
	}

	FlagResponseContentDisposition = cli.StringFlag{
		Name:  ResponseContentDisposition,
		Usage: T("Sets the Content-Disposition `HEADER` of the response."),
	}

	FlagResponseContentEncoding = cli.StringFlag{
		Name:  ResponseContentEncoding,
		Usage: T("Sets the Content-Encoding `HEADER` of the response."),
	}

	FlagResponseContentLanguage = cli.StringFlag{
		Name:  ResponseContentLanguage,
		Usage: T("Sets the Content-Language `HEADER` of the response."),
	}

	FlagResponseContentType = cli.StringFlag{
		Name:  ResponseContentType,
		Usage: T("Sets the Content-Type `HEADER` of the response."),
	}

	FlagResponseExpires = cli.StringFlag{
		Name:  ResponseExpires,
		Usage: T("Sets the Expires `HEADER` of the response."),
	}

	FlagUploadID = cli.StringFlag{
		Name:   UploadID,
		Usage:  T("Upload `ID` identifying the multipart upload."),
		Hidden: true,
	}

	FlagHMAC = cli.StringFlag{
		Name:  HMAC,
		Usage: T("Store HMAC credentials in the config."),
	}

	FlagCRN = cli.StringFlag{
		Name:  CRN,
		Usage: T("Store your service instance ID (`CRN`) in the config."),
	}

	FlagDDL = cli.StringFlag{
		Name:  DDL,
		Usage: T("Set the default location for downloads."),
	}

	FlagSwitch = cli.StringFlag{
		Name:  Switch,
		Usage: T("Switch between HMAC and IAM authentication."),
	}

	FlagIbmServiceInstanceID = cli.StringFlag{
		Name:  IbmServiceInstanceID,
		Usage: T("Sets the IBM Service Instance `ID` in the request."),
	}

	FlagKeyMarker = cli.StringFlag{
		Name:  KeyMarker,
		Usage: T("Together with upload-id-marker, this parameter specifies the multipart upload after which listing should begin."),
	}

	FlagUploadIDMarker = cli.StringFlag{
		Name:  UploadIDMarker,
		Usage: T("Together with key-marker, specifies the multipart upload after which listing should begin. If key-marker is not specified, the upload-id-marker parameter is ignored."),
	}

	FlagList = cli.BoolFlag{
		Name:  List,
		Usage: T("List configuration values."),
	}

	FlagMethod = cli.StringFlag{
		Name:  Method,
		Usage: T("Authentication `METHOD`."),
	}

	FlagURL = cli.StringFlag{
		Name:  URL,
		Usage: T("Regions endpoint `URL`."),
	}

	FlagClear = cli.BoolFlag{
		Name:  Clear,
		Usage: T("Clear the option value."),
	}
)
