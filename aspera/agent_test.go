package aspera

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials"
	"github.com/IBM/ibm-cos-sdk-go/aws/request"
	"github.com/IBM/ibm-cos-sdk-go/awstesting"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/aspera/mocks"
	sdk "github.com/IBM/ibmcloud-cos-cli/aspera/transfersdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var AsperaMetaJSON = `{
	"AccessKey":
		{
			"Id": "id",
			"Secret":"secret"
		},
	"ATSEndpoint": "https://zshengli.aspera.io:443"
}`

var AsperaMetaBrokenJSON = `{"AccessKey":`
var AsperaMetaErrorJSON = `{"__type":"UnknownError","message":"An error occurred."}`

func TestIsTransferdRunningTrue(t *testing.T) {
	mockClient := new(mocks.TransferServiceClient)
	mockClient.
		On("GetInfo", context.Background(), &sdk.InstanceInfoRequest{}).
		Return(&sdk.InstanceInfoResponse{}, nil)
	s3svc := &s3.S3{}
	a := agent{mockClient, s3svc, "apikey"}

	assert.True(t, a.IsTransferdRunning())
}

func TestIsTransferdRunningFalse(t *testing.T) {
	mockClient := new(mocks.TransferServiceClient)
	mockClient.
		On("GetInfo", context.Background(), &sdk.InstanceInfoRequest{}).
		Return(nil, fmt.Errorf("Error"))
	s3svc := &s3.S3{}
	a := agent{mockClient, s3svc, "apikey"}

	assert.False(t, a.IsTransferdRunning())
}

func TestGetBucketAsperaOK(t *testing.T) {
	s3svc := getMockS3(AsperaMetaJSON, false)
	transferdClient := new(mocks.TransferServiceClient)
	a := agent{transferdClient, s3svc, "apikey"}

	info, err := a.GetBucketAspera("a-bucket")
	expected := &BucketAsperaInfo{
		AccessKey:   &AccessKey{Id: aws.String("id"), Secret: aws.String("secret")},
		ATSEndpoint: aws.String("https://zshengli.aspera.io:443"),
	}
	assert.Nil(t, err)
	assert.Equal(t, expected, info)
}

func TestGetBucketAsperaError(t *testing.T) {
	s3svc := getMockS3(AsperaMetaErrorJSON, true)
	transferdClient := new(mocks.TransferServiceClient)
	a := agent{transferdClient, s3svc, "apikey"}

	_, err := a.GetBucketAspera("a-bucket")
	assert.NotNil(t, err)
}

func TestGetAsperaTransferSpecV2Upload(t *testing.T) {
	s3svc := getMockS3(AsperaMetaJSON, false)
	transferdClient := new(mocks.TransferServiceClient)
	a := agent{transferdClient, s3svc, "apikey"}

	sendPaths := []*sdk.Path{{Source: "/local/path/to/file", Destination: "key"}}
	sendSpec, err := a.GetAsperaTransferSpecV2("upload", "a-bucket", sendPaths)

	expectedSendSpec := `{"session_initiation":{"icos":{"api_key":"apikey","bucket":"a-bucket","ibm_service_instance_id":"SVCINTID","ibm_service_endpoint":"a-endpoint"}},"assets":{"destination_root":"/","paths":[{"source":"/local/path/to/file","destination":"key"}]},"direction":"send","title":"IBMCloud COS CLI"}`

	assert.Nil(t, err)
	assert.Equal(t, expectedSendSpec, sendSpec)
}
func TestGetAsperaTransferSpecV2Error(t *testing.T) {
	s3svc := getMockS3(AsperaMetaErrorJSON, true)
	transferdClient := new(mocks.TransferServiceClient)
	a := agent{transferdClient, s3svc, "apikey"}

	recvPaths := []*sdk.Path{{Source: "key", Destination: "/local/path/to/file"}}
	_, err := a.GetAsperaTransferSpecV2("download", "a-bucket", recvPaths)

	assert.NotNil(t, err)
}
func TestDoTransferOK(t *testing.T) {

	s3svc := getMockS3(AsperaMetaJSON, false)

	transferdClient := new(mocks.TransferServiceClient)
	monitorTransfersClient := new(mocks.TransferService_MonitorTransfersClient)

	// fake running transferd，so startServer will not be triggered
	transferdClient.
		On("GetInfo", context.Background(), &sdk.InstanceInfoRequest{}).
		Return(&sdk.InstanceInfoResponse{}, nil)
	a := agent{transferdClient, s3svc, "apikey"}

	startResp := &sdk.StartTransferResponse{
		TransferId: "test-download",
	}

	transResp := []*sdk.TransferResponse{
		{
			Status:       sdk.TransferStatus_QUEUED,
			TransferInfo: &sdk.TransferInfo{BytesTransferred: 0},
		},
		{
			Status:       sdk.TransferStatus_RUNNING,
			TransferInfo: &sdk.TransferInfo{BytesTransferred: 2048},
		},
		{
			Status:       sdk.TransferStatus_COMPLETED,
			TransferInfo: &sdk.TransferInfo{BytesTransferred: 2048},
		},
	}
	transferdClient.
		On("MonitorTransfers",
			mock.Anything,
			mock.MatchedBy(func(req *sdk.RegistrationRequest) bool { return req.TransferId[0] == "test-download" })).
		Return(monitorTransfersClient, nil)
	transferdClient.
		On("StartTransfer", mock.Anything, mock.Anything).
		Return(startResp, nil)
	for _, r := range transResp {
		monitorTransfersClient.
			On("Recv").
			Return(r, nil).Once()
	}

	out := bytes.NewBufferString("")
	bar := NewProgressBarSubscriber(2048, out)
	err := a.doTransfer(context.Background(), "download", &COSInput{
		Bucket: "a-bucket",
		Key:    "a-key",
		Path:   "/local/path",
		Sub:    bar,
	})

	assert.Nil(t, err)

	assert.Contains(t, out.String(), "Queued 0 B / 2.00 KiB")
	assert.Contains(t, out.String(), "Running 2.00 KiB / 2.00 KiB")
}

func TestDoStartTransferFail(t *testing.T) {
	s3svc := getMockS3(AsperaMetaJSON, false)

	transferdClient := new(mocks.TransferServiceClient)

	// fake running transferd，so startServer will not be triggered
	transferdClient.
		On("GetInfo", context.Background(), &sdk.InstanceInfoRequest{}).
		Return(&sdk.InstanceInfoResponse{}, nil)
	a := agent{transferdClient, s3svc, "apikey"}

	startResp := &sdk.StartTransferResponse{
		TransferId: "",
		Error:      &sdk.Error{Code: 400, Description: "Bad Request $_$"},
	}

	transferdClient.
		On("StartTransfer", mock.Anything, mock.Anything).
		Return(startResp, nil)

	err := a.doTransfer(context.Background(), "download", &COSInput{
		Bucket: "a-bucket",
		Key:    "a-key",
		Path:   "/local/path",
	})

	assert.NotNil(t, err)
	assert.EqualErrorf(t, err, "failed to start transfer: 400: Bad Request $_$", "formatted")
}

func TestDoTransferFail(t *testing.T) {
	s3svc := getMockS3(AsperaMetaJSON, false)

	transferdClient := new(mocks.TransferServiceClient)
	monitorTransfersClient := new(mocks.TransferService_MonitorTransfersClient)

	// fake running transferd，so startServer will not be triggered
	transferdClient.
		On("GetInfo", context.Background(), &sdk.InstanceInfoRequest{}).
		Return(&sdk.InstanceInfoResponse{}, nil)
	a := agent{transferdClient, s3svc, "apikey"}

	startResp := &sdk.StartTransferResponse{
		TransferId: "test-download",
	}

	failedResp := &sdk.TransferResponse{
		Status:       sdk.TransferStatus_FAILED,
		TransferInfo: &sdk.TransferInfo{ErrorDescription: "wasted"},
	}

	transferdClient.
		On("MonitorTransfers",
			mock.Anything,
			mock.MatchedBy(func(req *sdk.RegistrationRequest) bool { return req.TransferId[0] == "test-download" })).
		Return(monitorTransfersClient, nil)
	transferdClient.
		On("StartTransfer", mock.Anything, mock.Anything).
		Return(startResp, nil)

	monitorTransfersClient.
		On("Recv").
		Return(failedResp, nil).Once()

	err := a.doTransfer(context.Background(), "download", &COSInput{
		Bucket: "a-bucket",
		Key:    "a-key",
		Path:   "/local/path",
	})

	assert.NotNil(t, err)
	assert.Errorf(t, err, "transfer FAILED: wasted")
}

type stubProvider struct {
	creds   credentials.Value
	expired bool
	err     error
}

func (s *stubProvider) Retrieve() (credentials.Value, error) {
	s.expired = false
	s.creds.ProviderName = "stubProvider"
	return s.creds, s.err
}
func (s *stubProvider) IsExpired() bool {
	return s.expired
}

// Refer to ibm-cos-sdk-go/aws/request/request_test.go
func unmarshal(req *request.Request) {
	defer req.HTTPResponse.Body.Close()
	if req.Data != nil {
		json.NewDecoder(req.HTTPResponse.Body).Decode(req.Data)
	}
}

func body(str string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(str)))
}

func getMockS3(content string, fail bool) *s3.S3 {
	c := credentials.NewCredentials(&stubProvider{
		creds: credentials.Value{
			ServiceInstanceID: "SVCINTID",
		},
		expired: true,
	})

	s := awstesting.NewClient(&aws.Config{
		Credentials: c,
		Endpoint:    aws.String("a-endpoint"),
	})

	resp := &http.Response{
		StatusCode: 200,
		Body:       body(content),
	}
	if fail {
		resp = &http.Response{
			StatusCode: 400,
			Body:       body(content),
		}
	}

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.Send.Clear()
	s.Handlers.Send.PushBack(func(r *request.Request) {
		r.HTTPResponse = resp
	})

	svc := &s3.S3{Client: s}
	return svc
}
