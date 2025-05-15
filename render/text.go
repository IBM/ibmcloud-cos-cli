package render

import (
	"strconv"
	"strings"

	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"

	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"
)

const (
	bucket = "Bucket"
	object = "Object"

	timeFormat = "Jan 02, 2006 at 15:04:05"
)

var (
	found = T("Found ")
)

func bucketDetails(bucket string) string {
	return T("Details about bucket ") + terminal.EntityNameColor(bucket) + ":"
}
func bucketClass(class string) string {
	return T("Class: ") + terminal.EntityNameColor(class)
}
func bucketRegion(region string) string {
	return T("Region: ") + terminal.EntityNameColor(region)
}

var (
	badCastError       = awserr.New("render.text.badcast", "unable to cast to expected type", nil)
	invalidFormatError = awserr.New("render.text.source.invalidformat", "format does not match expected", nil)
)

type TextRender struct {
	terminal.UI
}

func NewTextRender(terminal terminal.UI) *TextRender {
	tmp := new(TextRender)
	tmp.UI = terminal
	return tmp
}

func (txtRender *TextRender) Display(
	input, output interface{},
	additionalParameters map[string]interface{}) (err error) {

	txtRender.Ok()
	switch castedOutput := output.(type) {
	case *GetBucketClassOutput:
		return txtRender.printGetBucketClass(input, castedOutput)
	case *s3.GetBucketLocationOutput:
		return txtRender.printGetBucketLocation(input, castedOutput)
	case *s3.DeleteBucketCorsOutput:
		return txtRender.printDeleteBucketCors(input, castedOutput)
	case *s3.GetBucketCorsOutput:
		return txtRender.printGetBucketCors(input, castedOutput)
	case *s3.PutBucketCorsOutput:
		return txtRender.printPutBucketCors(input, castedOutput)
	case *s3.CreateBucketOutput:
		return txtRender.printCreateBucket(input, castedOutput)
	case *s3.DeleteBucketOutput:
		return txtRender.printDeleteBucket(input, castedOutput)
	case *s3.HeadBucketOutput:
		return txtRender.printHeadBucket(input, castedOutput, additionalParameters)
	case *s3.ListBucketsOutput:
		return txtRender.printListBuckets(input, castedOutput)
	case *s3.AbortMultipartUploadOutput:
		return txtRender.printAbortMultipartUpload(input, castedOutput)
	case *s3.CompleteMultipartUploadOutput:
		return txtRender.printCompleteMultipartUpload(input, castedOutput)
	case *s3.CreateMultipartUploadOutput:
		return txtRender.printCreateMultipartUpload(input, castedOutput)
	case *s3.ListMultipartUploadsOutput:
		return txtRender.printListMultipartUploads(input, castedOutput)
	case *s3.CopyObjectOutput:
		return txtRender.printCopyObject(input, castedOutput)
	case *s3.DeleteObjectOutput:
		return txtRender.printDeleteObject(input, castedOutput)
	case *s3.DeleteObjectsOutput:
		return txtRender.printDeleteObjects(input, castedOutput)
	case *s3.HeadObjectOutput:
		return txtRender.printHeadObject(input, castedOutput)
	case *s3.ListObjectsOutput:
		return txtRender.printListObjects(input, castedOutput)
	case *s3.ListObjectsV2Output:
		return txtRender.printListObjectsV2(input, castedOutput)
	case *s3.ListObjectVersionsOutput:
		return txtRender.printListObjectVersions(input, castedOutput)
	case *s3.PutObjectOutput:
		return txtRender.printPutObject(input, castedOutput)
	case *s3.ListPartsOutput:
		return txtRender.printListParts(input, castedOutput)
	case *s3.UploadPartOutput:
		return txtRender.printUploadPart(input, castedOutput)
	case *s3.UploadPartCopyOutput:
		return txtRender.printUploadPartCopy(input, castedOutput)
	case *s3.ListBucketsExtendedOutput:
		return txtRender.printListBucketsExtended(input, castedOutput)
	case *s3.GetObjectOutput:
		return txtRender.printGetObject(input, castedOutput)
	case *s3manager.UploadOutput:
		return txtRender.printUpload(input, castedOutput)
	case *AsperaUploadOutput:
		return txtRender.printUpload(input, nil)
	case *DownloadOutput:
		return txtRender.printDownload(input, castedOutput)
	case *s3.GetBucketVersioningOutput:
		return txtRender.printGetBucketVersioning(input, castedOutput)
	case *s3.GetBucketWebsiteOutput:
		return txtRender.printGetBucketWebsite(input, castedOutput)
	case *s3.DeleteObjectTaggingOutput:
		return txtRender.printDeleteTaggingObject(input, castedOutput)
	case *s3.GetObjectTaggingOutput:
		return txtRender.printGetTaggingObject(input, castedOutput)
	case *s3.PutObjectTaggingOutput:
		return txtRender.printPutTaggingObject(input, castedOutput)
	case *s3.GetPublicAccessBlockOutput:
		return txtRender.printGetPublicAccessBlock(input, castedOutput)
	case *s3.GetBucketReplicationOutput:
		return txtRender.printGetBucketReplication(input, castedOutput)
	case *s3.GetObjectLockConfigurationOutput:
		return txtRender.printGetObjectLockConfiguration(input, castedOutput)
	case *s3.GetObjectLegalHoldOutput:
		return txtRender.printGetObjectLegalHold(input, castedOutput)
	case *s3.GetObjectRetentionOutput:
		return txtRender.printGetObjectRetention(input, castedOutput)
	case *s3.GetBucketLifecycleConfigurationOutput:
		return txtRender.printGetBucketLifecycleConfiguration(input, castedOutput)

	default:
		return
	}
}

func (txtRender *TextRender) printGetBucketClass(input interface{}, output *GetBucketClassOutput) (err error) {
	var castInput *s3.GetBucketLocationInput
	var ok bool
	if castInput, ok = input.(*s3.GetBucketLocationInput); !ok {
		return badCastError
	}
	class := renderClass(getClassFromLocationConstraint(aws.StringValue(output.LocationConstraint)))
	// Output the successful message
	txtRender.Say(bucketDetails(aws.StringValue(castInput.Bucket)))
	txtRender.Say(bucketClass(class))
	return
}

func (txtRender *TextRender) printGetBucketLocation(input interface{}, output *s3.GetBucketLocationOutput) (err error) {
	var castInput *s3.GetBucketLocationInput
	var ok bool
	if castInput, ok = input.(*s3.GetBucketLocationInput); !ok {
		return badCastError
	}
	class := renderClass(getClassFromLocationConstraint(aws.StringValue(output.LocationConstraint)))
	region := getRegionFromLocationConstraint(aws.StringValue(output.LocationConstraint))
	// Output the successful message
	txtRender.Say(bucketDetails(aws.StringValue(castInput.Bucket)))
	txtRender.Say(bucketRegion(region))
	txtRender.Say(bucketClass(class))
	return
}

func (txtRender *TextRender) printCreateBucket(input interface{}, _ *s3.CreateBucketOutput) (err error) {
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
	// Output the successful message
	txtRender.Say(bucketDetails(aws.StringValue(castInput.Bucket)))
	txtRender.Say(bucketRegion(region))
	txtRender.Say(bucketClass(class))
	return
}

func (txtRender *TextRender) printDeleteBucketCors(input interface{}, _ *s3.DeleteBucketCorsOutput) (err error) {
	var castInput *s3.DeleteBucketCorsInput
	var ok bool
	if castInput, ok = input.(*s3.DeleteBucketCorsInput); !ok {
		return badCastError
	}
	// Output the successful message

	txtRender.Say(T("Successfully deleted CORS configuration on bucket: {{.Bucket}}",
		map[string]interface{}{bucket: terminal.EntityNameColor(*castInput.Bucket)}))
	return
}

func (txtRender *TextRender) printGetBucketCors(input interface{}, output *s3.GetBucketCorsOutput) (err error) {
	var castInput *s3.GetBucketCorsInput
	var ok bool
	if castInput, ok = input.(*s3.GetBucketCorsInput); !ok {
		return badCastError
	}
	// Output the successful message
	txtRender.Say(T("The CORS configuration of ")+
		terminal.EntityNameColor(*castInput.Bucket)+
		": \n%s", output)
	return
}

func (txtRender *TextRender) printPutBucketCors(input interface{}, _ *s3.PutBucketCorsOutput) (err error) {
	var castInput *s3.PutBucketCorsInput
	var ok bool
	if castInput, ok = input.(*s3.PutBucketCorsInput); !ok {
		return badCastError
	}
	// Output the successful message
	txtRender.Say(T("Successfully set CORS configuration on bucket: {{.Bucket}}",
		map[string]interface{}{bucket: terminal.EntityNameColor(*castInput.Bucket)}))
	return
}

func (txtRender *TextRender) printDeleteBucket(input interface{}, output *s3.DeleteBucketOutput) (err error) {
	var castInput *s3.DeleteBucketInput
	var ok bool
	if castInput, ok = input.(*s3.DeleteBucketInput); !ok {
		return badCastError
	}
	txtRender.Say(T("Successfully deleted bucket '{{.Bucket}}'. The bucket '{{.Bucket}}' will be available for reuse after 15 minutes.",
		map[string]interface{}{bucket: terminal.EntityNameColor(*castInput.Bucket)}))
	return
}

func (txtRender *TextRender) printHeadBucket(input interface{},
	output *s3.HeadBucketOutput,
	additionalParameters map[string]interface{}) (err error) {
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
	txtRender.Say(T("Bucket '{{.Bucket}}' in region {{.region}} found in your IBM Cloud Object Storage account.",
		map[string]interface{}{bucket: terminal.EntityNameColor(aws.StringValue(castInput.Bucket)),
			"region": terminal.EntityNameColor(region)}))
	return
}

func (txtRender *TextRender) printListBuckets(input interface{}, output *s3.ListBucketsOutput) (err error) {
	switch len(output.Buckets) {
	case 0:
		txtRender.Say(T("No buckets found in your account."))
		return
	case 1:
		txtRender.Say(T("1 bucket found in your account:\n"))
	default:
		txtRender.Say(strconv.Itoa(len(output.Buckets)) + T(" buckets found in your account:\n"))
	}

	// Create a table object to display each bucket in an organized fashion.
	table := txtRender.Table([]string{T("Name"), T("Date Created (UTC)")})

	for _, b := range output.Buckets {
		// Add each bucket's name and date created to the table.
		t := aws.TimeValue(b.CreationDate)
		// Format the "Date Created" in a certain way.
		table.Add(terminal.EntityNameColor(aws.StringValue(b.Name)), t.Format(timeFormat))
	}
	table.Print()
	return
}

func (txtRender *TextRender) printAbortMultipartUpload(input interface{},
	output *s3.AbortMultipartUploadOutput) (err error) {
	var castInput *s3.AbortMultipartUploadInput
	var ok bool
	if castInput, ok = input.(*s3.AbortMultipartUploadInput); !ok {
		return badCastError
	}
	txtRender.Say(T("Successfully aborted a multipart upload instance with key '{{.Key}}' and bucket '{{.Bucket}}'.",
		map[string]interface{}{"Key": terminal.EntityNameColor(*castInput.Key),
			bucket: terminal.EntityNameColor(*castInput.Bucket)}))
	return
}

func (txtRender *TextRender) printCompleteMultipartUpload(input interface{},
	output *s3.CompleteMultipartUploadOutput) (err error) {
	var castInput *s3.CompleteMultipartUploadInput
	var ok bool
	if castInput, ok = input.(*s3.CompleteMultipartUploadInput); !ok {
		return badCastError
	}
	txtRender.Say(T("Successfully uploaded '{{.Key}}' to bucket '{{.Bucket}}'.",
		map[string]interface{}{"Key": terminal.EntityNameColor(*castInput.Key),
			bucket: terminal.EntityNameColor(*castInput.Bucket)}))
	if output.VersionId != nil {
		txtRender.Say(T("Version ID: ") + terminal.EntityNameColor(aws.StringValue(output.VersionId)))
	}
	return
}

func (txtRender *TextRender) printCreateMultipartUpload(_ interface{},
	output *s3.CreateMultipartUploadOutput) (err error) {
	txtRender.Say(T("Details about your multipart upload instance:"))
	txtRender.Say("Bucket: %s", terminal.EntityNameColor(aws.StringValue(output.Bucket)))
	txtRender.Say("Key: %s", terminal.EntityNameColor(aws.StringValue(output.Key)))
	txtRender.Say("Upload ID: %s", terminal.EntityNameColor(aws.StringValue(output.UploadId)))
	return
}

func (txtRender *TextRender) printListMultipartUploads(_ interface{},
	output *s3.ListMultipartUploadsOutput) (err error) {
	if len(output.CommonPrefixes) > 0 {
		table := txtRender.Table([]string{
			T("Common Prefixes:"),
		})
		for _, prefix := range output.CommonPrefixes {
			table.Add(aws.StringValue(prefix.Prefix))
		}
		table.Print()
		txtRender.Say("")
	}
	var foundString string
	switch len(output.Uploads) {
	case 0:
		foundString = T("no multipart uploads in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Bucket)) + "'.\n"
	case 1:
		foundString = T("1 multipart upload in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Bucket)) + "':\n"
	default:
		foundString = strconv.Itoa(len(output.Uploads)) + T(" multipart uploads in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Bucket)) + "':\n"
	}
	txtRender.Say(found + foundString)

	if len(output.Uploads) > 0 {
		table := txtRender.Table([]string{
			T("UploadId"),
			T("Key"),
			T("Initiated (UTC)"),
		})
		for _, upload := range output.Uploads {
			table.Add(
				aws.StringValue(upload.UploadId),
				aws.StringValue(upload.Key),
				aws.TimeValue(upload.Initiated).Format(timeFormat),
			)
		}
		table.Print()
		txtRender.Say("")
	}
	if aws.BoolValue(output.IsTruncated) {
		txtRender.Say(T("To retrieve the next set of multipart uploads, use the following markers in the next command:"))
		txtRender.Say("--key-marker %s --upload-id-marker %s",
			terminal.EntityNameColor(aws.StringValue(output.NextKeyMarker)),
			terminal.EntityNameColor(aws.StringValue(output.NextUploadIdMarker)))
		txtRender.Say("")
	}
	return
}

func (txtRender *TextRender) printCopyObject(input interface{}, output *s3.CopyObjectOutput) (err error) {
	var castInput *s3.CopyObjectInput
	var ok bool
	if castInput, ok = input.(*s3.CopyObjectInput); !ok {
		return badCastError
	}
	source := aws.StringValue(castInput.CopySource)
	sourceBucket := strings.Split(source, "/")[0] // <bucket>/<key>?versionId=<version>
	txtRender.Say(T("Successfully copied '{{.Object}}' from bucket '{{.bucket1}}' to bucket '{{.bucket2}}'.",
		map[string]interface{}{object: terminal.EntityNameColor(aws.StringValue(castInput.Key)),
			"bucket1": terminal.EntityNameColor(sourceBucket),
			"bucket2": terminal.EntityNameColor(aws.StringValue(castInput.Bucket))}))
	if output.CopySourceVersionId != nil {
		txtRender.Say(T("Copy Source Version ID: ") + terminal.EntityNameColor(aws.StringValue(output.CopySourceVersionId)))
	}
	if output.VersionId != nil {
		txtRender.Say(T("Version ID: ") + terminal.EntityNameColor(aws.StringValue(output.VersionId)))
	}
	return
}

func (txtRender *TextRender) printDeleteObject(input interface{}, output *s3.DeleteObjectOutput) (err error) {
	var castInput *s3.DeleteObjectInput
	var ok bool
	if castInput, ok = input.(*s3.DeleteObjectInput); !ok {
		return badCastError
	}
	if output.VersionId == nil {
		// Regular delete on a never-versioned bucket
		txtRender.Say(T("Delete '{{.Key}}' from bucket '{{.Bucket}}' ran successfully.", castInput))
	} else if aws.BoolValue(output.DeleteMarker) == true {
		// Simple versioned delete that created a delete marker
		txtRender.Say(T("Delete marker created for '{{.Key}}' from bucket '{{.Bucket}}'.", castInput))
		txtRender.Say(T("Delete marker version ID: ") + terminal.EntityNameColor(aws.StringValue(output.VersionId)))
	} else {
		// Targeted version delete that removed a specific version
		txtRender.Say(T("Delete '{{.Key}}' from bucket '{{.Bucket}}' with version ID '{{.VersionId}}' ran successfully.", castInput))
	}
	return
}

func (txtRender *TextRender) printDeleteObjects(input interface{}, _ *s3.DeleteObjectsOutput) (err error) {
	txtRender.Say(T("Delete multiple objects from bucket '{{.Bucket}}' ran successfully.", input))
	return
}

func (txtRender *TextRender) printHeadObject(input interface{}, output *s3.HeadObjectOutput) (err error) {
	var castInput *s3.HeadObjectInput
	var ok bool
	if castInput, ok = input.(*s3.HeadObjectInput); !ok {
		return badCastError
	}
	txtRender.Say(T("Object '{{.Key}}' was found in bucket '{{.Bucket}}'.",
		map[string]interface{}{
			fields.Key:    terminal.EntityNameColor(aws.StringValue(castInput.Key)),
			fields.Bucket: terminal.EntityNameColor(aws.StringValue(castInput.Bucket)),
		}))
	if output.VersionId != nil {
		txtRender.Say(T("Version ID: " + terminal.EntityNameColor(aws.StringValue(output.VersionId))))
	}
	txtRender.Say(T("Object Size: {{.objectsize}}", map[string]interface{}{
		"objectsize": FormatFileSize(aws.Int64Value(output.ContentLength))},
	))
	txtRender.Say(T("Last Modified (UTC): {{.lastmodified}}",
		map[string]interface{}{
			"lastmodified": output.LastModified.Format(timeFormat),
		}))
	return
}

func (txtRender *TextRender) printListObjects(_ interface{}, output *s3.ListObjectsOutput) (err error) {
	if len(output.CommonPrefixes) > 0 {
		table := txtRender.Table([]string{
			T("Common Prefixes:"),
		})
		for _, prefix := range output.CommonPrefixes {
			table.Add(aws.StringValue(prefix.Prefix))
		}
		table.Print()
		txtRender.Say("")
	}
	var foundString string
	switch len(output.Contents) {
	case 0:
		foundString = T("no objects in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "'.\n"
	case 1:
		foundString = T("1 object in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "':\n"
	default:
		foundString = strconv.Itoa(len(output.Contents)) + T(" objects in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "':\n"
	}
	txtRender.Say(found + foundString)

	if len(output.Contents) > 0 {
		table := txtRender.Table([]string{
			T("Name"),
			T("Last Modified (UTC)"),
			T("Object Size"),
		})
		for _, object := range output.Contents {
			table.Add(
				aws.StringValue(object.Key),
				object.LastModified.Format(timeFormat),
				FormatFileSize(aws.Int64Value(object.Size)),
			)
		}
		table.Print()
		txtRender.Say("")
	}
	if aws.BoolValue(output.IsTruncated) {
		txtRender.Say(T("To retrieve the next set of objects use this Key as your --marker for the next command: "))
		txtRender.Say(terminal.EntityNameColor(aws.StringValue(output.NextMarker)))
		txtRender.Say("")
	}
	return
}

func (txtRender *TextRender) printListObjectsV2(_ interface{}, output *s3.ListObjectsV2Output) (err error) {
	if len(output.CommonPrefixes) > 0 {
		table := txtRender.Table([]string{
			T("Common Prefixes:"),
		})
		for _, prefix := range output.CommonPrefixes {
			table.Add(aws.StringValue(prefix.Prefix))
		}
		table.Print()
		txtRender.Say("")
	}
	var foundString string

	if len(output.Contents) > 0 {
		table := txtRender.Table([]string{
			T("Name"),
			T("Last Modified (UTC)"),
			T("Object Size"),
		})
		for _, object := range output.Contents {
			table.Add(
				aws.StringValue(object.Key),
				object.LastModified.Format(timeFormat),
				FormatFileSize(aws.Int64Value(object.Size)),
			)
		}
		table.Print()
		txtRender.Say("")
	}
	switch len(output.Contents) {
	case 0:
		foundString = T("no objects in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "'."
	case 1:
		foundString = T("1 object in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "'"
	default:
		foundString = strconv.Itoa(len(output.Contents)) + T(" objects in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "'"
	}
	txtRender.Say(found + foundString)
	txtRender.Say("")

	if aws.BoolValue(output.IsTruncated) {
		txtRender.Say(T("To retrieve the next set of objects use this Token as your --starting-token for the next command: "))
		txtRender.Say(terminal.EntityNameColor(aws.StringValue(output.NextContinuationToken)))
		txtRender.Say("")
	}
	return
}

func (txtRender *TextRender) printListObjectVersions(_ interface{}, output *s3.ListObjectVersionsOutput) (err error) {
	if len(output.CommonPrefixes) > 0 {
		table := txtRender.Table([]string{
			T("Common Prefixes:"),
		})
		for _, prefix := range output.CommonPrefixes {
			table.Add(aws.StringValue(prefix.Prefix))
		}
		table.Print()
		txtRender.Say("")
	}
	var foundString string
	objectVersionsLength := len(output.Versions)
	switch objectVersionsLength {
	case 0:
		foundString = T("no object versions in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "'.\n"
	case 1:
		foundString = T("1 object version in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "':\n"
	default:
		foundString = strconv.Itoa(objectVersionsLength) + T(" object versions in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "':\n"
	}
	txtRender.Say(found + foundString)

	if objectVersionsLength > 0 {
		table := txtRender.Table([]string{
			T("Name"),
			T("Version ID"),
			T("Last Modified (UTC)"),
			T("Object Size"),
			T("Is Latest"),
		})
		for _, object := range output.Versions {
			table.Add(
				aws.StringValue(object.Key),
				aws.StringValue(object.VersionId),
				object.LastModified.Format(timeFormat),
				FormatFileSize(aws.Int64Value(object.Size)),
				strconv.FormatBool(aws.BoolValue(object.IsLatest)),
			)
		}
		table.Print()
		txtRender.Say("")
	}

	deleteMarkersLength := len(output.DeleteMarkers)
	switch deleteMarkersLength {
	case 0:
		foundString = T("no delete markers in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "'.\n"
	case 1:
		foundString = T("1 delete marker in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "':\n"
	default:
		foundString = strconv.Itoa(deleteMarkersLength) + T(" delete markers in bucket '") +
			terminal.EntityNameColor(aws.StringValue(output.Name)) + "':\n"
	}
	txtRender.Say(found + foundString)

	if deleteMarkersLength > 0 {
		table := txtRender.Table([]string{
			T("Name"),
			T("Version ID"),
			T("Last Modified (UTC)"),
			T("Is Latest"),
		})
		for _, marker := range output.DeleteMarkers {
			table.Add(
				aws.StringValue(marker.Key),
				aws.StringValue(marker.VersionId),
				marker.LastModified.Format(timeFormat),
				strconv.FormatBool(aws.BoolValue(marker.IsLatest)),
			)
		}
		table.Print()
		txtRender.Say("")
	}

	if aws.BoolValue(output.IsTruncated) {
		txtRender.Say(T("To retrieve the next set of objects use the following values for the next command: "))
		txtRender.Say(T("Key Marker: ") + terminal.EntityNameColor(aws.StringValue(output.NextKeyMarker)))
		txtRender.Say(T("Version ID Marker: ") + terminal.EntityNameColor(aws.StringValue(output.NextVersionIdMarker)))
		txtRender.Say("")
	}
	return
}

func (txtRender *TextRender) printPutObject(input interface{}, output *s3.PutObjectOutput) (err error) {
	var castInput *s3.PutObjectInput
	var ok bool
	if castInput, ok = input.(*s3.PutObjectInput); !ok {
		return badCastError
	}
	txtRender.Say(T("Successfully uploaded object '{{.Key}}' to bucket '{{.Bucket}}'.",
		map[string]interface{}{fields.Key: terminal.EntityNameColor(aws.StringValue(castInput.Key)),
			fields.Bucket: terminal.EntityNameColor(aws.StringValue(castInput.Bucket))}))
	if output.VersionId != nil {
		txtRender.Say(T("Version ID: ") + terminal.EntityNameColor(aws.StringValue(output.VersionId)))
	}
	return
}

func (txtRender *TextRender) printListParts(_ interface{}, output *s3.ListPartsOutput) (err error) {
	var foundString string
	switch len(output.Parts) {
	case 0:
		foundString = T("no parts in the multipart upload for '") +
			terminal.EntityNameColor(aws.StringValue(output.Key)) + "'.\n"
	case 1:
		foundString = T("1 part in the multipart upload for '") +
			terminal.EntityNameColor(aws.StringValue(output.Key)) + "':\n"
	default:
		foundString = strconv.Itoa(len(output.Parts)) + T(" parts in the multipart upload for '") +
			terminal.EntityNameColor(aws.StringValue(output.Key)) + "':\n"
	}
	txtRender.Say(found + foundString)

	if len(output.Parts) > 0 {
		table := txtRender.Table([]string{
			T("Part Number"),
			T("Last Modified (UTC)"),
			T("ETag"),
			T("Size"),
		})
		for _, part := range output.Parts {
			table.Add(
				strconv.FormatInt(aws.Int64Value(part.PartNumber), 10),
				part.LastModified.Format(timeFormat),
				aws.StringValue(part.ETag),
				FormatFileSize(aws.Int64Value(part.Size)),
			)
		}
		table.Print()
		txtRender.Say("")
	}
	if aws.BoolValue(output.IsTruncated) {
		txtRender.Say(T("To retrieve the next set of parts, use this marker as your part-number-marker for the next command: ") +
			terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(output.NextPartNumberMarker), 10)))
		txtRender.Say("")
	}
	return
}

func (txtRender *TextRender) printUploadPart(input interface{}, output *s3.UploadPartOutput) (err error) {
	var castInput *s3.UploadPartInput
	var ok bool
	if castInput, ok = input.(*s3.UploadPartInput); !ok {
		return badCastError
	}
	txtRender.Say(T("Uploaded part '{{.part}}' of object '{{.Object}}'.", map[string]interface{}{
		"part": terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(castInput.PartNumber), 10)),
		object: terminal.EntityNameColor(aws.StringValue(castInput.Key)),
	}))

	// We need to display the ETag to the user.
	txtRender.Say("ETag: %s", terminal.EntityNameColor(aws.StringValue(output.ETag)))
	return
}

func (txtRender *TextRender) printUploadPartCopy(input interface{}, output *s3.UploadPartCopyOutput) (err error) {
	var castInput *s3.UploadPartCopyInput
	var ok bool
	if castInput, ok = input.(*s3.UploadPartCopyInput); !ok {
		return badCastError
	}
	txtRender.Say(T("Uploaded part copy '{{.part}}' of object '{{.Object}}'.",
		map[string]interface{}{
			"part": terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(castInput.PartNumber), 10)),
			object: terminal.EntityNameColor(aws.StringValue(castInput.Key)),
		}))
	txtRender.Say(T("Copy Source Version ID: ") + terminal.EntityNameColor(aws.StringValue(output.CopySourceVersionId)))
	return
}

func (txtRender *TextRender) printListBucketsExtended(_ interface{}, output *s3.ListBucketsExtendedOutput) (err error) {
	var foundString string
	switch len(output.Buckets) {
	case 0:
		foundString = T("no buckets.") + "\n"
	case 1:
		foundString = T("1 bucket:") + "\n"
	default:
		foundString = strconv.Itoa(len(output.Buckets)) + T(" buckets:") + "\n"
	}
	txtRender.Say(found + foundString)

	if len(output.Buckets) > 0 {
		table := txtRender.Table([]string{
			T("Name"),
			T("Location Constraint"),
			T("Creation Date (UTC)"),
			T("Creation Template ID"),
		})
		for _, bucket := range output.Buckets {
			table.Add(
				aws.StringValue(bucket.Name),
				aws.StringValue(bucket.LocationConstraint),
				bucket.CreationDate.Format(timeFormat),
				aws.StringValue(bucket.CreationTemplateId),
			)
		}
		table.Print()
		txtRender.Say("")
	}
	if aws.BoolValue(output.IsTruncated) {
		lastBucket := output.Buckets[len(output.Buckets)-1]
		if lastBucket != nil {
			txtRender.Say(T("To retrieve the next set of buckets, use this value marker for the next command: ") +
				terminal.EntityNameColor(aws.StringValue(lastBucket.Name)))
			txtRender.Say("")
		}
	}
	return
}

func (txtRender *TextRender) printGetObject(input interface{}, output *s3.GetObjectOutput) (err error) {
	var castInput *s3.GetObjectInput
	var ok bool
	if castInput, ok = input.(*s3.GetObjectInput); !ok {
		return badCastError
	}

	txtRender.Say(T("Successfully downloaded '{{.Key}}' from bucket '{{.Bucket}}'",
		map[string]interface{}{
			"Key":  terminal.EntityNameColor(aws.StringValue(castInput.Key)),
			bucket: terminal.EntityNameColor(aws.StringValue(castInput.Bucket)),
		}))
	if output.VersionId != nil {
		txtRender.Say(T("Version ID: ") + terminal.EntityNameColor(aws.StringValue(output.VersionId)))
	}
	txtRender.Say(FormatFileSize(aws.Int64Value(output.ContentLength)) + T(" downloaded."))
	return
}

func (txtRender *TextRender) printUpload(input interface{}, _ *s3manager.UploadOutput) (err error) {
	var castInput *s3manager.UploadInput
	var ok bool
	if castInput, ok = input.(*s3manager.UploadInput); !ok {
		return badCastError
	}

	txtRender.Say(T("Successfully uploaded object '{{.Key}}' to bucket '{{.Bucket}}'.",
		map[string]interface{}{
			fields.Key:    terminal.EntityNameColor(aws.StringValue(castInput.Key)),
			fields.Bucket: terminal.EntityNameColor(aws.StringValue(castInput.Bucket)),
		}))
	return
}

func (txtRender *TextRender) printDownload(input interface{}, output *DownloadOutput) (err error) {
	var castInput *s3.GetObjectInput
	var ok bool
	if castInput, ok = input.(*s3.GetObjectInput); !ok {
		return badCastError
	}

	txtRender.Say(T("Successfully downloaded '{{.Key}}' from bucket '{{.Bucket}}'",
		map[string]interface{}{
			"Key":  terminal.EntityNameColor(aws.StringValue(castInput.Key)),
			bucket: terminal.EntityNameColor(aws.StringValue(castInput.Bucket)),
		}))

	txtRender.Say(FormatFileSize(output.TotalBytes) + T(" downloaded."))
	return
}

func (txtRender *TextRender) printGetBucketVersioning(input interface{}, output *s3.GetBucketVersioningOutput) (err error) {
	// Output the successful message
	txtRender.Say(T("Versioning Configuration"))
	if output.Status != nil {
		txtRender.Say(T("Status: ") + terminal.EntityNameColor(aws.StringValue(output.Status)))
	} else {
		txtRender.Say(T("(empty response from server; versioning has never been configured for this bucket)"))
	}
	return
}

func (txtRender *TextRender) printGetBucketWebsite(input interface{}, output *s3.GetBucketWebsiteOutput) (err error) {
	errorDocument := output.ErrorDocument
	indexDocument := output.IndexDocument
	redirectRequests := output.RedirectAllRequestsTo

	// Output the successful message
	txtRender.Say(T("Website Configuration"))
	if indexDocument.Suffix != nil {
		txtRender.Say(T("Index Suffix: ") + terminal.EntityNameColor(aws.StringValue(indexDocument.Suffix)))
	}
	if errorDocument.Key != nil {
		txtRender.Say(T("Error Document: ") + terminal.EntityNameColor(aws.StringValue(errorDocument.Key)))
	}
	if redirectRequests != nil && redirectRequests.HostName != nil {
		txtRender.Say(T("Redirect Hostname: ") + terminal.EntityNameColor(aws.StringValue(redirectRequests.HostName)))
	}
	if redirectRequests != nil && redirectRequests.Protocol != nil {
		txtRender.Say(T("Redirect Protocol: ") + terminal.EntityNameColor(aws.StringValue(redirectRequests.Protocol)))
	}
	return
}

func (txtRender *TextRender) printDeleteTaggingObject(input interface{}, output *s3.DeleteObjectTaggingOutput) (err error) {
	// Output the successful message
	if output.VersionId != nil {
		txtRender.Say(T("Version ID: ") + terminal.EntityNameColor(aws.StringValue(output.VersionId)))
	}
	return
}

func (txtRender *TextRender) printPutTaggingObject(input interface{}, output *s3.PutObjectTaggingOutput) (err error) {
	// Output the successful message
	if output.VersionId != nil {
		txtRender.Say(T("Version ID: ") + terminal.EntityNameColor(aws.StringValue(output.VersionId)))
	}
	return
}

func (txtRender *TextRender) printGetTaggingObject(input interface{}, output *s3.GetObjectTaggingOutput) (err error) {
	// Output the successful message
	if output.VersionId != nil {
		txtRender.Say(T("Version ID: ") + terminal.EntityNameColor(aws.StringValue(output.VersionId)))
	}

	if output.TagSet != nil && len(output.TagSet) > 0 {
		// Create a table object to display each tag set in an organized fashion.
		table := txtRender.Table([]string{T("Key"), T("Value")})
		for _, entry := range output.TagSet {
			table.Add(terminal.EntityNameColor(aws.StringValue(entry.Key)), terminal.EntityNameColor(aws.StringValue(entry.Value)))
		}
		table.Print()
		txtRender.Say("")
	} else {
		txtRender.Say(T("No tags returned"))
	}
	return
}

func (txtRender *TextRender) printGetBucketReplication(input interface{}, output *s3.GetBucketReplicationOutput) (err error) {

	config := output.ReplicationConfiguration.Rules[0]
	txtRender.Say(T("Replication Configuration"))
	txtRender.Say(T("Priority: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.Priority), 10)))
	txtRender.Say(T("Status: ") + terminal.EntityNameColor(aws.StringValue(config.Status)))

	var buildString strings.Builder
	if config.Filter.Prefix == nil && config.Filter.And == nil && config.Filter.Tag == nil {
		txtRender.Say(T("Filter: ") + terminal.EntityNameColor(T("Empty")))
	} else {
		if config.Filter.Prefix != nil {
			txtRender.Say(T("Filter by prefix: ") + terminal.EntityNameColor(aws.StringValue(config.Filter.Prefix)))
		} else if config.Filter.Tag != nil {
			txtRender.Say(T("Filter by tag: ") + "(" + terminal.EntityNameColor("'"+aws.StringValue(config.Filter.Tag.Key)+"': '"+aws.StringValue(config.Filter.Tag.Value)+"'") + ")")
		} else {
			if config.Filter.And.Prefix == nil {
				for _, tag := range config.Filter.And.Tags {
					buildString.WriteString("'" + aws.StringValue(tag.Key) + "': '" + aws.StringValue(tag.Value) + "', ")
				}
				output := strings.TrimSuffix(buildString.String(), ", ")
				txtRender.Say(T("Filter by tags: ") + "(" + terminal.EntityNameColor(output) + ")")
			} else {
				// prefix and tag(s)
				for _, tag := range config.Filter.And.Tags {
					buildString.WriteString("'" + aws.StringValue(tag.Key) + "': '" + aws.StringValue(tag.Value) + "', ")
				}
				output := strings.TrimSuffix(buildString.String(), ", ")
				txtRender.Say(T("Filter by prefix: ") + terminal.EntityNameColor("'"+aws.StringValue(config.Filter.And.Prefix)+"'") + T(" and tags: ") + "(" + terminal.EntityNameColor(output) + ")")
			}
		}
	}
	txtRender.Say(T("Destination bucket: ") + terminal.EntityNameColor(aws.StringValue(config.Destination.Bucket)))
	return
}

func (txtRender *TextRender) printGetBucketLifecycleConfiguration(input interface{}, output *s3.GetBucketLifecycleConfigurationOutput) (err error) {

	if output.Rules != nil && len(output.Rules) > 0 {
		txtRender.Say(T("Lifecycle Configuration"))
		for _, config := range output.Rules {
			txtRender.Say(T("ID: ") + terminal.EntityNameColor(aws.StringValue(config.ID)))
			txtRender.Say(T("Status: ") + terminal.EntityNameColor(aws.StringValue(config.Status)))
			if config.Filter == nil || (config.Filter.Prefix == nil && config.Filter.And == nil && config.Filter.Tag == nil && config.Filter.ObjectSizeGreaterThan == nil && config.Filter.ObjectSizeLessThan == nil) {
				txtRender.Say(T("Filter: ") + terminal.EntityNameColor(T("Empty")))
			} else {
				if config.Filter.Prefix != nil {
					txtRender.Say(T("Filter by prefix: ") + terminal.EntityNameColor(aws.StringValue(config.Filter.Prefix)))
				}
				if config.Filter.Tag != nil {
					txtRender.Say(T("Filter by tag: ") + terminal.EntityNameColor("'"+aws.StringValue(config.Filter.Tag.Key)+"': '"+aws.StringValue(config.Filter.Tag.Value)+"'"))
				}
				if config.Filter.ObjectSizeGreaterThan != nil {
					txtRender.Say(T("Filter object with size greater than: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.Filter.ObjectSizeGreaterThan), 10)+" bytes"))
				}
				if config.Filter.ObjectSizeLessThan != nil {
					txtRender.Say(T("Filter object with Size less than: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.Filter.ObjectSizeLessThan), 10)+" bytes"))
				}
				if config.Filter.And != nil {
					txtRender.Say(T("Filter by And operator: "))
					if config.Filter.And.Tags != nil {
						txtRender.Say(T("Tags: "))
						if len(config.Filter.And.Tags) > 0 {
							table := txtRender.Table([]string{
								T("key"),
								T("value"),
							})
							for _, tag := range config.Filter.And.Tags {
								table.Add(*tag.Key, *tag.Value)
							}
							txtRender.Say("-------")
							table.Print()
							txtRender.Say("-------")

						}
					}
					if config.Filter.And.Prefix != nil {
						txtRender.Say(T("Prefix: ") + terminal.EntityNameColor(aws.StringValue(config.Filter.And.Prefix)))
					}
					if config.Filter.And.ObjectSizeGreaterThan != nil {
						txtRender.Say(T("ObjectSizeGreaterThan: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.Filter.And.ObjectSizeGreaterThan), 10)+" bytes"))
					}
					if config.Filter.And.ObjectSizeLessThan != nil {
						txtRender.Say(T("ObjectSizeLessThan: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.Filter.And.ObjectSizeLessThan), 10)+" bytes"))
					}
				}
			}

			if config.Expiration == nil || (config.Expiration.Date == nil && config.Expiration.Days == nil && config.Expiration.ExpiredObjectDeleteMarker == nil) {
				txtRender.Say(T("Expiration: ") + terminal.EntityNameColor(T("Empty")))
			} else {
				if config.Expiration.Days != nil {
					txtRender.Say(T("Expiration in days: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.Expiration.Days), 10)))
				}
				if config.Expiration.Date != nil {
					txtRender.Say(T("Expiration by date (UTC): ") + terminal.EntityNameColor(aws.TimeValue(config.Expiration.Date).Format(timeFormat)))
				}
				if config.Expiration.ExpiredObjectDeleteMarker != nil {
					txtRender.Say(T("Remove delete markers on expired objects: ") + terminal.EntityNameColor(strconv.FormatBool(aws.BoolValue(config.Expiration.ExpiredObjectDeleteMarker))))
				}
			}

			if config.AbortIncompleteMultipartUpload != nil {
				if config.AbortIncompleteMultipartUpload.DaysAfterInitiation == nil {
					txtRender.Say(T("AbortIncompleteMultipartUpload: ") + terminal.EntityNameColor(T("Empty")))
				} else {
					txtRender.Say(T("Abort incomplete multipart upload initiated after, in days: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.AbortIncompleteMultipartUpload.DaysAfterInitiation), 10)))
				}
			}

			if config.NoncurrentVersionExpiration != nil {
				if config.NoncurrentVersionExpiration.NewerNoncurrentVersions == nil && config.NoncurrentVersionExpiration.NoncurrentDays == nil {
					txtRender.Say(T("NoncurrentVersionExpiration: ") + terminal.EntityNameColor(T("Empty")))
				} else {
					if config.NoncurrentVersionExpiration.NewerNoncurrentVersions != nil {
						txtRender.Say(T("NoncurrentVersionExpiration by newer noncurrent versions: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.NoncurrentVersionExpiration.NewerNoncurrentVersions), 10)))
					}
					if config.NoncurrentVersionExpiration.NoncurrentDays != nil {
						txtRender.Say(T("NoncurrentVersionExpiration by noncurrent days: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.NoncurrentVersionExpiration.NoncurrentDays), 10)))
					}
				}
			}

			var foundNonCurrentVersionTransitions string
			switch len(config.NoncurrentVersionTransitions) {
			case 0:
				foundNonCurrentVersionTransitions = T("no noncurrentVersionTransitions in bucket lifecycle configuration")
			case 1:
				foundNonCurrentVersionTransitions = T("1 noncurrentVersionTransition in bucket lifecycle configuration")
			default:
				foundNonCurrentVersionTransitions = strconv.Itoa(len(config.NoncurrentVersionTransitions)) + T(" noncurrentVersionTransitions in bucket lifecycle configuration")
			}
			txtRender.Say(found + foundNonCurrentVersionTransitions)

			if len(config.NoncurrentVersionTransitions) > 0 {
				table := txtRender.Table([]string{
					T("Newer noncurrent versions"),
					T("noncurrent days"),
					T("storage class"),
				})
				for _, noncurrentVersionTransition := range config.NoncurrentVersionTransitions {
					newerNoncurrentVersions := "null"
					NoncurrentDays := "null"
					storageClass := "null"
					if noncurrentVersionTransition.NewerNoncurrentVersions != nil {
						newerNoncurrentVersions = strconv.FormatInt(aws.Int64Value(noncurrentVersionTransition.NewerNoncurrentVersions), 10)
					}
					if noncurrentVersionTransition.NoncurrentDays != nil {
						NoncurrentDays = strconv.FormatInt(aws.Int64Value(noncurrentVersionTransition.NoncurrentDays), 10)
					}
					if noncurrentVersionTransition.StorageClass != nil {
						storageClass = aws.StringValue(noncurrentVersionTransition.StorageClass)
					}
					table.Add(newerNoncurrentVersions, NoncurrentDays, storageClass)
				}
				table.Print()
				txtRender.Say("")
			}

			var foundTransitions string
			switch len(config.Transitions) {
			case 0:
				foundTransitions = T("no transitions in bucket lifecycle configuration")
			case 1:
				foundTransitions = T("1 transitions in bucket lifecycle configuration")
			default:
				foundTransitions = strconv.Itoa(len(config.Transitions)) + T(" transitions in bucket lifecycle configuration")
			}
			txtRender.Say(found + foundTransitions)

			if len(config.Transitions) > 0 {
				table := txtRender.Table([]string{
					T("days"),
					T("date (UTC)"),
					T("storage class"),
				})
				for _, transition := range config.Transitions {
					days := "null"
					date := "null"
					class := "null"
					if transition.Days != nil {
						days = strconv.FormatInt(aws.Int64Value(transition.Days), 10)
					}
					if transition.Date != nil {
						date = aws.TimeValue(transition.Date).Format(timeFormat)
					}
					if transition.StorageClass != nil {
						class = aws.StringValue(transition.StorageClass)
					}
					table.Add(days, date, class)
				}
				table.Print()
				txtRender.Say("")
			}
		}
	} else {
		txtRender.Say(T("No lifecycle configuration rules returned"))
	}

	return
}

func (txtRender *TextRender) printGetPublicAccessBlock(input interface{}, output *s3.GetPublicAccessBlockOutput) (err error) {
	config := output.PublicAccessBlockConfiguration
	txtRender.Say(T("Public Access Block Configuration"))
	txtRender.Say(T("Block Public ACLs: ") + terminal.EntityNameColor(strconv.FormatBool(aws.BoolValue(config.BlockPublicAcls))))
	txtRender.Say(T("Ignore Public ACLs: ") + terminal.EntityNameColor(strconv.FormatBool(aws.BoolValue(config.IgnorePublicAcls))))
	return
}

func (txtRender *TextRender) printGetObjectLockConfiguration(input interface{}, output *s3.GetObjectLockConfigurationOutput) (err error) {
	config := output.ObjectLockConfiguration
	txtRender.Say(T("Object Lock Configuration"))
	txtRender.Say(T("Object Lock Status: ") + terminal.EntityNameColor((aws.StringValue(config.ObjectLockEnabled))))

	if config.Rule != nil {
		txtRender.Say(T("Retention Mode: ") + terminal.EntityNameColor((aws.StringValue(config.Rule.DefaultRetention.Mode))))
		if config.Rule.DefaultRetention.Years != nil {
			txtRender.Say(T("Retention Period: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.Rule.DefaultRetention.Years), 10)) + " Years")
		} else {
			txtRender.Say(T("Retention Period: ") + terminal.EntityNameColor(strconv.FormatInt(aws.Int64Value(config.Rule.DefaultRetention.Days), 10)) + " Days")
		}
	}
	return
}

func (txtRender *TextRender) printGetObjectLegalHold(input interface{}, output *s3.GetObjectLegalHoldOutput) (err error) {
	config := output.LegalHold
	txtRender.Say(T("Legal Hold Status: ") + terminal.EntityNameColor((aws.StringValue(config.Status))))
	return
}

func (txtRender *TextRender) printGetObjectRetention(input interface{}, output *s3.GetObjectRetentionOutput) (err error) {
	config := output.Retention
	txtRender.Say(T("Retention"))
	txtRender.Say(T("Mode: ") + terminal.EntityNameColor((aws.StringValue(config.Mode))))
	txtRender.Say(T("Retain Until Date (UTC): ") + terminal.EntityNameColor(aws.TimeValue(config.RetainUntilDate).Format(timeFormat)))
	return
}
