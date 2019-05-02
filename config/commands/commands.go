package commands

import (
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/functions"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
	"github.com/urfave/cli"
)

var (
	CommandBucketCreate = cli.Command{
		Name:        CreateBucket,
		Description: T("Create a new bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagIbmServiceInstanceID,
			flags.FlagClass,
			flags.FlagRegion,
		},
		Action: functions.BucketCreate,
	}

	CommandBucketDelete = cli.Command{
		Name:        DeleteBucket,
		Description: T("Delete an existing bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
			flags.FlagForce,
		},
		Action: functions.BucketDelete,
	}

	CommandBucketGetLocation = cli.Command{
		Name:        GetBucketLocation,
		Description: T("Get the region and class of a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
		},
		Action: functions.BucketLocation,
	}

	CommandBucketGetClass = cli.Command{
		Name:        GetBucketClass,
		Description: T("Returns the class type of the specified bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
		},
		Action: functions.BucketClass,
	}

	CommandBucketHead = cli.Command{
		Name:        HeadBucket,
		Description: T("Determine if a specified bucket exists in your account."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
		},
		Action: functions.BucketHead,
	}

	CommandBucketsList = cli.Command{
		Name:        ListBuckets,
		Description: T("List all the buckets in your IBM Cloud Object Storage account."),
		Flags: []cli.Flag{
			flags.FlagIbmServiceInstanceID,
		},
		Action: functions.BucketsList,
	}

	CommandObjectGet = cli.Command{
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
		},
		ArgsUsage: "[OUTFILE]",
		Action:    functions.ObjectGet,
	}

	CommandObjectHead = cli.Command{
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
		},
		Action: functions.ObjectHead,
	}

	CommandObjectPut = cli.Command{
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
		},
		Action: functions.ObjectPut,
	}

	CommandObjectDelete = cli.Command{
		Name:        DeleteObject,
		Description: T("Delete an object from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagRegion,
			flags.FlagForce,
		},
		Action: functions.ObjectDelete,
	}

	CommandObjectDeletes = cli.Command{
		Name:        DeleteObjects,
		Description: T("Delete multiple objects from a bucket"),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagDelete,
			flags.FlagRegion,
		},
		Action: functions.ObjectDeletes,
	}

	CommandObjectCopy = cli.Command{
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
		},
		Action: functions.ObjectCopy,
	}

	CommandObjectList = cli.Command{
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
		},
		Action: functions.ObjectsList,
	}

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
		},
		Action: functions.MultipartCreate,
	}

	CommandMPUPartUpload = cli.Command{
		Name:        UploadPart,
		Description: T("Upload a part of an object."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagContentMD5,
			flags.FlagContentLength,
			flags.FlagUploadID,
			flags.FlagPartNumber,
			flags.FlagBody,
			flags.FlagRegion,
		},
		Action: functions.PartUpload,
	}

	CommandMPUPartUploadCopy = cli.Command{
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
		},
		Action: functions.PartUploadCopy,
	}

	CommandMPUPartsList = cli.Command{
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
		},
		Action: functions.PartsList,
	}

	CommandMPUAbort = cli.Command{
		Name:        AbortMultipartUpload,
		Description: T("Aborts a multipart upload instance."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagRegion,
		},
		Action: functions.MultipartAbort,
	}

	CommandMPUComplete = cli.Command{
		Name:        CompleteMultipartUpload,
		Description: T("Complete an existing multipart upload instance."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagKey,
			flags.FlagUploadID,
			flags.FlagMultipartUpload,
			flags.FlagRegion,
		},
		Action: functions.MultipartComplete,
	}

	CommandBucketCorsPut = cli.Command{
		Name:        PutBucketCors,
		Description: T("Set the CORS configuration on a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagCorsConfiguration,
			flags.FlagRegion,
		},
		Action: functions.BucketCorsPut,
	}

	CommandBucketCorsGet = cli.Command{
		Name:        GetBucketCors,
		Description: T("Get the CORS configuration from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
		},
		Action: functions.BucketCorsGet,
	}

	CommandBucketCorsDelete = cli.Command{
		Name:        DeleteBucketCors,
		Description: T("Delete the CORS configuration from a bucket."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
		},
		Action: functions.BucketCorsDelete,
	}

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

	CommandBucketExists = cli.Command{
		Name:        BucketExists,
		Description: T("Wait until 200 response is received when polling with head-bucket.  It will poll every 5 seconds until a successful state has been reached. This will exit with a return code of 255 after 20 failed checks."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
		},
		Action: functions.WaitBucketExists,
	}

	CommandBucketNotExists = cli.Command{
		Name:        BucketNotExists,
		Description: T("Wait until 404 response is received when polling with head-bucket.  It will poll every 5 seconds until a successful state has been reached.  This will exit with a return code of 255 after 20 failed checks."),
		Flags: []cli.Flag{
			flags.FlagBucket,
			flags.FlagRegion,
		},
		Action: functions.WaitBucketNotExists,
	}

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

	CommandMPUList = cli.Command{
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
		},
		Action: functions.MultiPartList,
	}

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
			CommandRegionsEndpointURL,
		},
	}

	CommandList = cli.Command{
		Name:        List,
		Description: T("List configuration"),
		Action:      functions.ConfigList,
	}

	CommandRegion = cli.Command{
		Name:        Region,
		Description: T("Store Default Region in the config."),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagRegion,
		},
		Action: functions.ConfigChangeDefaultRegion,
	}

	CommandDDL = cli.Command{
		Name:        DDL,
		Description: T("Store Default Download Location in the config."),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagDDL,
		},
		Action: functions.ConfigSetDLLocation,
	}

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

	CommandHMAC = cli.Command{
		Name:        HMAC,
		Description: T("Store HMAC credentials in the config."),
		Flags: []cli.Flag{
			flags.FlagList,
		},
		Action: functions.ConfigAmazonHMAC,
	}

	CommandAuth = cli.Command{
		Name:        Auth,
		Description: T("Switch between HMAC and IAM authentication"),
		Flags: []cli.Flag{
			flags.FlagList,
			flags.FlagMethod,
		},
		Action: functions.ConfigSetAuthMethod,
	}

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
)
