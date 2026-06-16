package valorantapi

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	client http.Client
)

func init() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = http.Client{Transport: tr}

}

type valorantApiCall struct {
	URL      string
	Body     io.Reader
	Method   string
	Request  *http.Request
	Response *http.Response
}

/* --[[ API Requesters ]]-- */

func newAPICall() *valorantApiCall {

	call := valorantApiCall{}

	return &call

}

func (call *valorantApiCall) SetBasicAuth(context ValorantAPIContext) {

	auth := "riot:" + context.localLockfile.Password
	call.AddHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))

}

func (call *valorantApiCall) SetBearerAuth(context ValorantAPIContext) error {

	entitlement, err := context.GetEntitlementToken()
	if err != nil {
		return err
	}

	clientPlatform, err := GetClientPlatformBase64()
	if err != nil {
		return err
	}

	clientVersion, err := context.GetValorantVersion()
	if err != nil {
		return err
	}

	call.AddHeader("Authorization", "Bearer "+entitlement.AccessToken)
	call.AddHeader("X-Riot-Entitlements-JWT", entitlement.Token)
	call.AddHeader("X-Riot-ClientPlatform", clientPlatform)
	call.AddHeader("X-Riot-ClientVersion", clientVersion.Data.RiotClientVersion)

	return nil

}

func GetClientPlatformBase64() (string, error) {

	var client_platform = map[string]string{
		"platformType":      "PC",
		"platformOS":        "Windows",
		"platformOSVersion": "10.0.19042.1.256.64bit",
		"platformChipset":   "Unknown",
	}

	client_platform_json, err := json.Marshal(client_platform)

	if err != nil {
		return "", err
	}

	client_platform_base64 := base64.StdEncoding.EncodeToString(client_platform_json)

	return client_platform_base64, nil
}

func (call *valorantApiCall) AddCookies(cookies []*http.Cookie) {

	for _, v := range cookies {
		call.Request.AddCookie(v)
	}

}

func (call *valorantApiCall) AddHeader(Key, Value string) {

	call.Request.Header.Set(Key, Value)
}

func (call *valorantApiCall) SetRequest() error {

	req, err := http.NewRequest(call.Method, call.URL, call.Body)
	if err != nil {
		return err
	}

	call.Request = req

	return nil

}

func (call *valorantApiCall) Call() error {

	var err error

	call.Request.Header.Set("Content-Type", "application/json")

	call.Response, err = client.Do(call.Request)
	if err != nil {
		return err
	}

	if call.Response.StatusCode == 429 {
		return errors.New(VALORANT_API_TIMEOUT)
	}

	if call.Response.StatusCode == 401 {
		return errors.New(VALORANT_API_UNAUTHORIZED)
	}

	if call.Response.StatusCode != 200 {
		return errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return nil

}

/*
func requestAPI(url string, body io.Reader, method string, cookies []*http.Cookie) (*http.Response, error) {

	fmt.Println("Requested:", "'"+url+"'")

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for _, v := range cookies {
		req.AddCookie(v)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if response.StatusCode == 429 {
		return nil, errors.New(VALORANT_API_TIMEOUT)
	}

	return response, nil

}*/
