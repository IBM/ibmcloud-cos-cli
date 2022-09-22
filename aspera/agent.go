package aspera

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/IBM/ibm-cos-sdk-go/aws/request"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	sdk "github.com/IBM/ibmcloud-cos-cli/aspera/transfersdk"
)

type AccessKey struct {
	_      struct{} `type:"structure"`
	Id     *string  `type:"string"`
	Secret *string  `type:"string"`
}

type BucketAsperaInfo struct {
	AccessKey   *AccessKey `type:"structure"`
	ATSEndpoint *string    `type:"string"`
}

type COSInput struct {
	Bucket string
	Key    string
	Path   string
	Sub    Subscriber
}

type agent struct {
	sdk.TransferServiceClient
	svc    *s3.S3
	apikey string
}

func New(c *s3.S3, apikey string) (a *agent, err error) {
	cc, err := defaultConnection()
	if err != nil {
		return
	}
	client := sdk.NewTransferServiceClient(cc)

	a = &agent{client, c, apikey}
	return
}

// StartServer start transferd if it's not running
func (a *agent) StartServer(ctx context.Context) error {
	if a.IsTransferdRunning() {
		return nil
	}
	transferd := TransferdBinPath()
	err := exec.CommandContext(ctx, transferd).Start()
	if err != nil {
		return fmt.Errorf("failed to start asperatransferd(%s): %v", transferd, err)
	}
	return nil
}

// IsTransferdRunning test if transferd service is running
func (a *agent) IsTransferdRunning() bool {
	if _, err := a.GetInfo(context.Background(), &sdk.InstanceInfoRequest{}); err != nil {
		return false
	}
	return true
}

// GetBucketAspera get aspera meta information for bucket
// Add this to cos-sdk if possible
func (a *agent) GetBucketAspera(bucket string) (info *BucketAsperaInfo, err error) {
	opGetBucketAspera := &request.Operation{
		Name:       "GetBucketAspera",
		HTTPMethod: "GET",
		HTTPPath:   fmt.Sprintf("/%s?faspConnectionInfo", bucket),
	}

	output := &BucketAsperaInfo{}
	req := a.svc.NewRequest(opGetBucketAspera, nil, output)
	if err := req.Send(); err != nil {
		return nil, err
	}
	return output, nil
}

//GetICOSSpec get COS spec for a bucket
func (a *agent) GetICOSSpec(bucket string) *sdk.ICOSSpec {
	creds, _ := a.svc.Config.Credentials.Get()
	ICOSSpec := &sdk.ICOSSpec{
		ApiKey:               a.apikey,
		Bucket:               bucket,
		IbmServiceInstanceId: creds.ServiceInstanceID,
		IbmServiceEndpoint:   a.svc.Endpoint,
	}
	return ICOSSpec
}

//GetAsperaTransferSpecV2 get transfer spec version 2
func (a *agent) GetAsperaTransferSpecV2(action string, bucket string, paths []*sdk.Path) (spec string, err error) {

	ICOSSpec := a.GetICOSSpec(bucket)
	direction := "recv"
	if action == "upload" {
		direction = "send"
	}

	transferSpec := &sdk.TransferSpecV2{
		SessionInitiation: &sdk.Initiation{
			Icos: ICOSSpec,
		},
		Direction: direction,
		Assets: &sdk.Assets{
			Paths:           paths,
			DestinationRoot: "/",
		},
		Title: "IBMCloud COS CLI",
	}

	data, err := json.Marshal(transferSpec)
	if err != nil {
		return "", fmt.Errorf("unable to marshal transferspecv2: %s", err)
	}

	spec = string(data)
	return
}

// GetAsperaTransferSpecV1 get transfer spec version 1. Not being used currently.
func (a *agent) GetAsperaTransferSpecV1(action string, bucket string, paths []*sdk.Path) (spec string, err error) {
	meta, err := a.GetBucketAspera(bucket)
	if err != nil {
		return "", fmt.Errorf("unable to get aspera metadata: %s", err)
	}
	creds, err := a.svc.Config.Credentials.Get()
	if err != nil {
		return "", fmt.Errorf("unable to get aws credentials: %s", err)
	}
	credentials := fmt.Sprintf(`{"type":"token","token":{ "delegated_refresh_token": "%s"}}`, creds.Token.AccessToken)

	j, err := json.Marshal(paths)
	if err != nil {
		return "", err
	}
	jsonPaths := string(j)

	// The type of `tags` in the original proto file is string
	// so I can't use TransferSpecV1 struct directly.
	// Reported to Aspera team
	data := fmt.Sprintf(`{
		"transfer_requests": [
		  {
			"transfer_request": {
			  "paths": %s,
			  "tags": {
				"aspera": {
				  "node": {
					"storage_credentials": %s
				  }
				}
			  }
			}
		  }
		]
	  }
	  `, jsonPaths, credentials)

	url := fmt.Sprintf("%s/files/%s_setup", *meta.ATSEndpoint, action)
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Aspera-Storage-Credentials", credentials)
	req.SetBasicAuth(*meta.AccessKey.Id, *meta.AccessKey.Secret)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to %s: %v", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP Error: %v : %s", res.StatusCode, res.Request.URL)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	type reqErr struct {
		Code    int    `json:"code,omitempty"`
		Reason  string `json:"reason,omitempty"`
		Message string `json:"user_message,omitempty"`
	}

	if matched, _ := regexp.Match(`"error":`, body); matched {
		var e reqErr
		if err := json.Unmarshal(body, &e); err != nil {
			return "", fmt.Errorf("request error: %d: %s", e.Code, e.Message)
		}
	}

	specs := map[string][]map[string]map[string]interface{}{
		"transfer_specs": {
			{"transfer_spec": map[string]interface{}{}},
		},
	}
	if err := json.Unmarshal(body, &specs); err != nil {
		return "", fmt.Errorf("failed to unmarshal transferspecs: %s: %s", err, string(body))
	}
	spec_map := specs["transfer_specs"][0]["transfer_spec"]
	j, err = json.Marshal(spec_map)
	if err != nil {
		return "", fmt.Errorf("failed to marshal spec: %s", err)
	}
	spec = string(j)
	return
}

// Download download object to local, could be file or directory
func (a *agent) Download(ctx context.Context, input *COSInput) (err error) {
	return a.doTransfer(ctx, "download", input)
}

// Upload upload file or directory to COS
func (a *agent) Upload(ctx context.Context, input *COSInput) (err error) {
	// When uploading directory, the local path can't be relative path.
	// Transferd will raise no such file or directory error.
	// This should be a bug of transferd because there is no such problem with faspmanager2 backend.
	var absPath string
	if absPath, err = filepath.Abs(input.Path); err == nil {
		input.Path = absPath
	}
	return a.doTransfer(ctx, "upload", input)
}

func (a *agent) doTransfer(ctx context.Context, action string, input *COSInput) (err error) {
	rpcCtx := context.TODO()
	if err = a.StartServer(rpcCtx); err != nil {
		return
	}

	p := &sdk.Path{Source: input.Key, Destination: input.Path}
	if action == "upload" {
		p = &sdk.Path{Source: input.Path, Destination: input.Key}
	}

	transferSpec, err := a.GetAsperaTransferSpecV2(action, input.Bucket, []*sdk.Path{p})
	if err != nil {
		return
	}

	req := &sdk.TransferRequest{
		TransferSpec: transferSpec,
		Config: &sdk.TransferConfig{
			Retry: &sdk.RetryStrategy{},
		},
		TransferType: sdk.TransferType_FILE_REGULAR,
	}
	transferResp, err := a.StartTransfer(ctx, req)
	if err != nil {
		return
	}
	// `StartTransfer` method only returns error when the RPC message was not sent successfully,
	// it doesn't care if the transfer is accepted or not, so if we want to know the result,
	// we need to check the response explicitly.
	if e := transferResp.GetError(); e != nil {
		return fmt.Errorf("failed to start transfer: %d: %s", e.GetCode(), e.GetDescription())
	}
	transferId := transferResp.GetTransferId()

	go func() {
		<-ctx.Done()
		stop := &sdk.StopTransferRequest{TransferId: []string{transferId}}
		if _, err := a.StopTransfer(rpcCtx, stop); err != nil {
			log.Println("failed to stop transfer:", err)
		}
	}()

	stream, err := a.MonitorTransfers(rpcCtx, &sdk.RegistrationRequest{TransferId: []string{transferId}})
	if err != nil {
		return
	}

	if input.Sub == nil {
		input.Sub = &DefaultSubscriber{}
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch resp.Status {
		case sdk.TransferStatus_QUEUED:
			input.Sub.Queued(resp)
		case sdk.TransferStatus_RUNNING:
			input.Sub.Running(resp)
		case sdk.TransferStatus_FAILED, sdk.TransferStatus_CANCELED:
			description := strings.TrimSpace(resp.TransferInfo.GetErrorDescription())
			return fmt.Errorf("transfer %s: %s", resp.Status, description)
		case sdk.TransferStatus_COMPLETED:
			input.Sub.Done(resp)
			// MonitorTransfers doesn't works like StartTransferWithMonitor,
			// the response it returns doesn't emit EOF because of multiple transfers
			// so the loop will block infinitely even the transfer's finished.
			// I have to return here explicitly.
			return nil
		}
	}

	return
}
