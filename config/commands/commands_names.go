package commands

const (
	// AbortMultipartUpload Command
	AbortMultipartUpload = "abort-multipart-upload"

	// Buckets Command
	Buckets = "buckets"

	// BucketClassGet Command
	BucketClassGet = "bucket-class-get"

	// BucketCorsDelete Command
	BucketCorsDelete = "bucket-cors-delete"

	// BucketCorsGet Command
	BucketCorsGet = "bucket-cors-get"

	// BucketCorsPut Command
	BucketCorsPut = "bucket-cors-put"

	// BucketCreate Command
	BucketCreate = "bucket-create"

	// BucketDelete Command
	BucketDelete = "bucket-delete"

	// BucketExists Waiter Command
	BucketExists = "bucket-exists"

	// BucketsExtended Command
	BucketsExtended = "buckets-extended"

	// BucketHead Command
	BucketHead = "bucket-head"

	// BucketNotExists Waiter Command
	BucketNotExists = "bucket-not-exists"

	// BucketLocationGet Command
	BucketLocationGet = "bucket-location-get"

	// BucketReplicationDelete Command
	BucketReplicationDelete = "bucket-replication-delete"

	// BucketReplicationGet Command
	BucketReplicationGet = "bucket-replication-get"

	// BucketReplicationPut Command
	BucketReplicationPut = "bucket-replication-put"

	// ObjectLockGet Command
	ObjectLockGet = "object-lock-configuration-get"

	// ObjectLockPut Command
	ObjectLockPut = "object-lock-configuration-put"

	// ObjectLegalHoldGet Command
	ObjectLegalHoldGet = "object-legal-hold-get"

	// ObjectLegalHoldPut Command
	ObjectLegalHoldPut = "object-legal-hold-put"

	// ObjectRetentionGet Command
	ObjectRetentionGet = "object-retention-get"

	// ObjectRetentionPut Command
	ObjectRetentionPut = "object-retention-put"

	// BucketVersioningGet Command
	BucketVersioningGet = "bucket-versioning-get"

	// BucketVersioningPut Command
	BucketVersioningPut = "bucket-versioning-put"

	// BucketWebsiteDelete Command
	BucketWebsiteDelete = "bucket-website-delete"

	// BucketWebsiteGet Command
	BucketWebsiteGet = "bucket-website-get"

	// BucketWebsitePut Command
	BucketWebsitePut = "bucket-website-put"

	// CompleteMultipartUpload Command
	CompleteMultipartUpload = "complete-multipart-upload"

	// Config Command
	Config = "config"

	// CopyObject Command
	CopyObject = "copy-object"

	// CreateBucket Command
	CreateBucket = "create-bucket"

	// CreateMultipartUpload Command
	CreateMultipartUpload = "create-multipart-upload"

	// DeleteBucket Command
	DeleteBucket = "delete-bucket"

	// DeleteBucketCors Command
	DeleteBucketCors = "delete-bucket-cors"

	// DeleteObject Command
	DeleteObject = "delete-object"

	// DeleteObjects Command
	DeleteObjects = "delete-objects"

	// Download Command from S3Manager
	Download = "download"

	//AsperaDownload
	AsperaDownload = "aspera-download"

	//AsperaUpload
	AsperaUpload = "aspera-upload"

	// GetBucketClass Command
	GetBucketClass = "get-bucket-class"

	// GetBucketCors Command
	GetBucketCors = "get-bucket-cors"

	// GetBucketLocation Command
	GetBucketLocation = "get-bucket-location"

	// GetObject Command
	GetObject = "get-object"

	// HeadBucket Command
	HeadBucket = "head-bucket"

	// HeadObject Command
	HeadObject = "head-object"

	// ListBuckets Command
	ListBuckets = "list-buckets"

	// ListBucketsExtended Command
	ListBucketsExtended = "list-buckets-extended"

	// ListMultipartUploads Command
	ListMultipartUploads = "list-multipart-uploads"

	// ListObjects Command
	ListObjects = "list-objects"

	// ListParts Command
	ListParts = "list-parts"

	// MultipartUploadAbort Command
	MultipartUploadAbort = "multipart-upload-abort"

	// MultipartUploadComplete Command
	MultipartUploadComplete = "multipart-upload-complete"

	// MultipartUploadCreate Command
	MultipartUploadCreate = "multipart-upload-create"

	// MultipartUploads Command
	MultipartUploads = "multipart-uploads"

	// ObjectCopy Command
	ObjectCopy = "object-copy"

	// ObjectDelete Command
	ObjectDelete = "object-delete"

	// ObjectsDelete Command
	ObjectsDelete = "objects-delete"

	// ObjectGet Command
	ObjectGet = "object-get"

	// ObjectHead Command
	ObjectHead = "object-head"

	// ObjectPut Command
	ObjectPut = "object-put"

	// ObjectExists Waiter Command
	ObjectExists = "object-exists"

	// ObjectNotExists Waiter Command
	ObjectNotExists = "object-not-exists"

	// Objects Command
	Objects = "objects"

	// ObjectTaggingDelete Command
	ObjectTaggingDelete = "object-tagging-delete"

	// ObjectTaggingGet Command
	ObjectTaggingGet = "object-tagging-get"

	// ObjectTaggingPut Command
	ObjectTaggingPut = "object-tagging-put"

	// ObjectVersions Command
	ObjectVersions = "object-versions"

	// PartUpload Command
	PartUpload = "part-upload"

	// PartUploadCopy Command
	PartUploadCopy = "part-upload-copy"

	// Parts Command
	Parts = "parts"

	// PutBucketCors Command
	PutBucketCors = "put-bucket-cors"

	// PutObject Command
	PutObject = "put-object"

	// PublicAccessBlockDelete Command
	PublicAccessBlockDelete = "public-access-block-delete"

	// PublicAccessBlockGet Command
	PublicAccessBlockGet = "public-access-block-get"

	// PublicAccessBlockPut Command
	PublicAccessBlockPut = "public-access-block-put"

	// Upload Command from S3Manager
	Upload = "upload"

	// UploadPart Command
	UploadPart = "upload-part"

	// UploadPartCopy Command
	UploadPartCopy = "upload-part-copy"

	// Wait Command
	Wait = "wait"

	// HMAC Subcommand for Config Auth
	HMAC = "hmac"

	// CRN Subcommand for Config CRN
	CRN = "crn"

	// DDL Subcommand for Config
	DDL = "ddl"

	// Auth Subcommand for Config
	Auth = "auth"

	// Region Subcommand for Config
	Region = "region"

	// List Subcommand for Config
	List = "list"

	// RegionsEndpointURL Subcommand for Config
	RegionsEndpointURL = "regions-endpoint"

	// URLStyle that identifies bucket and location
	URLStyle = "url-style"

	// SetEndpoint Set Service Endpoint URL
	SetEndpoint = "endpoint-url"
)
