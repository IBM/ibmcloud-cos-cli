package commands

const (
	// AbortMultipartUpload Command
	AbortMultipartUpload = "abort-multipart-upload"

	// BucketExists Waiter Command
	BucketExists = "bucket-exists"

	// BucketNotExists Waiter Command
	BucketNotExists = "bucket-not-exists"

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

	// ObjectExists Waiter Command
	ObjectExists = "object-exists"

	// ObjectNotExists Waiter Command
	ObjectNotExists = "object-not-exists"

	// PutBucketCors Command
	PutBucketCors = "put-bucket-cors"

	// PutObject Command
	PutObject = "put-object"

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
)
