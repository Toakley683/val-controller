package valorantapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

/* [ - Get Auth Cookies - ] */
/* [ URL: 'https://auth.riotgames.com/api/v1/authorization' ] */

type AuthCookie struct {
	Type      string `json:"type"`
	Country   string `json:"country"`
	SessionID *http.Cookie
	ClientID  *http.Cookie
}

func GetAuthCookie() (*AuthCookie, error) {

	JsonData :=
		`{
			"client_id":"play-valorant-web-prod",
			"nonce":"1",
			"redirect_uri": "https://playvalorant.com/opt_in",
			"response_type": "token id_token",
			"scope": "account openid"
		}`

	APICall := newAPICall()

	APICall.URL = "https://auth.riotgames.com/api/v1/authorization"
	APICall.Body = strings.NewReader(JsonData)
	APICall.Method = "POST"

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

	var SessionID *http.Cookie
	var ClientID *http.Cookie

	for _, v := range APICall.Response.Cookies() {
		switch v.Name {
		case "asid":
			SessionID = v
		case "clid":
			ClientID = v
		}
	}

	Response := &AuthCookie{}

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	Response.SessionID = SessionID
	Response.ClientID = ClientID

	return Response, nil

}

/* [ - Cookie Reauth - ] */
/* [ URL: 'https://auth.riotgames.com/authorize' ] */

/* Beware that his does require the Redirect URI to be verified in your applications portal 'https://developer.riotgames.com/' */
func CookieReauthenticate(redirect_uri string) (string, error) {

	params := url.Values{}

	params.Add("redirect_uri", redirect_uri)
	params.Add("client_id", "play-valorant-web-prod")
	params.Add("response_type", "token id_token")
	params.Add("nonce", "1")
	params.Add("scope", "account openid")

	FinalURL := "https://auth.riotgames.com/authorize?" + params.Encode()

	return FinalURL, nil

}

/* [ - Player Info - ] */
/* [ URL: 'https://auth.riotgames.com/userinfo' ] */

type AuthPlayerInfoFederatedDetails struct {
	ProviderName        string      `json:"provider_name"`
	ProviderEnvironment interface{} `json:"provider_environment"`
}

type AuthPlayerInfo struct {
	Country                    string                           `json:"country"`
	LocalUUID                  string                           `json:"sub"`
	EmailVerified              bool                             `json:"email_verified"`
	PlayerPlocale              interface{}                      `json:"player_plocale"`
	CountryAt                  int64                            `json:"country_at"`
	PhoneNumberVerified        bool                             `json:"phone_number_verified"`
	LinkedIdentityDetails      []interface{}                    `json:"linked_identity_details"`
	PreferredUsername          string                           `json:"preferred_username"`
	AccountVerified            bool                             `json:"account_verified"`
	Ppid                       interface{}                      `json:"ppid"`
	FederatedIdentityDetails   []AuthPlayerInfoFederatedDetails `json:"federated_identity_details"`
	FederatedIdentityProviders []string                         `json:"federated_identity_providers"`
	PlayerLocale               string                           `json:"player_locale"`
	EmailSet                   bool                             `json:"email_set"`
	Account                    RSOUserInfoAcct                  `json:"acct"`
	Age                        int                              `json:"age"`
	Jti                        string                           `json:"jti"`
	Username                   string                           `json:"username"`
	Affinity                   struct {
		Pp string `json:"pp"`
	} `json:"affinity"`
}

func (context ValorantAPIContext) GetLocalPlayerInfo() (*AuthPlayerInfo, error) {

	APICall := newAPICall()

	APICall.URL = "https://auth.riotgames.com/userinfo"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBearerAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &AuthPlayerInfo{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}
