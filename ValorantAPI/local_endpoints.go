package valorantapi

import (
	"encoding/json"
	"io"
	"strings"
)

/* [ - Get Local Help - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/help' ] */

type LocalHelp struct {
	Events    map[string]string `json:"events"`
	Functions map[string]string `json:"functions"`
	Types     map[string]string `json:"types"`
}

func (context ValorantAPIContext) GetLocalHelp() (*LocalHelp, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/help"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &LocalHelp{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Sessions - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/product-session/v1/external-sessions' ] */

type LaunchConfiguration struct {
	Arguments        []string `json:"arguments"`
	Executable       string   `json:"executable"`
	Locale           *string  `json:"locale"`      // nullable
	VoiceLocale      *string  `json:"voiceLocale"` // nullable
	WorkingDirectory string   `json:"workingDirectory"`
}

type SessionDetail struct {
	ExitCode            int                 `json:"exitCode"`
	ExitReason          *string             `json:"exitReason"` // nullable
	IsInternal          bool                `json:"isInternal"`
	LaunchConfiguration LaunchConfiguration `json:"launchConfiguration"`
	PatchlineFullName   string              `json:"patchlineFullName"` // "VALORANT" or "riot_client"
	PatchlineId         string              `json:"patchlineId"`       // "" | "live" | "pbe"
	Phase               string              `json:"phase"`
	ProductId           string              `json:"productId"` // "valorant" or "riot_client"
	Version             string              `json:"version"`
}

type SessionsResponse map[string]SessionDetail

func (context ValorantAPIContext) GetSessions() (*SessionsResponse, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/product-session/v1/external-sessions"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &SessionsResponse{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get RSOUserInfo - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/rso-auth/v1/authorization/userinfo' ] */

type RSOUserInfoWrapper struct {
	UserInfo string `json:"userInfo"`
}

type RSOUserInfo struct {
	Country             string          `json:"country"`
	PlayerUUID          string          `json:"sub"`
	EmailVerified       bool            `json:"email_verified"`
	PlayerPLocale       interface{}     `json:"player_plocale"`
	CountryAt           int64           `json:"country_at"`
	Pw                  RSOUserInfoPw   `json:"pw"`
	Lol                 interface{}     `json:"lol"`
	OriginalPlatformID  interface{}     `json:"original_platform_id"`
	OriginalAccountID   interface{}     `json:"original_account_id"`
	PhoneNumberVerified bool            `json:"phone_number_verified"`
	PreferredUsername   string          `json:"preferred_username"`
	Ban                 RSOUserInfoBan  `json:"ban"`
	Ppid                interface{}     `json:"ppid"`
	LolRegion           []string        `json:"lol_region"`
	PlayerLocale        string          `json:"player_locale"`
	PvpnetAccountID     interface{}     `json:"pvpnet_account_id"`
	Acct                RSOUserInfoAcct `json:"acct"`
	Jti                 string          `json:"jti"`
	Username            string          `json:"username"`
}

type RSOUserInfoPw struct {
	CngAt     int64 `json:"cng_at"`
	Reset     bool  `json:"reset"`
	MustReset bool  `json:"must_reset"`
}

type RSOUserInfoBan struct {
	Restrictions []interface{} `json:"restrictions"`
}

type RSOUserInfoAcct struct {
	Type      int    `json:"type"`
	State     string `json:"state"`
	Adm       bool   `json:"adm"`
	GameName  string `json:"game_name"`
	TagLine   string `json:"tag_line"`
	CreatedAt int64  `json:"created_at"`
}

func (context ValorantAPIContext) GetRSOUserInfo() (*RSOUserInfo, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/rso-auth/v1/authorization/userinfo"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	wrapper := &RSOUserInfoWrapper{}

	err = json.Unmarshal(FullBody, &wrapper)
	if err != nil {
		return nil, err
	}

	Response := &RSOUserInfo{}

	err = json.Unmarshal([]byte(wrapper.UserInfo), Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Locale - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/riotclient/region-locale' ] */

type RegionLocale struct {
	Locale      string `json:"locale"`
	Region      string `json:"region"`
	WebLanguage string `json:"webLanguage"`
	WebRegion   string `json:"webRegion"`
}

func (context ValorantAPIContext) GetRegionLocale() (*RegionLocale, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/riotclient/region-locale"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &RegionLocale{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Account Alias - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/player-account/aliases/v1/active' ] */

type AccountAlias struct {
	Active          bool   `json:"active"`
	CreatedDatetime int64  `json:"created_datetime"`
	GameName        string `json:"game_name"`
	Summoner        bool   `json:"summoner"`
	TagLine         string `json:"tag_line"`
}

func (context ValorantAPIContext) GetAccountAlias() (*AccountAlias, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/player-account/aliases/v1/active"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &AccountAlias{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Entitlement Token - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/entitlements/v1/token' ] */

type EntitlementToken struct {
	AccessToken  string        `json:"accessToken"`
	Entitlements []interface{} `json:"entitlements"`
	Issuer       string        `json:"issuer"`
	Subject      string        `json:"subject"`
	Token        string        `json:"token"`
}

func (context ValorantAPIContext) GetEntitlementToken() (*EntitlementToken, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/entitlements/v1/token"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &EntitlementToken{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Chat Session - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/chat/v1/session' ] */

type ChatSession struct {
	Federated  bool   `json:"federated"`
	GameName   string `json:"game_name"`
	GameTag    string `json:"game_tag"`
	Loaded     bool   `json:"loaded"`
	Name       string `json:"name"`
	Pid        string `json:"pid"`
	PlayerUUID string `json:"puuid"`
	Region     string `json:"region"`
	Resource   string `json:"resource"`
	State      string `json:"state"`
}

func (context ValorantAPIContext) GetChatSession() (*ChatSession, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/chat/v1/session"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &ChatSession{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Friends - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/chat/v4/friends' ] */

type LocalFriend struct {
	ActivePlatform string `json:"activePlatform"`
	DisplayGroup   string `json:"displayGroup"`
	GameName       string `json:"game_name"`
	GameTag        string `json:"game_tag"`
	Group          string `json:"group"`
	LastOnlineTs   int64  `json:"last_online_ts"`
	Name           string `json:"name"`
	Note           string `json:"note"`
	Pid            string `json:"pid"`
	PlayerUUID     string `json:"puuid"`
	Region         string `json:"region"`
}

type LocalFriendsList struct {
	Friends []LocalFriend `json:"friends"`
}

func (context ValorantAPIContext) GetFriends() (*LocalFriendsList, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/chat/v4/friends"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &LocalFriendsList{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Send Friend Request - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/chat/v4/friendrequests' ] */

func (context ValorantAPIContext) SendFriendRequest(GameName, GameTag string) error {

	JsonData :=
		`{
			"game_name":"` + GameName + `",
			"game_tag":"` + GameTag + `"
		}`

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/chat/v4/friends"
	APICall.Body = strings.NewReader(JsonData)
	APICall.Method = "POST"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return err
	}

	defer APICall.Response.Body.Close()

	return nil

}

/* [ - Remove Friend Request - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/chat/v4/friendrequests' ] */

func (context ValorantAPIContext) RemoveFriendRequest(PUUID string) error {

	JsonData :=
		`{
			"puuid":"` + PUUID + `"
		}`

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/chat/v4/friends"
	APICall.Body = strings.NewReader(JsonData)
	APICall.Method = "DELETE"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return err
	}

	defer APICall.Response.Body.Close()

	return nil

}

/* [ - Get Presences - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/chat/v4/presences' ] */

type LocalPresence struct {
	ActivePlatform string        `json:"activePlatform"`
	Actor          interface{}   `json:"actor"`
	Basic          string        `json:"basic"`
	Details        interface{}   `json:"details"`
	GameName       string        `json:"game_name"`
	GameTag        string        `json:"game_tag"`
	Location       interface{}   `json:"location"`
	Msg            interface{}   `json:"msg"`
	Name           string        `json:"name"`
	Parties        []interface{} `json:"parties"`
	Patchline      interface{}   `json:"patchline"`
	Pid            string        `json:"pid"`
	Platform       string        `json:"platform"`
	Private        string        `json:"private"`
	PrivateJwt     interface{}   `json:"privateJwt"`
	Product        string        `json:"product"`
	PlayerUUID     string        `json:"puuid"`
	Region         string        `json:"region"`
	Resource       string        `json:"resource"`
	State          string        `json:"state"`
	Summary        string        `json:"summary"`
	Time           int64         `json:"time"`
}

type LocalPresences struct {
	Presences []LocalPresence `json:"presences"`
}

func (context ValorantAPIContext) GetPresences() (*LocalPresences, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/chat/v4/presences"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &LocalPresences{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}

/* [ - Get Friend Requests - ] */
/* [ URL: 'https://127.0.0.1:{ Local Lockfile Port }/chat/v4/friendrequests' ] */

type FriendRequest struct {
	GameName     string `json:"game_name"`
	GameTag      string `json:"game_tag"`
	Name         string `json:"name"`
	Note         string `json:"note"`
	Pid          string `json:"pid"`
	Platform     string `json:"platform"`
	PlayerUUID   string `json:"puuid"`
	Region       string `json:"region"`
	Subscription string `json:"subscription"`
}

type FriendRequests struct {
	Requests []FriendRequest `json:"requests"`
}

func (context ValorantAPIContext) GetFriendRequests() (*FriendRequests, error) {

	APICall := newAPICall()

	APICall.URL = "https://127.0.0.1:" + context.localLockfile.Port + "/chat/v4/friendrequests"
	APICall.Body = nil
	APICall.Method = "GET"

	var err error

	err = APICall.SetRequest()
	if err != nil {
		return nil, err
	}

	APICall.SetBasicAuth(context)

	err = APICall.Call()
	if err != nil {
		return nil, err
	}

	defer APICall.Response.Body.Close()

	FullBody, err := io.ReadAll(APICall.Response.Body)
	if err != nil {
		return nil, err
	}

	Response := &FriendRequests{}

	err = json.Unmarshal(FullBody, &Response)
	if err != nil {
		return nil, err
	}

	return Response, nil

}
