# IBM Cloud Object Storage CLI plug-in

This plug-in for the IBM Cloud CLI allows users to interact with [IBM Cloud Object Storage][ibm-cos] services entirely from the command line.

```yaml
NAME:
  ibmcloud cos -

USAGE:
  ibmcloud cos command [arguments...] [command options]

COMMANDS:
  aspera-download                         Download objects via Aspera
  aspera-upload                           Upload files or directories via Aspera
  bucket-class-get                        Get the location and billing tier of a bucket
  bucket-cors-delete                      Delete the CORS configuration from a bucket
  bucket-cors-get                         Get the CORS configuration for a bucket
  bucket-cors-put                         Set the CORS configuration on a bucket
  bucket-create                           Create a new bucket
  bucket-delete                           Delete an existing bucket
  bucket-head                             Determine if a specified bucket exists in the target region
  bucket-lifecycle-configuration-delete   Delete the lifecycle configuration from a bucket
  bucket-lifecycle-configuration-get      Get the lifecycle configuration for a bucket
  bucket-lifecycle-configuration-put      Set the lifecycle configuration on a bucket
  bucket-location-get                     Get the location and billing tier of a bucket
  bucket-replication-delete               Delete the replication configuration from a bucket
  bucket-replication-get                  Get the replication configuration for a bucket
  bucket-replication-put                  Set the replication configuration on a bucket
  bucket-versioning-get                   Get the versioning configuration for a bucket
  bucket-versioning-put                   Set the versioning configuration on a bucket
  bucket-website-delete                   Remove static website configuration from a bucket
  bucket-website-get                      Get the static website configuration on a bucket
  bucket-website-put                      Set static website configuration on a bucket
  buckets                                 List all buckets in a service instance
  buckets-extended                        List all buckets in a service instance and their provisioning codes
  config                                  Change plugin configuration
  download                                Download an object using a managed multipart transfer
  list-objects-v2                         List all objects in a specific bucket
  multipart-upload-abort                  Abort an existing multipart upload
  multipart-upload-complete               Complete an existing multipart upload
  multipart-upload-create                 Initiate a new multipart upload
  multipart-uploads                       List active multipart uploads
  object-copy                             Copy an object from one bucket to another
  object-delete                           Delete an object from a bucket
  object-get                              Download an object from a bucket
  object-head                             Get an object's size and last modified date
  object-legal-hold-get                   Get legal hold for a object
  object-legal-hold-put                   Set the legal hold on a object
  object-lock-configuration-get           Get the object lock configuration for a bucket
  object-lock-configuration-put           Set the object lock configuration on a bucket
  object-put                              Upload an object to a bucket
  object-retention-get                    Get retention on a object 
  object-retention-put                    Set retention on a object
  object-tagging-delete                   Remove tags from an object
  object-tagging-get                      Get tags for an object
  object-tagging-put                      Set tags on an object
  object-versions                         List all object versions in a specific bucket
  objects                                 List all objects in a specific bucket
  objects-delete                          Delete multiple objects from a bucket
  part-upload                             Upload a part
  part-upload-copy                        Upload a part by copying data from an existing object
  parts                                   List parts of an active multipart upload
  public-access-block-delete              Remove public access block configuration from a bucket
  public-access-block-get                 Get the public access block configuration on a bucket
  public-access-block-put                 Set public access block configuration on a bucket
  upload                                  Upload an object using a managed multipart transfer
  version                                 Print the version
  wait                                    Poll an API until a particular condition is satisfied
  help, h                                 Show help
```

## Prerequisites

- An [IBM Cloud][ibm-cloud] account
- An instance of [IBM Cloud Object Storage][cos-docs]
- The [IBM Cloud CLI][ibmcloud-cli-install]

### Install the plug-in

After you've installed `ibmcloud`, you need to install the plug-in.

1. Login to IBM Cloud (if you haven't already) with the command `ibmcloud login`.
2. Install the plug-in with `ibmcloud plugin install cloud-object-storage`.

## Getting Started

You need to provide a Service Instance ID (CRN) for the IBM Cloud Object Storage instance you want to interact with by typing `ibmcloud cos config crn`. You can find the CRN with `ibmcloud resource service-instance INSTANCE_NAME`.  Alternatively, you might open the web-based console, select **Service credentials** in the sidebar, and create a new set of credentials (or view an existing credential file that you already created).

Next, copy and paste the `resource_instance_id` from the credentials file into the CLI when it prompts you for a service instance ID. You can view your current Cloud Object Storage credentials by prompting `ibmcloud cos config list`. As the configuration file is generated by the plug-in, it's best not to edit the file manually.

**NOTE:** You must set the environment variable `IBMCLOUD_API_KEY=xxxxxxx` in order to use `aspera-upload` or `apera-download`.

### Example CLI usage

- Create a bucket in your IBM Cloud Object Storage account.
  - EXAMPLE: `ibmcloud cos bucket-create --bucket test-cli-bucket`

- Put an object in an existing bucket.
  - EXAMPLE: `ibmcloud cos object-put --bucket test-cli-bucket --key my-file --body testfile.jpg`

- Get an existing object from a bucket.
  - EXAMPLE: `ibmcloud cos object-get --bucket test-cli-bucket --key my-file output-file.jpg`

- Delete an object in an existing bucket.
  - EXAMPLE: `ibmcloud cos object-delete --bucket test-cli-bucket --key my-file`

For information on other commands, go to our plug-in [page](https://cloud.ibm.com/docs/cloud-object-storage-cli-plugin?topic=cloud-object-storage-cli-ic-use-the-ibm-cli).

## Build the plug-in from source

Building the IBM Cloud CLI COS plug-in requires the following utilities:

- [The Go programming language][golang]
- `make`

First, you need to [install Go][go-install].

To build and install the plug-in from source, run the four simple steps using Go modules:

```sh
git clone git@github.com:IBM/ibmcloud-cos-cli.git
cd ibmcloud-cos-cli
make
make install
```

**NOTE:** If you're refreshing the dependencies, use ```make clean``` option to remove the dependency files and then rebuild.

## Getting Help

Feel free to use GitHub issues for tracking bugs and feature requests, but for help use one of the following resources:

- Read a quick start guide in [IBM Docs](https://cloud.ibm.com/docs/cloud-object-storage-cli-plugin).
- Ask a question on [Stack Overflow](https://stackoverflow.com/) and tag it with `ibm` and `object-storage`.
- Open a support ticket with [IBM Cloud Support](https://cloud.ibm.com/unifiedsupport/supportcenter).
- If it turns out that you find a bug, open an [issue](https://github.com/IBM/ibmcloud-cos-cli/issues/new).

[ibm-cos]: https://cloud.ibm.com/catalog/services/cloud-object-storage
[ibmcloud-cli-install]: https://cloud.ibm.com/docs/cli?topic=cloud-cli-ibmcloud_cli
[go-install]: https://golang.org/doc/install
[golang]: https://golang.org/
[cos-docs]: https://cloud.ibm.com/docs/services/cloud-object-storage?topic=cloud-object-storage-getting-started
[ibm-cloud]: https://cloud.ibm.com
