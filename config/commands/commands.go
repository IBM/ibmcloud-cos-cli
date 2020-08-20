package commands

import (
	"github.com/urfave/cli"

	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/functions"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
)

var (
	// CommandCreateBucket - Create a bucket (Legacy version)
	// command:
	//	 ibmcloud cos create-bucket
	CommandCreateBucket = cli.Command{
		Name:        CreateBucket,
		Description: T("Create a new bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagIbmServiceInstanceID,
			flags.FlagClass,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketCreate,
		Hidden: true,
	}

	// CommandBucketCreate - Create a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos bucket-create
	CommandBucketCreate = cli.Command{
		Name:        BucketCreate,
		Description: T("Create a new bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagIbmServiceInstanceID,
			flags.FlagClass,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketCreate,
	}

	// CommandDeleteBucket - Delete a bucket (Legacy version)
	// command:
	//	 ibmcloud cos delete-bucket
	CommandDeleteBucket = cli.Command{
		Name:        DeleteBucket,
		Description: T("Delete an existing bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
			flags.FlagForce,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketDelete,
		Hidden: true,
	}

	// CommandBucketDelete - Delete a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos bucket-delete
	CommandBucketDelete = cli.Command{
		Name:        BucketDelete,
		Description: T("Delete an existing bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
			flags.FlagForce,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketDelete,
	}

	// CommandGetBucketLocation - Get the location of a bucket (Legacy version)
	// command:
	//	 ibmcloud cos get-bucket-location
	CommandGetBucketLocation = cli.Command{
		Name:        GetBucketLocation,
		Description: T("Get the region and class of a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketClassLocation,
		Hidden: true,
	}

	// CommandBucketLocationGet - Get the location of a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos bucket-location-get
	CommandBucketLocationGet = cli.Command{
		Name:        BucketLocationGet,
		Description: T("Get the region and class of a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketClassLocation,
	}

	// CommandGetBucketClass - Get the class of a bucket (Legacy version)
	// command:
	//	 ibmcloud cos get-bucket-class
	CommandGetBucketClass = cli.Command{
		Name:        GetBucketClass,
		Description: T("Returns the class type of the specified bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketClassLocation,
		Hidden: true,
	}

	// CommandBucketClassGet - Get the class of a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos bucket-class-get
	CommandBucketClassGet = cli.Command{
		Name:        BucketClassGet,
		Description: T("Returns the class type of the specified bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketClassLocation,
	}

	// CommandHeadBucket - Head a bucket (Legacy version)
	// command:
	//	 ibmcloud cos head-bucket
	CommandHeadBucket = cli.Command{
		Name:        HeadBucket,
		Description: T("Determine if a specified bucket exists in your account."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketHead,
		Hidden: true,
	}

	// CommandBucketHead - Head a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos bucket-head
	CommandBucketHead = cli.Command{
		Name:        BucketHead,
		Description: T("Determine if a specified bucket exists in your account."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketHead,
	}

	// CommandListBuckets - List all buckets (Legacy version)
	// command:
	//	 ibmcloud cos list-buckets
	CommandListBuckets = cli.Command{
		Name:        ListBuckets,
		Description: T("List all the buckets in your IBM Cloud Object Storage account."),
		Flags: []cli.Flag{
			flags.FlagIbmServiceInstanceID,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketsList,
		Hidden: true,
	}

	// CommandBuckets - List all buckets (OneCloud version)
	// command:
	//	 ibmcloud cos buckets
	CommandBuckets = cli.Command{
		Name:        Buckets,
		Description: T("List all the buckets in your IBM Cloud Object Storage account."),
		Flags: []cli.Flag{
			flags.FlagIbmServiceInstanceID,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketsList,
	}

	// CommandDeleteBucketCors - Delete CORS configuration from a bucket (Legacy version)
	// command:
	//	 ibmcloud cos delete-bucket-cors
	CommandDeleteBucketCors = cli.Command{
		Name:        DeleteBucketCors,
		Description: T("Delete the CORS configuration from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketCorsDelete,
		Hidden: true,
	}

	// CommandBucketCorsDelete - Delete CORS configuration from a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos delete-bucket-cors
	CommandBucketCorsDelete = cli.Command{
		Name:        BucketCorsDelete,
		Description: T("Delete the CORS configuration from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketCorsDelete,
	}

	// CommandGetBucketCors - Get CORS configuration from a bucket (Legacy version)
	// command:
	//	 ibmcloud cos get-bucket-cors
	CommandGetBucketCors = cli.Command{
		Name:        GetBucketCors,
		Description: T("Get the CORS configuration from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketCorsGet,
		Hidden: true,
	}

	// CommandBucketCorsGet - Get CORS configuration from a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos bucket-cors-get
	CommandBucketCorsGet = cli.Command{
		Name:        BucketCorsGet,
		Description: T("Get the CORS configuration from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketCorsGet,
	}

	// CommandPutBucketCors - Sets CORS configuration on a bucket (Legacy version)
	// command:
	//	 ibmcloud cos put-bucket-cors
	CommandPutBucketCors = cli.Command{
		Name:        PutBucketCors,
		Description: T("Set the CORS configuration on a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagCorsConfiguration,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketCorsPut,
		Hidden: true,
	}

	// CommandBucketCorsPut - Sets CORS configuration on a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos bucket-cors-put
	CommandBucketCorsPut = cli.Command{
		Name:        BucketCorsPut,
		Description: T("Set the CORS configuration on a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagCorsConfiguration,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketCorsPut,
	}

	// CommandListBucketsExtended - List all the extended buckets (Legacy version)
	// command:
	//	 ibmcloud cos list-bucket-extended
	CommandListBucketsExtended = cli.Command{
		Name:        ListBucketsExtended,
		Description: T("List all the extended buckets with pagination support."),
		Flags: []cli.Flag{
			flags.FlagIbmServiceInstanceID,
			flags.FlagMarker,
			flags.FlagPrefix,
			flags.FlagPageSize,
			flags.FlagMaxItems,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketsListExtended,
		Hidden: true,
	}

	// CommandBucketsExtended - List all the extended buckets (OneCloud version)
	// command:
	//	 ibmcloud cos buckets-extended
	CommandBucketsExtended = cli.Command{
		Name:        BucketsExtended,
		Description: T("List all the extended buckets with pagination support."),
		Flags: []cli.Flag{
			flags.FlagIbmServiceInstanceID,
			flags.FlagMarker,
			flags.FlagPrefix,
			flags.FlagPageSize,
			flags.FlagMaxItems,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.BucketsListExtended,
	}

	// CommandGetObject - Get object from a bucket (Legacy version)
	// command:
	//	 ibmcloud cos get-object
	CommandGetObject = cli.Command{
		Name:        GetObject,
		Description: T("Download an object from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagIfMatch,
			flags.FlagIfModifiedSince,
			flags.FlagIfNoneMatch,
			flags.FlagIfUnmodifiedSince,
			flags.FlagRange,
			flags.FlagResponseCacheControl,
			flags.FlagResponseContentDisposition,
			flags.FlagResponseContentEncoding,
			flags.FlagResponseContentLanguage,
			flags.FlagResponseContentType,
			flags.FlagResponseExpires,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		ArgsUsage: "[OUTFILE]",
		Action:    functions.ObjectGet,
		Hidden:    true,
	}

	// CommandObjectGet - Get object from a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos object-get
	CommandObjectGet = cli.Command{
		Name:        ObjectGet,
		Description: T("Download an object from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagIfMatch,
			flags.FlagIfModifiedSince,
			flags.FlagIfNoneMatch,
			flags.FlagIfUnmodifiedSince,
			flags.FlagRange,
			flags.FlagResponseCacheControl,
			flags.FlagResponseContentDisposition,
			flags.FlagResponseContentEncoding,
			flags.FlagResponseContentLanguage,
			flags.FlagResponseContentType,
			flags.FlagResponseExpires,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		ArgsUsage: "[OUTFILE]",
		Action:    functions.ObjectGet,
	}

	// CommandHeadObject - Head object from a bucket (Legacy version)
	// command:
	//	 ibmcloud cos head-object
	CommandHeadObject = cli.Command{
		Name:        HeadObject,
		Description: T("Determine if an object exists within a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagIfMatch,
			flags.FlagIfModifiedSince,
			flags.FlagIfNoneMatch,
			flags.FlagIfUnmodifiedSince,
			flags.FlagRange,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectHead,
		Hidden: true,
	}

	// CommandObjectHead - Head object from a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos object-head
	CommandObjectHead = cli.Command{
		Name:        ObjectHead,
		Description: T("Determine if an object exists within a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagIfMatch,
			flags.FlagIfModifiedSince,
			flags.FlagIfNoneMatch,
			flags.FlagIfUnmodifiedSince,
			flags.FlagRange,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectHead,
	}

	// CommandPutObject - Upload an object to a bucket (Legacy version)
	// command:
	//	 ibmcloud cos put-object
	CommandPutObject = cli.Command{
		Name:        PutObject,
		Description: T("Upload an object to a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagBody,
			flags.FlagCacheControl,
			flags.FlagContentDisposition,
			flags.FlagContentEncoding,
			flags.FlagContentLanguage,
			flags.FlagContentLength,
			flags.FlagContentMD5,
			flags.FlagContentType,
			flags.FlagMetadata,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectPut,
		Hidden: true,
	}

	// CommandObjectPut - Upload an object to a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos object-put
	CommandObjectPut = cli.Command{
		Name:        ObjectPut,
		Description: T("Upload an object to a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagBody,
			flags.FlagCacheControl,
			flags.FlagContentDisposition,
			flags.FlagContentEncoding,
			flags.FlagContentLanguage,
			flags.FlagContentLength,
			flags.FlagContentMD5,
			flags.FlagContentType,
			flags.FlagMetadata,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectPut,
	}

	// CommandDeleteObject - Delete an object from a bucket (Legacy version)
	// command:
	//	 ibmcloud cos delete-object
	CommandDeleteObject = cli.Command{
		Name:        DeleteObject,
		Description: T("Delete an object from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagRegion,
			flags.FlagForce,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectDelete,
		Hidden: true,
	}

	// CommandObjectDelete - Delete an object from a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos object-delete
	CommandObjectDelete = cli.Command{
		Name:        ObjectDelete,
		Description: T("Delete an object from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagRegion,
			flags.FlagForce,
			flags.FlagJSON,
		},
		Action: functions.ObjectDelete,
	}

	// CommandDeleteObjects - Delete multiple objects from a bucket (Legacy version)
	// command:
	//	 ibmcloud cos delete-objects
	CommandDeleteObjects = cli.Command{
		Name:        DeleteObjects,
		Description: T("Delete multiple objects from a bucket"),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagDelete,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectDeletes,
		Hidden: true,
	}

	// CommandObjectsDelete - Delete multiple object from a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos objects-delete
	CommandObjectsDelete = cli.Command{
		Name:        ObjectsDelete,
		Description: T("Delete multiple objects from a bucket"),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagDelete,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectDeletes,
	}

	// CommandCopyObject - Copy an object from one bucket to another (Legacy version)
	// command:
	//	 ibmcloud cos copy-object
	CommandCopyObject = cli.Command{
		Name:        CopyObject,
		Description: T("Copy an object from one bucket to another."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagCopySource,
			flags.FlagCacheControl,
			flags.FlagContentDisposition,
			flags.FlagContentEncoding,
			flags.FlagContentLanguage,
			flags.FlagContentType,
			flags.FlagCopySourceIfMatch,
			flags.FlagCopySourceIfModifiedSince,
			flags.FlagCopySourceIfNoneMatch,
			flags.FlagCopySourceIfUnmodifiedSince,
			flags.FlagMetadata,
			flags.FlagMetadataDirective,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectCopy,
		Hidden: true,
	}

	// CommandObjectCopy - Copy an object from one bucket to another (OneCloud version)
	// command:
	//	 ibmcloud cos object-copy
	CommandObjectCopy = cli.Command{
		Name:        ObjectCopy,
		Description: T("Copy an object from one bucket to another."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagCopySource,
			flags.FlagCacheControl,
			flags.FlagContentDisposition,
			flags.FlagContentEncoding,
			flags.FlagContentLanguage,
			flags.FlagContentType,
			flags.FlagCopySourceIfMatch,
			flags.FlagCopySourceIfModifiedSince,
			flags.FlagCopySourceIfNoneMatch,
			flags.FlagCopySourceIfUnmodifiedSince,
			flags.FlagMetadata,
			flags.FlagMetadataDirective,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectCopy,
	}

	// CommandListObjects - List all objects in a bucket (Legacy version)
	// command:
	//	 ibmcloud cos list-objects
	CommandListObjects = cli.Command{
		Name:        ListObjects,
		Description: T("List all objects in a specific bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagDelimiter,
			flags.FlagEncodingType,
			flags.FlagPrefix,
			flags.FlagMarker,
			flags.FlagPageSize,
			flags.FlagMaxItems,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectsList,
		Hidden: true,
	}

	// CommandObjects - List all objects in a bucket (OneCloud version)
	// command:
	//	 ibmcloud cos objects
	CommandObjects = cli.Command{
		Name:        Objects,
		Description: T("List all objects in a specific bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagDelimiter,
			flags.FlagEncodingType,
			flags.FlagPrefix,
			flags.FlagMarker,
			flags.FlagPageSize,
			flags.FlagMaxItems,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.ObjectsList,
	}

	// CommandMPUCreate - Create a new multipart upload instance (Legacy version)
	// command:
	//	 ibmcloud cos create-multipart-upload
	CommandMPUCreate = cli.Command{
		Name:        CreateMultipartUpload,
		Description: T("Create a new multipart upload instance."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagCacheControl,
			flags.FlagContentDisposition,
			flags.FlagContentEncoding,
			flags.FlagContentLanguage,
			flags.FlagContentType,
			flags.FlagMetadata,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.MultipartCreate,
		Hidden: true,
	}

	// CommandCreateMPU - Create a new multipart upload instance (OneCloud version)
	// command:
	//	 ibmcloud cos multipart-upload-create
	CommandCreateMPU = cli.Command{
		Name:        MultipartUploadCreate,
		Description: T("Create a new multipart upload instance."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagCacheControl,
			flags.FlagContentDisposition,
			flags.FlagContentEncoding,
			flags.FlagContentLanguage,
			flags.FlagContentType,
			flags.FlagMetadata,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.MultipartCreate,
	}

	// CommandAbortMPU - Abort a multipart upload instance (Legacy version)
	// command:
	//	 ibmcloud cos abort-multipart-upload
	CommandAbortMPU = cli.Command{
		Name:        AbortMultipartUpload,
		Description: T("Abort a multipart upload instance."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.MultipartAbort,
		Hidden: true,
	}

	// CommandMPUAbort - Abort a multipart upload instance (OneCloud version)
	// command:
	//	 ibmcloud cos multipart-upload-abort
	CommandMPUAbort = cli.Command{
		Name:        MultipartUploadAbort,
		Description: T("Abort a multipart upload instance."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.MultipartAbort,
	}

	// CommandCompleteMPU - Complete an existing multipart upload instance (Legacy version)
	// command:
	//	 ibmcloud cos complete-multipart-upload
	CommandCompleteMPU = cli.Command{
		Name:        CompleteMultipartUpload,
		Description: T("Complete an existing multipart upload instance."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagMultipartUpload,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.MultipartComplete,
		Hidden: true,
	}

	// CommandMPUComplete - Complete an existing multipart upload instance (OneCloud version)
	// command:
	//	 ibmcloud cos multipart-upload-complete
	CommandMPUComplete = cli.Command{
		Name:        MultipartUploadComplete,
		Description: T("Complete an existing multipart upload instance."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagMultipartUpload,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.MultipartComplete,
	}

	// CommandListMPUs - List in-progress multipart uploads (Legacy version)
	// command:
	//	 ibmcloud cos list-multipart-uploads
	CommandListMPUs = cli.Command{
		Name:        ListMultipartUploads,
		Description: T("This operation lists in-progress multipart uploads."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagDelimiter,
			flags.FlagEncodingType,
			flags.FlagPrefix,
			flags.FlagKeyMarker,
			flags.FlagUploadIDMarker,
			flags.FlagPageSize,
			flags.FlagMaxItems,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.MultiPartList,
		Hidden: true,
	}

	// CommandMPUs - List in-progress multipart uploads (OneCloud version)
	// command:
	//	 ibmcloud cos multipart-uploads
	CommandMPUs = cli.Command{
		Name:        MultipartUploads,
		Description: T("This operation lists in-progress multipart uploads."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagDelimiter,
			flags.FlagEncodingType,
			flags.FlagPrefix,
			flags.FlagKeyMarker,
			flags.FlagUploadIDMarker,
			flags.FlagPageSize,
			flags.FlagMaxItems,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.MultiPartList,
	}

	// CommandUploadPart - Upload a part of an object (Legacy version)
	// command:
	//	 ibmcloud cos upload-part
	CommandUploadPart = cli.Command{
		Name:        UploadPart,
		Description: T("Upload a part of an object."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagPartNumber,
			flags.FlagContentMD5,
			flags.FlagContentLength,
			flags.FlagBody,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.PartUpload,
		Hidden: true,
	}

	// CommandPartUpload - Upload a part of an object (OneCloud version)
	// command:
	//	 ibmcloud cos part-upload
	CommandPartUpload = cli.Command{
		Name:        PartUpload,
		Description: T("Upload a part of an object."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagPartNumber,
			flags.FlagContentMD5,
			flags.FlagContentLength,
			flags.FlagBody,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.PartUpload,
		Hidden: true,
	}

	// CommandCopyUploadPart - Upload a part of an object (Legacy version)
	// command:
	//	 ibmcloud cos upload-part-copy
	CommandCopyUploadPart = cli.Command{
		Name:        UploadPartCopy,
		Description: T("Upload a part by copying data from an existing object."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagPartNumber,
			flags.FlagCopySource,
			flags.FlagCopySourceIfMatch,
			flags.FlagCopySourceIfModifiedSince,
			flags.FlagCopySourceIfNoneMatch,
			flags.FlagCopySourceIfUnmodifiedSince,
			flags.FlagCopySourceRange,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.PartUploadCopy,
		Hidden: true,
	}

	// CommandPartUploadCopy - Upload a part of an object (OneCloud version)
	// command:
	//	 ibmcloud cos part-upload-copy
	CommandPartUploadCopy = cli.Command{
		Name:        PartUploadCopy,
		Description: T("Upload a part by copying data from an existing object."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagPartNumber,
			flags.FlagCopySource,
			flags.FlagCopySourceIfMatch,
			flags.FlagCopySourceIfModifiedSince,
			flags.FlagCopySourceIfNoneMatch,
			flags.FlagCopySourceIfUnmodifiedSince,
			flags.FlagCopySourceRange,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.PartUploadCopy,
	}

	// CommandListParts - List all uploaded parts of an object (Legacy version)
	// command:
	//	 ibmcloud cos parts
	CommandListParts = cli.Command{
		Name:        ListParts,
		Description: T("Display the list of uploaded parts of an object."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagPartNumberMarker,
			flags.FlagPageSize,
			flags.FlagMaxItems,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.PartsList,
		Hidden: true,
	}

	// CommandParts - List all uploaded parts of an object (OneCloud version)
	// command:
	//	 ibmcloud cos parts
	CommandParts = cli.Command{
		Name:        Parts,
		Description: T("Display the list of uploaded parts of an object."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagPartNumberMarker,
			flags.FlagPageSize,
			flags.FlagMaxItems,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.PartsList,
	}

	// CommandDownload - Download objects concurrently using S3 Transfer Manager
	// command:
	//	 ibmcloud cos download
	CommandDownload = cli.Command{
		Name:        Download,
		Description: T("Download objects from S3 concurrently."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagConcurrency,
			flags.FlagPartSize,
			flags.FlagIfMatch,
			flags.FlagIfModifiedSince,
			flags.FlagIfNoneMatch,
			flags.FlagIfUnmodifiedSince,
			flags.FlagRange,
			flags.FlagResponseCacheControl,
			flags.FlagResponseContentDisposition,
			flags.FlagResponseContentEncoding,
			flags.FlagResponseContentLanguage,
			flags.FlagResponseContentType,
			flags.FlagResponseExpires,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		ArgsUsage: "[OUTFILE]",
		Action:    functions.Download,
	}

	// CommandUpload - Upload objects concurrently using S3 Transfer Manager
	// command:
	//	 ibmcloud cos upload
	CommandUpload = cli.Command{
		Name:        Upload,
		Description: T("Upload objects from S3 concurrently."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagFile,
			flags.FlagConcurrency,
			flags.FlagMaxUploadParts,
			flags.FlagPartSize,
			flags.FlagLeavePartsOnErrors,
			flags.FlagCacheControl,
			flags.FlagContentDisposition,
			flags.FlagContentEncoding,
			flags.FlagContentLanguage,
			flags.FlagContentLength,
			flags.FlagContentMD5,
			flags.FlagContentType,
			flags.FlagMetadata,
			flags.FlagRegion,
			flags.FlagOutput,
			flags.FlagJSON,
		},
		Action: functions.Upload,
	}

	// CommandWait - Wait until a particular condition is satisfied.
	// command:
	//	 ibmcloud cos wait
	CommandWait = cli.Command{
		Name:        Wait,
		Description: T("Wait until a particular condition is satisfied.  Each subcommand polls an API until the listed requirement is met."),
		Subcommands: cli.Commands{
			CommandBucketExists,
			CommandBucketNotExists,
			CommandObjectExists,
			CommandObjectNotExists,
		},
	}

	// CommandBucketExists ...
	CommandBucketExists = cli.Command{
		Name:        BucketExists,
		Description: T("Wait until 200 response is received when polling with head-bucket.  It will poll every 5 seconds until a successful state has been reached. This will exit with a return code of 255 after 20 failed checks."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
		},
		Action: functions.WaitBucketExists,
	}

	// CommandBucketNotExists ...
	CommandBucketNotExists = cli.Command{
		Name:        BucketNotExists,
		Description: T("Wait until 404 response is received when polling with head-bucket.  It will poll every 5 seconds until a successful state has been reached.  This will exit with a return code of 255 after 20 failed checks."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
		},
		Action: functions.WaitBucketNotExists,
	}

	// CommandObjectExists ...
	CommandObjectExists = cli.Command{
		Name:        ObjectExists,
		Description: T("Wait until 200 response is received when polling with head-object.  It will poll every 5 seconds until a successful state has been reached.  This will exit with a return code of 255 after 20 failed checks."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagIfMatch,
			flags.FlagIfModifiedSince,
			flags.FlagIfNoneMatch,
			flags.FlagIfUnmodifiedSince,
			flags.FlagRange,
			flags.FlagPartNumber,
			flags.FlagRegion,
		},
		Action: functions.WaitObjectExists,
	}

	// CommandObjectNotExists ...
	CommandObjectNotExists = cli.Command{
		Name:        ObjectNotExists,
		Description: T("Wait until 404 response is received when polling with head-object.  It will poll every 5 seconds until a successful state has been reached.  This will exit with a return code of 255 after 20 failed checks."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagIfMatch,
			flags.FlagIfModifiedSince,
			flags.FlagIfNoneMatch,
			flags.FlagIfUnmodifiedSince,
			flags.FlagRange,
			flags.FlagPartNumber,
			flags.FlagRegion,
		},
		Action: functions.WaitObjectNotExists,
	}

	// CommandConfig ...
	CommandConfig = cli.Command{
		Name:        Config,
		Description: T("Changes plugin configuration"),
		Subcommands: cli.Commands{
			CommandList,
			CommandRegion,
			CommandDDL,
			CommandCRN,
			CommandHMAC,
			CommandAuth,
			CommandURLStyle,
			CommandRegionsEndpointURL,
			CommandSetEndpoint,
		},
	}

	// CommandList - (subcommand for Config)
	CommandList = cli.Command{
		Name:        List,
		Description: T("List configuration"),
		Action:      functions.ConfigList,
	}

	// CommandRegion - (subcommand for Config)
	CommandRegion = cli.Command{
		Name:        Region,
		Description: T("Store Default Region in the config."),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagRegion,
		},
		Action: functions.ConfigChangeDefaultRegion,
	}

	// CommandDDL - (subcommand for Config)
	CommandDDL = cli.Command{
		Name:        DDL,
		Description: T("Store Default Download Location in the config."),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagDDL,
		},
		Action: functions.ConfigSetDLLocation,
	}

	// CommandCRN - (subcommand for Config)
	CommandCRN = cli.Command{
		Name:        CRN,
		Description: T("Store CRN in the config."),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagCRN,
			flags.FlagForce,
		},
		Action: functions.ConfigCRN,
	}

	// CommandHMAC - (subcommand for Config)
	CommandHMAC = cli.Command{
		Name:        HMAC,
		Description: T("Store HMAC credentials in the config."),
		Flags: []cli.Flag{
			flags.FlagList,
		},
		Action: functions.ConfigAmazonHMAC,
	}

	// CommandAuth - (subcommand for Config)
	CommandAuth = cli.Command{
		Name:        Auth,
		Description: T("Switch between HMAC and IAM authentication"),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagMethod,
		},
		Action: functions.ConfigSetAuthMethod,
	}

	// CommandRegionsEndpointURL - (subcommand for Config)
	CommandRegionsEndpointURL = cli.Command{
		Name:        RegionsEndpointURL,
		Description: T("Set regions endpoint URL."),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagClear,
			flags.FlagURL,
		},
		Action: functions.ConfigSetRegionsEndpointURL,
		Hidden: true,
	}

	// CommandURLStyle - (subcommand for Config)
	CommandURLStyle = cli.Command{
		Name:        URLStyle,
		Description: T("Switch between VHost and Path URL style."),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagStyle,
		},
		Action: functions.ConfigSetURLStyle,
	}

	CommandSetEndpoint = cli.Command{
		Name:        SetEndpoint,
		Description: T("Set custom Service Endpoint for all operations."),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagClear,
			flags.FlagURL,
		},
		Action: functions.ConfigSetEndpointURL,
	}
)
