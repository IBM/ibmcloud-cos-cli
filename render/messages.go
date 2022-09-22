package render

import (
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
)

//
// WarningDeleteBucket - bucket delete message
func WarningDeleteBucket(input *s3.DeleteBucketInput) string {
	return T("WARNING: This will permanently delete the bucket '{{.Bucket}}' from your account.", input)
}

// WarningDeleteBucketWebsite - bucket website delete message
func WarningDeleteBucketWebsite(input *s3.DeleteBucketWebsiteInput) string {
	return T("WARNING: This will permanently delete any bucket website configuration from the bucket '{{.Bucket}}'.", input)
}

// WarningDeleteObject - object delete message
func WarningDeleteObject(input *s3.DeleteObjectInput) string {
	return T("WARNING: This will permanently delete the object '{{.Key}}' from the bucket '{{.Bucket}}'.", input)
}

// WarningGetObject - object get message
func WarningGetObject(file, location string) string {
	return T("WARNING: An object with the name '{{.file}}' already exists at '{{.dl}}'.",
		map[string]interface{}{"file": file, "dl": location})
}

// MessageConfirmationContinue Confirmation message
func MessageConfirmationContinue() string { return T("Are you sure you would like to continue?") }

// MessageOperationCanceled - Operation cancellation message
func MessageOperationCanceled() string { return T("Operation canceled.") }

func MessageAsperaBinaryNotFound() string {
	return T("Aspera Transferd binary not found. Installing now...")
}
