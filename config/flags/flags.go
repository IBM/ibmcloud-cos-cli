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

	FlagConcurrency = cli.StringFlag{
		Name:  Concurrency,
		Usage: T("The number of goroutines to spin up in parallel per call to Upload when sending parts. Default value is 5."),
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

	FlagContentLength = cli.StringFlag{
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
		Usage: T("A `STRUCTURE` using JSON syntax in a file. See IBM Cloud Documentation."),
	}

	FlagDelete = cli.StringFlag{
		Name:   Delete,
		Usage:  T("A `STRUCTURE` using either shorthand or JSON syntax. See IBM Cloud Documentation."),
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

	FlagFile = cli.StringFlag{
		Name:   File,
		Usage:  T("The `PATH` to the file to upload."),
		Hidden: true,
	}

	FlagMultipartUpload = cli.StringFlag{
		Name:   MultipartUpload,
		Usage:  T("A `STRUCTURE` using either shorthand or JSON syntax. See IBM Cloud Documentation."),
		Hidden: true,
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

	FlagKmsEncryptionAlgorithm = cli.StringFlag{
		Name:   KmsEncryptionAlgorithm,
		Usage:  T("The `ALGORITHM` and `SIZE` to use with the encryption key stored by using key protect."),
		Hidden: true,
	}

	FlagKmsRootKeyCrn = cli.StringFlag{
		Name:   KmsRootKeyCrn,
		Usage:  T("The `CUSTOMERROOTKEYCRN` of the KMS  root key associated with the bucket for data encryption."),
		Hidden: true,
	}

	FlagLeavePartsOnErrors = cli.BoolFlag{
		Name:  LeavePartsOnErrors,
		Usage: T("Setting this value to true will cause the SDK to avoid calling AbortMultipartUpload on a failure, leaving all successfully uploaded parts on S3 for manual recovery."),
	}

	FlagMarker = cli.StringFlag{
		Name:  Marker,
		Usage: T("Specifies the `KEY` to start with when listing objects in a bucket."),
	}

	FlagMaxItems = cli.StringFlag{
		Name:  MaxItems,
		Usage: T("The total `NUMBER` of items to return in the command's output."),
	}

	FlagMaxUploadParts = cli.StringFlag{
		Name:  MaxUploadParts,
		Usage: T("Max number of `PARTS` which will be uploaded to S3 that calculates the part size of the object to be uploaded.  Limit is 10,000 parts."),
	}

	FlagMetadataDirective = cli.StringFlag{
		Name:  MetadataDirective,
		Usage: T("Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request. `DIRECTIVE` values: COPY, REPLACE."),
	}

	FlagMetadata = cli.StringFlag{
		Name:  Metadata,
		Usage: T("A `STRUCTURE` using JSON syntax. See IBM Cloud Documentation."),
	}

	FlagPageSize = cli.StringFlag{
		Name:  PageSize,
		Usage: T("The `SIZE` of each page to get in the service call. This does not affect the number of items returned in the command's output. Setting a smaller page size results in more calls to the COS service, retrieving fewer items in each call. This can help prevent the service calls from timing out."),
	}

	FlagPartNumber = cli.StringFlag{
		Name:   PartNumber,
		Usage:  T("Part `NUMBER` of part being uploaded. This is a positive integer between 1 and 10,000."),
		Hidden: true,
	}

	FlagPartNumberMarker = cli.StringFlag{
		Name:  PartNumberMarker,
		Usage: T("Part number `VALUE` after which listing begins"),
	}

	FlagPartSize = cli.StringFlag{
		Name:  PartSize,
		Usage: T("The buffer `SIZE` (in bytes) to use when buffering data into chunks and ending them as parts to S3. The minimum allowed part size is 5MB."),
	}

	FlagPrefix = cli.StringFlag{
		Name:  Prefix,
		Usage: T("Limits the response to keys that begin with the specified `PREFIX`."),
	}

	FlagPublicAccessBlockConfiguration = cli.StringFlag{
		Name:  PublicAccessBlockConfiguration,
		Usage: T("A `STRUCTURE` using JSON syntax. See IBM Cloud Documentation."),
	}

	FlagRange = cli.StringFlag{
		Name:  Range,
		Usage: T("Downloads the specified `RANGE` bytes of an object. For more information about the HTTP Range header, go to http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.35."),
	}

	FlagRegion = cli.StringFlag{
		Name:  Region,
		Usage: T("The `REGION` where the bucket is present. If this flag is not provided, the program will use the default option specified in config."),
	}

	FlagReplicationConfiguration = cli.StringFlag{
		Name:  ReplicationConfiguration,
		Usage: T("A `STRUCTURE` using JSON syntax. See IBM Cloud Documentation."),
	}

	FlagObjectLockConfiguration = cli.StringFlag{
		Name:  ObjectLockConfiguration,
		Usage: T("A `STRUCTURE` using JSON syntax. See IBM Cloud Documentation."),
	}

	FlagObjectLegalHold = cli.StringFlag{
		Name:  LegalHold,
		Usage: T("A `STRUCTURE` using JSON syntax. See IBM Cloud Documentation."),
	}

	FlagObjectRetention = cli.StringFlag{
		Name:  Retention,
		Usage: T("A `STRUCTURE` using JSON syntax. See IBM Cloud Documentation."),
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

	FlagTagging = cli.StringFlag{
		Name:  Tagging,
		Usage: T("A tag-set for the object encoded as URL query parameters (`Key1=Value1`). For `object-tagging-put`, a `STRUCTURE` using JSON syntax. See IBM Cloud Documentation."),
	}

	FlagTaggingDirective = cli.StringFlag{
		Name:  TaggingDirective,
		Usage: T("Specifies whether to copy the tag-set from the source object or replace it with the tag-set provided"),
	}

	FlagVersioningConfiguration = cli.StringFlag{
		Name:  VersioningConfiguration,
		Usage: T("A `STRUCTURE` using JSON syntax. See IBM Cloud Documentation."),
	}

	FlagVersionId = cli.StringFlag{
		Name:  VersionId,
		Usage: T("The object version identifier"),
	}

	FlagVersionIdMarker = cli.StringFlag{
		Name:  VersionIdMarker,
		Usage: T("Together with key-marker, specifies the object version from which to start listing object versions in a bucket"),
	}

	FlagWebsiteConfiguration = cli.StringFlag{
		Name:  WebsiteConfiguration,
		Usage: T("A `STRUCTURE` using JSON syntax. See IBM Cloud Documentation."),
	}

	FlagWebsiteRedirectLocation = cli.StringFlag{
		Name:  WebsiteRedirectLocation,
		Usage: T("If a bucket is configured as a website, redirects requests for the object key to another object in the same bucket or to a URL (`LOCATION`)."),
	}

	FlagUploadID = cli.StringFlag{
		Name:   UploadID,
		Usage:  T("Upload `ID` identifying the multipart upload."),
		Hidden: true,
	}

	FlagCRN = cli.StringFlag{
		Name:  CRN,
		Usage: T("Store your service instance ID (`CRN`) in the config."),
	}

	FlagDDL = cli.StringFlag{
		Name:  DDL,
		Usage: T("Set the default location for downloads."),
	}

	FlagIbmServiceInstanceID = cli.StringFlag{
		Name:  IbmServiceInstanceID,
		Usage: T("Sets the IBM Service Instance `ID` in the request."),
	}

	FlagKeyMarker = cli.StringFlag{
		Name:  KeyMarker,
		Usage: T("Together with another listing parameter, specifies the key from which to start listing"),
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

	FlagStyle = cli.StringFlag{
		Name:  Style,
		Usage: T("URL `STYLE` can be VHost or Path."),
	}

	FlagJSON = cli.BoolFlag{
		Name:  JSON,
		Usage: T("[Deprecated] Output returned in raw JSON format."),
	}

	FlagOutput = cli.StringFlag{
		Name:  Output,
		Usage: T("Output `FORMAT` can be only json or text."),
	}

	FlagContinuationToken = cli.StringFlag{
		Name:  ContinuationToken,
		Usage: T("A `Starting Token` to specify where to start paginating. This is the NextContinuationToken from a previously truncated response."),
	}

	FlagFetchOwner = cli.BoolFlag{
		Name:  FetchOwner,
		Usage: T("The `Boolean` is not present in listV2 by default, if you want to return owner field with each key in the result then set the fetch owner field to true."),
	}

	FlagStartAfter = cli.StringFlag{
		Name:  StartAfter,
		Usage: T(" `Start After` is where you want S3 to start listing from. S3 starts listing after this specified key. StartAfter can be any key in the bucket."),
	}

)
