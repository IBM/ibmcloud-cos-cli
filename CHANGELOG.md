# CHANGELOG

## 1.7.0

### Content

#### Features

* Key Protect and Hyper Protect Crypto Services(HPCS) Support
* S3 Compatible Object Lock Support
* Update to use Go SDK 1.10.0

#### Defect Fixes

* COSSDK-99295: <https://github.com/IBM/ibmcloud-cos-cli/issues/4>

## 1.6.0

### Content

#### Features

* One Rate Active Plan Support
* Aspera High-Speed Transfer Support

## 1.5.0

### Content

#### Features

* Cloud Object Storage Replication
* Update to use Go SDK 1.9.0

## 1.4.0

### Content

#### Features

* Object Versioning
* Object Tagging
* Public Block Access
* IBM Cloud Security and Compliance Center Support
* Update to use Go SDK 1.8.0

## 1.3.1

### Content

#### Defect Fixes

* Update Default COS Endpoint URL

## 1.3.0

### Content

#### Features

* Static Website support
* Smart Tier support

## 1.2.4

### Content

#### Features

* COS CLI support for s390x platform
* Update to use Go SDK 1.7.0

## 1.2.3

### Content

#### Defect Fixes

* Enable trace logging only when environment variable is explicitly set
* Update to use Go SDK 1.6.1

## 1.2.2

### Content

#### Defect Fixes

* Update to use Go SDK 1.6.0

## 1.2.1

### Content

#### Defect Fixes

* Added "Deprecated" to description of `--json` flag.

## 1.2.0

### Content

* IBM OneCloud Compliance
  * The following commands have changed names.  Legacy names remain for backwards compatibility but are deemed deprecated.
    * `abort-multipart-upload` -> `multipart-upload-abort`
    * `complete-multipart-upload` -> `multipart-upload-complete`
    * `copy-object` -> `object-copy`
    * `create-bucket` -> `bucket-create`
    * `create-multipart-upload` -> `multipart-upload-create`
    * `delete-bucket` -> `bucket-delete`
    * `delete-bucket-cors` -> `bucket-cors-delete`
    * `delete-object` -> `object-delete`
    * `delete-objects` -> `objects-delete`
    * `get-bucket-cors` -> `bucket-cors-get`
    * `get-bucket-location` -> `bucket-location-get`
    * `get-object` -> `object-get`
    * `head-bucket` -> `bucket-head`
    * `head-object` -> `object-head`
    * `list-buckets` -> `buckets`
    * `list-multipart-uploads` -> `multipart-uploads`
    * `list-objects` -> `objects`
    * `list-parts` -> `parts`
    * `put-bucket-cors` -> `bucket-cors-put`
    * `put-object` -> `object-put`
    * `upload-part` -> `part-upload`
    * `upload-part-copy` -> `part-copy-upload`
  * The `--json` flag has changed to `--output json`.  `--json` remains for backwards compatibility but is deemed deprecated.
* Configurable Cloud Object Storage Endpoints

## 1.1.3

### Content

#### Features

* COS CLI support for ppc64le platform
* Update to use Go SDK 1.3.2

#### Defect Fixes

* COSSDK-68546: <https://github.com/IBM/ibmcloud-cos-cli/issues/1>

## 1.1.2

### Content

#### Defect Fixes

* COSSDK-62346: <https://github.ibm.com/objectstore/cases/issues/410>

## 1.1.1

### Content

#### Features

* Update to use Go SDK 1.3.0

#### Defect Fixes

* COSSDK-62419: <https://github.ibm.com/objectstore/objectstorage-issues/issues/679>
* COSSDK-63162: <https://github.ibm.com/Bluemix/bluemix-cli/issues/2753>

## 1.1.0

### Content

#### Features

* Extended List Buckets Support
* Concurrent S3Manager Upload and Download Support

#### Defect Fixes

* COSSDK-56304: <https://github.ibm.com/objectstore/objectstorage-issues/issues/667>
* COSSDK-56336: <https://github.ibm.com/objectstore/objectstorage-issues/issues/657>

## 1.0.0

### Content

#### Features

* Initial release
