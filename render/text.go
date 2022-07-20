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
	table := txtRender.Table([]string{T("Name"), T("Date Created")})

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
			T("Initiated"),
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
	txtRender.Say(T("Last Modified: {{.lastmodified}}",
		map[string]interface{}{
			"lastmodified": output.LastModified.Format("Monday, January 02, 2006 at 15:04:05"),
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
			T("Last Modified"),
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
			T("Last Modified"),
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
			T("Last Modified"),
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
			T("Last Modified"),
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
			T("Creation Date"),
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
			// only prefix
			txtRender.Say(T("Filter by prefix: ") + terminal.EntityNameColor(aws.StringValue(config.Filter.Prefix)))
		} else if config.Filter.Tag != nil {
			// only tag
			txtRender.Say(T("Filter by tag: ") + "(" + terminal.EntityNameColor("'"+aws.StringValue(config.Filter.Tag.Key)+"': '"+aws.StringValue(config.Filter.Tag.Value)+"'") + ")")
		} else {
			if config.Filter.And.Prefix == nil {
				// only tag(s)
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

func (txtRender *TextRender) printGetPublicAccessBlock(input interface{}, output *s3.GetPublicAccessBlockOutput) (err error) {
	config := output.PublicAccessBlockConfiguration
	txtRender.Say(T("Public Access Block Configuration"))
	txtRender.Say(T("Block Public ACLs: ") + terminal.EntityNameColor(strconv.FormatBool(aws.BoolValue(config.BlockPublicAcls))))
	txtRender.Say(T("Ignore Public ACLs: ") + terminal.EntityNameColor(strconv.FormatBool(aws.BoolValue(config.IgnorePublicAcls))))
	return
}
