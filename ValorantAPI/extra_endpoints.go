package valorantapi

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

/* [ - Get Valorant Version - ] */
/* [ URL: 'https://valorant-api.com/v1/version' ] */

type ExtraValorantVersion struct {
	Status int `json:"status"`
	Data   struct {
		ManifestID        string    `json:"manifestId"`
		Branch            string    `json:"branch"`
		Version           string    `json:"version"`
		BuildVersion      string    `json:"buildVersion"`
		EngineVersion     string    `json:"engineVersion"`
		RiotClientVersion string    `json:"riotClientVersion"`
		RiotClientBuild   string    `json:"riotClientBuild"`
		BuildDate         time.Time `json:"buildDate"`
	} `json:"data"`
}

func (context ValorantAPIContext) GetValorantVersion() (*ExtraValorantVersion, error) {

	APICall := newAPICall()

	APICall.URL = "https://valorant-api.com/v1/version"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &ExtraValorantVersion{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	if Response.Status != 200 {
		return nil, errors.New(VALORANT_NO_VALID_RESPONSE)
	}

	return Response, nil

}
