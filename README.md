# IBM Cloud Object Storage CLI plug-in

This plug-in for the IBM Cloud CLI allows users to interact with [IBM Cloud Object Storage][ibm-cos] services entirely from the command line.

```yaml
NAME:
  ibmcloud cos -

USAGE:
  ibmcloud cos command [arguments...] [command options]

COMMANDS:
  bucket-class-get            Returns the class type of the specified bucket.
  bucket-cors-delete          Delete the CORS configuration from a bucket.
  bucket-cors-get             Get the CORS configuration from a bucket.
  bucket-cors-put             Set the CORS configuration on a bucket.
  bucket-create               Create a new bucket.
  bucket-delete               Delete an existing bucket.
  bucket-head                 Determine if a specified bucket exists in your account.
  bucket-location-get         Get the region and class of a bucket.
  bucket-website-delete       Remove static website configuration from a bucket.
  bucket-website-get          Get the static website configuration on a bucket.
  bucket-website-put          Set static website configuration on a bucket.
  buckets                     List all the buckets in your IBM Cloud Object Storage account.
  buckets-extended            List all the extended buckets with pagination support.
  config                      Changes plugin configuration
  download                    Download objects from S3 concurrently.
  multipart-upload-abort      Abort a multipart upload instance.
  multipart-upload-complete   Complete an existing multipart upload instance.
  multipart-upload-create     Create a new multipart upload instance.
  multipart-uploads           This operation lists in-progress multipart uploads.
  object-copy                 Copy an object from one bucket to another.
  object-delete               Delete an object from a bucket.
  object-get                  Download an object from a bucket.
  object-head                 Determine if an object exists within a bucket.
  object-put                  Upload an object to a bucket.
  objects                     List all objects in a specific bucket.
  objects-delete              Delete multiple objects from a bucket
  part-upload-copy            Upload a part by copying data from an existing object.
  parts                       Display the list of uploaded parts of an object.
  upload                      Upload objects from S3 concurrently.
  wait                        Wait until a particular condition is satisfied.  
  help, h                     Show help
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
