package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"

	"github.com/IBM/ibmcloud-cos-cli/config"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam/token"
)

// Provider constants
const (
	providerName = "ic-bridge"
	providerType = "oauth"
)

/// /// /// Token Time /// /// ///
type tokenTime struct {
	IAt int64 `json:"iat"`
	Exp int64 `json:"exp"`

	*credentials.Expiry
}

// set the expiry of the token based on the token expire and emit values
func (tt *tokenTime) populateExpiry() {
	if tt != nil && tt.IAt != 0 && tt.Exp != 0 {
		iat := time.Unix(tt.IAt, 0)
		exp := time.Unix(tt.Exp, 0)

		// as SDK IAM Token Manager It expires at 17% of the TTL
		timeDelta := time.Duration(float64(exp.Sub(iat).Nanoseconds())*0.17) * time.Nanosecond

		if tt.Expiry == nil {
			tt.Expiry = new(credentials.Expiry)
		}

		tt.Expiry.SetExpiration(exp, timeDelta)
	}
}

// checks if the token is expired
// mandatory to be used as provider by sdk
func (tt *tokenTime) IsExpired() bool {
	// no expiration, the token is static / forever
	if tt == nil || tt.Expiry == nil {
		return false
	}
	return tt.Expiry.IsExpired()
}

/// /// /// /// /// /// ///

/// /// /// Bluemix / IBM Cloud CLI Bridge /// /// ///
type bxIAM struct {
	// Plugin Context
	pluginContext plugin.PluginContext

	// IAM Token
	token *token.Token

	// IAM Token Expiration time
	expiry *tokenTime
}

// Check if IAM Token is expired or not
func (bx *bxIAM) IsExpired() bool {
	return bx.expiry.IsExpired()
}

// Grabs current authentication token from cli and wrap it in an SDK Credentials
// mandatory to be used as provider by sdk
func (bx *bxIAM) Retrieve() (credentials.Value, error) {

	// Initialize token if not given
	if bx.token == nil {
		err := bx.init()
		if err != nil {
			return credentials.Value{}, err
		}
	}

	// Check if token is expired - refresh
	if bx.IsExpired() {
		err := bx.refresh()
		if err != nil {
			return credentials.Value{}, err
		}
	}

	// Obtain CRN from the Plugin Context
	crn, err := bx.pluginContext.PluginConfig().GetString(config.CRN)
	if err != nil {
		return credentials.Value{}, err
	}

	// Build Credentials value with token and provider and CRN (IAM)
	value := credentials.Value{
		Token:             *bx.token,
		ProviderName:      providerName,
		ProviderType:      providerType,
		ServiceInstanceID: crn,
	}

	return value, nil
}

// initialize the provider with current iam token from cli
func (bx *bxIAM) init() error {

	// Raw Token passed in from IBM Cloud CLI that
	// IAMToken returns the IAM access tokenValue
	if !bx.pluginContext.IsLoggedIn() {
		return awserr.New("auth.iam.logout", "Loggin Required", nil)
	}

	rawToken := bx.pluginContext.IAMToken()
	//fmt.Println(, bx.pluginContext.IsLoggedInWithServiceID())

	// Retains the times of the tokenValue
	timing, err := getTokenTimes(rawToken)
	if err != nil {
		return err
	}

	// Build Token locally for the CLI App
	// Token and its expiration time
	tokenValue, err := buildToken(rawToken)
	bx.token = tokenValue
	bx.expiry = timing
	return nil
}

// refresh the token when expired
func (bx *bxIAM) refresh() error {

	// Raw Token passed in from IBMCloud CLI that RefreshIAMToken
	// refreshes and returns the IAM access tokenValue
	rawToken, err := bx.pluginContext.RefreshIAMToken()
	if err != nil {
		return err
	}

	// Retains the times of the tokenValue
	timing, err := getTokenTimes(rawToken)
	if err != nil {
		return err
	}

	// Build Token locally for the CLI App
	// Token and its expiration time
	tokenValue, err := buildToken(rawToken)
	bx.token = tokenValue
	bx.expiry = timing
	return nil
}

// parse the JWT token parses it and gets its timings
func getTokenTimes(token string) (*tokenTime, error) {
	// Initialize Token time
	timing := new(tokenTime)

	// Split JWT token into three pieces by periods
	tokenSplit := strings.Split(token, ".")
	if len(tokenSplit) != 3 {
		return nil, awserr.New("auth.iam.token.decode", "Unable to decode IAM token", nil)
	}

	// Decode Token string
	bytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(tokenSplit[1])
	if err != nil {
		return nil, err
	}

	// Json unmarshal
	err = json.Unmarshal(bytes, timing)
	if err != nil {
		return nil, err
	}
	timing.populateExpiry()
	return timing, nil
}

// build an sdk token from the cli oauth token
func buildToken(rawToken string) (*token.Token, error) {
	// Build token into three values
	// Entire Token
	// IAM Access Token Type
	// IAM Access Token
	tokenSplit := strings.Split(rawToken, " ")
	if len(tokenSplit) != 2 {
		return nil, awserr.New("auth.iam.token.decode", "Unable to decode IAM token", nil)
	}
	tokenValue := new(token.Token)
	tokenValue.TokenType = tokenSplit[0]
	tokenValue.AccessToken = tokenSplit[1]

	return tokenValue, nil
}

/// /// /// /// /// /// ///
/// /// /// Wrappers to Credentials /// /// ///

// NewBxBridgeProvider helper to make it easier register the bridge as credentials provider
func NewBxBridgeProvider(pluginContext plugin.PluginContext) *bxIAM {
	return &bxIAM{pluginContext: pluginContext}
}

// NewBxBridgeCredentials helper to make it easier register the bridge as credentials in the config of the sdk
func NewBxBridgeCredentials(pluginContext plugin.PluginContext) *credentials.Credentials {
	return credentials.NewCredentials(NewBxBridgeProvider(pluginContext))
}
