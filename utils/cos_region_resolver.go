package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"
	"github.com/IBM/ibm-cos-sdk-go/aws/endpoints"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
)

// ListKnownRegions interface to hold the GetAllRegions operation
type ListKnownRegions interface {
	// GetAllRegions function to get all known regions
	GetAllRegions() ([]string, error)
}

var (
	// DefaultCOSEndPointsWS default regions endpoint url
	DefaultCOSEndPointsWS = `https://control.cloud-object-storage.cloud.ibm.com/v2/endpoints`

	// Classes ibm cos storage classes
	// needed to disambiguation between us-south vs us-geo vs us-flex in region constraint
	Classes        = []string{"standard", "vault", "cold", "flex", "smart"}
	classTokenizer = strings.Join(Classes, `|`)
	regionSpec     = fmt.Sprintf(`(?i)^(\w+(?:-\w+)??)(-geo)?(?:-(%s|\*))?$`, classTokenizer)
	// RegionDecoderRegex regular expression to break down regions and storage class
	RegionDecoderRegex = regexp.MustCompile(regionSpec)
)

// IBMEndPoints maps the struct of the endpoints region
type IBMEndPoints struct {
	// Service urls organised by resiliency level
	Service *ResiliencyEP `json:"service-endpoints"`
	// Identity service url
	Identity *Identity `json:"identity-endpoints"`
}

// Identity maps part of the struct of the endpoints region
type Identity struct {
	// IamToken service url
	IamToken string `json:"iam-token"`
	// IamPolicy service url
	IamPolicy string `json:"iam-policy"`
}

// ResiliencyEP maps part of the struct of the endpoints region
type ResiliencyEP struct {
	// CrossRegion service urls
	CrossRegion LocationEP `json:"cross-region"`
	// Regional service urls
	Regional LocationEP `json:"regional"`
	// SingleSite service urls
	SingleSite LocationEP `json:"single-site"`
}

// LocationEP maps part of the struct of the endpoints region
type LocationEP map[string]PublicPrivateEP

// PublicPrivateEP maps part of the struct of the endpoints region
type PublicPrivateEP struct {
	// Public service urls
	Public map[string]string `json:"public"`
	// Private service urls
	Private map[string]string `json:"private"`
}

// COSEndPointsWSClient represents a client to be used to retrieve regions endpoint value
type COSEndPointsWSClient struct {
	client       *http.Client
	EndpointsURL string
	endPoints    *IBMEndPoints

	disableSSL bool

	mutex     sync.RWMutex
	lastFetch *time.Time
	etag      string

	time.Ticker // ticker to refresh every N time ...
}

// NewIBMEndPoints creates a new client receiving the config and regions endpoint url as parameters
func NewIBMEndPoints(config *aws.Config, endpointsURL string) (*COSEndPointsWSClient, error) {
	// allocates memory for the client struct
	c := COSEndPointsWSClient{}

	// checks if a regions endpoint url was provided, if not provided use default value
	c.EndpointsURL = endpointsURL
	if c.EndpointsURL == "" {
		c.EndpointsURL = DefaultCOSEndPointsWS
	}

	// check if a config was provided
	if config != nil {
		// if config provided use its httpclient
		c.client = config.HTTPClient
	}
	// if value of http client is empty
	// use default httpclient
	if c.client == nil {
		c.client = http.DefaultClient
	}

	//c.endPoints = &IBMEndPoints{}
	//err := c.fetchEndPointsDetails()
	//if err != nil {
	//	return nil, err
	//}

	// check if config provided
	// if config provided apply its SSL policy
	disableSSL := aws.Bool(false)
	if config != nil {
		disableSSL = config.DisableSSL
	}
	c.disableSSL = aws.BoolValue(disableSSL)

	return &c, nil
}

// connects to the endpoint and retrieves the values
func (c *COSEndPointsWSClient) fetchEndPointsDetails() (err error) {
	// full lock of the client struct to avoid race conditions accessing it
	c.mutex.Lock()
	// defer the unlock to guarantee it will be unlock when function exit
	defer c.mutex.Unlock()

	// builds a new get request
	var request *http.Request
	request, err = http.NewRequest(http.MethodGet, c.EndpointsURL, http.NoBody)
	// check if error on build get request
	if err != nil {
		return
	}

	//if c.etag != "" {
	//	request.Header.Set("If-None-Match", c.etag)
	//}
	// // Commented Out because the server does not comply
	// // https://tools.ietf.org/html/rfc7232#section-3.3
	// // "A recipient MUST ignore If-Modified-Since if the request contains an If-None-Match header field;"
	// // After Testing 200/304 does not bring performance differences
	//if c.lastFetch != nil {
	//	request.Header.Set("If-Modified-Since",c.lastFetch.Format(http.TimeFormat) )
	//}

	// do the request
	var response *http.Response
	response, err = c.client.Do(request)
	// check if error on request
	if nil != err {
		return
	}
	// defer a function to flush the request buffers
	defer func() {
		io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()
	// update client struct last fetch time
	now := time.Now()
	if response.StatusCode == 304 {
		// no change return
		c.lastFetch = &now
		return
	}

	// check response code
	if 200 < response.StatusCode || 300 <= response.StatusCode {
		err = fmt.Errorf("url:'%s' Status: %s", c.EndpointsURL, response.Status)
		return
	}
	// read response body to a buffer
	var data []byte
	data, err = ioutil.ReadAll(response.Body)
	if nil != err {
		return
	}
	// parses the content of the response buffer
	tmp := &IBMEndPoints{}
	err = json.Unmarshal(data, tmp)
	if nil != err {
		return
	}
	c.endPoints = tmp
	// check the parse of the buffer was successfully
	if c.endPoints.Service == nil || c.endPoints.Identity == nil {
		c.endPoints = nil
		err = awserr.New("ibm.Resolver.ErrorFetchingEndpoints", "Error Fetching And Parsing Endpoints JSON", nil)
	}
	// update last last check values ( check time and ETag )
	c.lastFetch = &now
	c.etag = response.Header.Get("ETag")
	return
}

// Refresh refresh the values fetched earlier
func (c *COSEndPointsWSClient) Refresh() (err error) {
	return c.fetchEndPointsDetails()
}

// EndpointFor method required by the sdk to use the client as sdk Resolver
func (c *COSEndPointsWSClient) EndpointFor(service, region string,
	opts ...func(options *endpoints.Options)) (endpoint endpoints.ResolvedEndpoint, err error) {
	// apply a read lock to guarantee other can read but no write can be dne until lock released
	c.mutex.RLock()
	// defer the release of the read lock
	defer c.mutex.RUnlock()

	// if no endpoints present
	// release read lock
	// acquire a full lock
	if c.endPoints == nil {
		c.mutex.RUnlock()
		err := c.fetchEndPointsDetails()
		c.mutex.RLock()
		if err != nil {
			return endpoint, err
		}
	}
	// check the request is for S3 service
	if service == s3.ServiceName {
		// use the RegionDecoderRegex to break down region and class storage if present
		regionDetails := RegionDecoderRegex.FindStringSubmatch(region)
		if regionDetails != nil {
			//region,geo,class := regionDetails[1],regionDetails[2],regionDetails[3]
			region, geo := regionDetails[1], regionDetails[2]
			// set -geo value for when look up the cross region
			if geo == "" {
				geo = "-geo"
			}
			var url string
			c.mutex.RLock()
			defer c.mutex.RUnlock()
			url = c.endPoints.Service.CrossRegion[region].Public[region+geo]
			if url == "" {
				url = c.endPoints.Service.Regional[region].Public[region]
			}
			if url == "" {
				url = c.endPoints.Service.SingleSite[region].Public[region]
			}
			if url != "" {

				var opt endpoints.Options
				opt.DisableSSL = c.disableSSL
				opt.Set(opts...)

				var prefix string
				if opt.DisableSSL {
					prefix = "http://"
				} else {
					prefix = "https://"
				}

				endpoint.URL = prefix + url
				endpoint.SigningRegion = region
				endpoint.SigningName = service
				return
			}
			// return error if region not found
			return endpoint, awserr.New("ibm.Resolver.RegionNotFound", "Region Not Found: "+region, nil)
		}
		// return error if region parse error
		return endpoint, awserr.New("ibm.Resolver.ErrorParsingRegion", "Error parsing Region: "+region, nil)
	}
	// return error if service not s3
	return endpoint, awserr.New("ibm.Resolver.ServiceNotSupported", service+": Service Not Supported", nil)
}

// GetAllRegions method that returns all known regions
// useful when looking for a bucket location in unknown region
func (c *COSEndPointsWSClient) GetAllRegions() ([]string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	// check if regions present
	// if not fetch them
	if c.endPoints == nil {
		c.mutex.RUnlock()
		err := c.fetchEndPointsDetails()
		c.mutex.RLock()
		if err != nil {
			return nil, err
		}
	}
	if c.endPoints.Service == nil {
		return []string{}, nil
	}

	// allocate memory for a slice for all regions
	regions := make([]string,
		len(c.endPoints.Service.CrossRegion)+len(c.endPoints.Service.Regional)+len(c.endPoints.Service.SingleSite))
	idx := 0

	// add cross region regions to slice
	for region := range c.endPoints.Service.CrossRegion {
		regions[idx] = region
		idx++
	}

	// add regional regions to slice
	for region := range c.endPoints.Service.Regional {
		regions[idx] = region
		idx++
	}

	// add single site regions to slice
	for region := range c.endPoints.Service.SingleSite {
		regions[idx] = region
		idx++
	}

	return regions, nil
}
